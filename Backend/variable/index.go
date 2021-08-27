package variable

import (
	"context"
	"database/sql"
	"feorm/config"
	"feorm/db"
	"feorm/normalize"
	"fmt"
	"github.com/bluele/gcache"
	"regexp"
	"time"
)

type Variable struct {
	name         string
	sqlQuery     string
	ttl          int
	cache        gcache.Cache
	arguments    []string
	db           db.DB
	defaultValue interface{}
}

var variableMap = make(map[string]*Variable)

var variableRegexp = regexp.MustCompile(`@(\w+)`)

func processSql(s string, placeHolderFunc func(idx int) string) (string, []string, error) {
	varMatches := variableRegexp.FindAllStringSubmatch(s, -1)
	vars := make([]string, 0, len(varMatches))
	for _, v := range varMatches {
		vars = append(vars, v[1])
	}
	var index = 0
	processed := variableRegexp.ReplaceAllStringFunc(s, func(s string) string {
		index++
		return placeHolderFunc(index - 1)
	})
	return processed, vars, nil
}

func New(config config.VariableDefine, db db.DB) (*Variable, error) {
	query, arguments, err := processSql(config.Sql, db.PlaceHolderFunc())
	if err != nil {
		return nil, err
	}
	v := Variable{
		db:           db,
		name:         config.Name,
		sqlQuery:     query,
		ttl:          config.Expire,
		cache:        gcache.New(4096).ARC().Build(),
		arguments:    arguments,
		defaultValue: config.Default,
	}
	variableMap[v.name] = &v
	return &v, nil
}

func (v *Variable) fetchData(uid string) (interface{}, error) {
	variables := make([]interface{}, len(v.arguments))
	for i := range v.arguments {
		var err error
		variables[i], err = GetVariable(v.arguments[i], uid)
		if err != nil {
			return nil, err
		}
	}
	row := v.db.QueryRow(context.Background(), v.sqlQuery, variables...)
	var ret interface{}
	err := row.Scan(&ret)
	if err == sql.ErrNoRows {
		return v.defaultValue, nil
	}
	return normalize.NormalizeType(ret), err
}
func (v *Variable) get(uid string) (interface{}, error) {
	return v.cache.Get(uid)
}

func GetVariable(name string, uid string) (interface{}, error) {
	if name == "uid" {
		return uid, nil
	}
	v := variableMap[name]
	if v == nil {
		return nil, fmt.Errorf("variable %s 未定义", name)
	}
	if v.ttl == 0 {
		return v.fetchData(uid)
	}
	ret, err := v.get(uid)
	if err != nil && err == gcache.KeyNotFoundError {
		ret, err = v.fetchData(uid)
		if err != nil {
			return nil, err
		}
		if v.ttl > 0 {
			_ = v.cache.SetWithExpire(uid, ret, time.Duration(v.ttl)*time.Second)
		}
	} else if err != nil {
		panic(err)
	}
	return ret, nil
}

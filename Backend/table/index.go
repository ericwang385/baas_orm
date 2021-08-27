package table

import (
	"context"
	"database/sql"
	"feorm/config"
	"feorm/db"
	"feorm/expr"
	"strings"
)

var TableMap = map[string]*Table{}

type Table struct {
	config.TableDefine
	selectGuardAst     *expr.AstNode
	updateGuardAst     *expr.AstNode
	deleteGuardAst     *expr.AstNode
	insertGuardProgram *expr.Program

	columnRules []*columnRule
	db          db.DB
	//  包含数据和rule需要的
	AllColumns       []string
	allColumnsIdx    map[string]int
	allColumnsString string
	columnInfo       []sql.ColumnType
	lazyColumnInfo   []sql.ColumnType
	allColumnInfo    []sql.ColumnType

	pkColumn      []string
	lazyColumn    map[string]interface{}
	returnColumns map[string]bool
}

func InitTable(conf config.TableDefine, db db.DB) error {
	selectGuardAst, err := expr.CompileAst(conf.SelectGuard)
	if err != nil {
		return err
	}
	updateGuardAst, err := expr.CompileAst(conf.UpdateGuard)
	if err != nil {
		return err
	}
	deleteGuardAst, err := expr.CompileAst(conf.DeleteGuard)
	if err != nil {
		return err
	}
	insertGuardProgram, err := expr.Compile(conf.InsertGuard)
	if err != nil {
		return err
	}
	rules := make([]*columnRule, len(conf.ColumnRules))
	for i := range rules {
		var err error
		rules[i], err = newColumnRule(conf.ColumnRules[i])
		if err != nil {
			return err
		}
	}
	t := &Table{
		TableDefine:        conf,
		selectGuardAst:     selectGuardAst,
		updateGuardAst:     updateGuardAst,
		deleteGuardAst:     deleteGuardAst,
		insertGuardProgram: insertGuardProgram,
		columnRules:        rules,
		db:                 db,
	}
	TableMap[t.Schema+"."+t.Name] = t
	columns, err := t.db.GetTableColumns(context.Background(), t.Schema, t.Name)
	if err != nil {
		return err
	}
	columns_info, err := t.db.GetTableColumnT(context.Background(), t.Schema, t.Name)
	if err != nil {
		return err
	}
	t.pkColumn, err = t.db.GetPrimaryKeys(context.Background(), t.Schema, t.Name)
	if err != nil {
		return err
	}
	t.AllColumns = make([]string, 0)
	t.columnInfo = make([]sql.ColumnType, 0)
	for i, c := range columns {
		for _, cc := range t.HiddenColumns {
			if c == cc {
				goto end
			}
		}
		for _, cc := range t.LazyColumns {
			if c == cc {
				t.lazyColumnInfo = append(t.lazyColumnInfo, *columns_info[i])
				t.AllColumns = append(t.AllColumns, c)
				goto end
			}
		}
		t.AllColumns = append(t.AllColumns, c)
		t.columnInfo = append(t.columnInfo, *columns_info[i])
	end:
	}
	t.allColumnsString = strings.Join(t.AllColumns, ",")
	t.allColumnsIdx = map[string]int{}
	t.allColumnInfo = append(t.columnInfo, t.lazyColumnInfo...)
	for i, c := range t.AllColumns {
		t.allColumnsIdx[c] = i
	}
	return nil
}

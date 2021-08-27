package table

import (
	"database/sql"
	"errors"
	"strings"
)

func (t *Table) Update(pkcol string, colName []string, value []string) (sql.Result, error) {
	if pkcol == "" {
		return nil, errors.New("update op missing PrimaryKey")
	}
	if value == nil {
		return nil, errors.New("update op missing Value")
	}
	updateBody := make([]string, 0)
	if len(value) == len(colName) {
		for i := 0; i < len(value); i++ {
			updateBody = append(updateBody, strings.Join([]string{colName[i], "=", value[i]}, ""))

		}
	}
	sql := strings.Join([]string{
		"UPDATE",
		t.Schema + "." + t.Name,
		"SET",
		strings.Join(updateBody, ","),
		"WHERE",
		t.pkColumn[0] + "=" + pkcol,
	}, " ")
	out, err := t.db.ExecQuery(sql, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

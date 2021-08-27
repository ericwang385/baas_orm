package table

import (
	"database/sql"
	"errors"
	"strings"
)

func (t *Table) Delete(pkcol string, colName []string, value []string) (sql.Result, error) {
	if pkcol == "" {
		return nil, errors.New("delete op missing PrimaryKey")
	}
	if value == nil {
		return nil, errors.New("delete op missing Value")
	}
	updateBody := make([]string, 0)
	if len(value) == len(colName) {
		for i := 0; i < len(value); i++ {
			updateBody = append(updateBody, strings.Join([]string{colName[i], "=", value[i]}, ""))

		}
	}
	sqlStr := strings.Join([]string{
		"DELETE FROM",
		t.Schema + "." + t.Name,
		"WHERE",
		t.pkColumn[0] + "=" + pkcol,
	}, " ")
	out, err := t.db.ExecQuery(sqlStr, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

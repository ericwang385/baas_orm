package table

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"
)

func (t *Table) Lazy(pkcol string, colName string) (*sqlx.Rows, error) {
	if pkcol == "" {
		return nil, errors.New("lazyload op missing PrimaryKey")
	}
	if colName == "" {
		return nil, errors.New("lazyload op missing ColName")
	}

	sql := strings.Join([]string{
		"SELECT " + colName + " FROM",
		t.Schema + "." + t.Name,
		"WHERE",
		t.pkColumn[0] + "=" + pkcol,
	}, " ")
	out, err := t.db.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	return out, nil
}

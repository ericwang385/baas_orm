package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var driversMap = map[string]func(dsn string) (DB, error){}

type DB interface {
	PlaceHolderFunc() func(idx int) string
	Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	GetTableColumns(ctx context.Context, schema, table string) ([]string, error)
	GetTableColumnT(ctx context.Context, schema, table string) ([]*sql.ColumnType, error)
	GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error)
	ExecQuery(query string, args ...interface{}) (sql.Result, error)
}

func Register(name string, factory func(dsn string) (DB, error)) {
	if driversMap[name] != nil {
		panic("duplicated  driver " + name)
	}
	driversMap[name] = factory
}

func Open(driver string, dsn string) (DB, error) {
	factory := driversMap[driver]
	if factory == nil {
		return nil, fmt.Errorf("driver %s not  found", driver)
	}
	return factory(dsn)
}

package db

import (
	"context"
	"database/sql"
	"feorm/util"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type postgres struct {
	db *sqlx.DB
}

func (p *postgres) GetTableColumns(ctx context.Context, schema, table string) ([]string, error) {
	rows, err := p.db.QueryxContext(ctx, fmt.Sprintf("select * from %s.%s limit 0",
		util.SanitizeSqlObject(schema), util.SanitizeSqlObject(table)))
	if err != nil {
		return nil, err
	}
	return rows.Columns()
}

func (p *postgres) GetTableColumnT(ctx context.Context, schema, table string) ([]*sql.ColumnType, error) {
	rows, err := p.db.QueryxContext(ctx, fmt.Sprintf("select * from %s.%s limit 0",
		util.SanitizeSqlObject(schema), util.SanitizeSqlObject(table)))
	if err != nil {
		return nil, err
	}
	return rows.ColumnTypes()
}

func (p *postgres) GetPrimaryKeys(ctx context.Context, schema, table string) ([]string, error) {
	query := fmt.Sprintf("SELECT a.attname\nFROM   pg_index i\nJOIN   pg_attribute a ON a.attrelid = i.indrelid\nAND a.attnum = ANY(i.indkey)\nWHERE  i.indrelid = '%s'::regclass\nAND    i.indisprimary;", table)
	rows, err := p.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for rows.Next() {
		r, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}
		out = append(out, r[0].(string))
	}
	return out, nil
}

func (p *postgres) PlaceHolderFunc() func(idx int) string {
	return func(idx int) string {
		return "$" + strconv.Itoa(idx+1)
	}
}

func (p *postgres) Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return p.db.QueryxContext(ctx, query, args...)
}

func (p *postgres) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return p.db.QueryRowxContext(ctx, query, args...)
}

func (p *postgres) ExecQuery(query string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(query)
}

func init() {
	Register("postgres", func(dsn string) (DB, error) {
		d, err := sqlx.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}
		return &postgres{db: d}, nil
	})
}

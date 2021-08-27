package table

import (
	"context"
	"feorm/expr"
	"feorm/normalize"
	"feorm/query"
	"feorm/util"
	"strconv"
	"strings"
)

func (t *Table) Fetch(query *query.Ast, args []interface{}, uid string, limit int, offset int, orderBy string,
	desc bool) ([][]interface{}, error) {
	var filter string
	var idx int
	if query == nil {
		filter = ""
	} else {
		var err error
		filter, idx, err = query.Compile2Sql(t.db.PlaceHolderFunc())
		if err != nil {
			return nil, err
		}
	}
	selectGuard, varFunc := t.selectGuardAst.CompileToSql(idx, t.db.PlaceHolderFunc())
	filters := make([]string, 0, 2)
	if filter != "" {
		filters = append(filters, filter)
	}
	if selectGuard != "" {
		filters = append(filters, selectGuard)
	}
	filterClause := strings.Join(filters, " AND ")
	if filterClause != "" {
		filterClause = "WHERE " + filterClause
	}
	var orderByClause string
	if orderBy != "" {
		if desc {
			orderByClause = "ORDER BY " + util.SanitizeSqlObject(orderBy) + " DESC"
		} else {
			orderByClause = "ORDER BY " + util.SanitizeSqlObject(orderBy)
		}
	}
	var limitClause string
	if limit > 0 {
		limitClause = "LIMIT " + strconv.Itoa(limit)
	}
	var offsetClause string
	if offset > 0 {
		offsetClause = "OFFSET " + strconv.Itoa(offset)
	}
	sql := strings.Join([]string{
		"SELECT",
		t.allColumnsString,
		"FROM",
		t.Schema + "." + t.Name,
		filterClause,
		orderByClause,
		limitClause,
		offsetClause,
	}, " ")
	vars, err := varFunc(uid)
	if err != nil {
		return nil, err
	}
	rows, err := t.db.Query(context.Background(), sql, append(args, vars...)...)
	if err != nil {
		return nil, err
	}
	dataRows := make([][]interface{}, 0, 8)
	for rows.Next() {
		r, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}
		dataRows = append(dataRows, r)
	}
	result, err := t.buildData(dataRows, uid)
	return result, err
}

func (t *Table) buildData(data [][]interface{}, uid string) ([][]interface{}, error) {
	s := expr.NewStack()
	dataRow := &expr.RowData{
		ColumnIndices: t.allColumnsIdx,
	}
	result := make([][]interface{}, 0, len(data))
	for _, row := range data {
		for i := range row {
			row[i] = normalize.NormalizeType(row[i])
		}
		dataRow.Data = row
		for _, rule := range t.columnRules {
			if rule.sel {
				ret, err := rule.Validate(s, uid, dataRow)
				if err != nil {
					return nil, err
				}
				if ret {
					for _, col := range rule.Columns {
						dataRow.SetInvalid(col, "rule")
					}
				}
			}
		}
		for _, lazyCol := range t.LazyColumns {
			dataRow.SetInvalid(lazyCol, "lazy")
		}
		result = append(result, dataRow.Data)
	}
	return result, nil
}

package table

import (
	"feorm/config"
	"feorm/expr"
	"fmt"
)

type columnRule struct {
	config.ColumnRule
	program *expr.Program
	deny    bool
	insert  bool
	update  bool
	// select
	sel bool
}

func newColumnRule(conf config.ColumnRule) (*columnRule, error) {
	p, err := expr.Compile(conf.Match)
	if err != nil {
		return nil, fmt.Errorf("column rule compile error: %v", err)
	}
	cr := columnRule{
		ColumnRule: conf,
		program:    p,
		deny:       conf.Action == "deny",
		insert:     false,
		update:     false,
		sel:        false,
	}
	for _, o := range conf.Operations {
		switch o {
		case "select":
			cr.sel = true
		case "insert":
			cr.insert = true
		case "update":
			cr.update = true
		default:
			return nil, fmt.Errorf("unsupport operation for column rule: %s", o)
		}
	}
	return &cr, nil
}

func (cr *columnRule) Validate(s *expr.Stack, uid string, data *expr.RowData) (bool, error) {
	ret, err := cr.program.Run(uid, data, s)
	if err != nil {
		// todo 调试时更详细的错误，并且可以在正式运行时忽略错误
		return false, err
	}
	return !expr.Falsy(ret), err
}

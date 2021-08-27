package query

import (
	"fmt"
	"regexp"
	"strings"
)

var whiteList = map[string]bool{
	"+":    true,
	"-":    true,
	"*":    true,
	"/":    true,
	">":    true,
	">=":   true,
	"=":    true,
	"<":    true,
	"<=":   true,
	"!=":   true,
	"and":  true,
	"or":   true,
	"not":  true,
	"like": true,
	"is":   true,
}

type Ast struct {
	// value/bop/uop/call/col/ph
	T string
	V interface{}
	C []Ast
}

type sqlCompileCtx struct {
	currentIndex int
}

func (a *Ast) Compile2Sql(phFunc func(idx int) string) (code string, count int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", recover())
		}
		return
	}()
	ctx := &sqlCompileCtx{currentIndex: 0}
	code, err = a.compileToSqlInternal(ctx, phFunc)
	return code, ctx.currentIndex, err
}

var columnCheckRegexp = regexp.MustCompile(`^\w+$`)

func (a *Ast) compileToSqlInternal(ctx *sqlCompileCtx, phFunc func(int) string) (string, error) {
	switch a.T {
	case "value":
		if a.V == nil {
			return "null", nil
		}
		switch a.V.(type) {
		case string:
			return "'" + strings.ReplaceAll(a.V.(string), "'", "''") + "'", nil
		default:
			return fmt.Sprint(a.V), nil
		}
	case "bop":
		if whiteList[a.V.(string)] == false {
			return "", fmt.Errorf("invalid binary operator %s", a.V)
		}
		if len(a.C) != 2 {
			return "", fmt.Errorf("invalid op number %d, should be 2", len(a.C))
		}
		l, err := a.C[0].compileToSqlInternal(ctx, phFunc)
		if err != nil {
			return "", err
		}
		r, err := a.C[1].compileToSqlInternal(ctx, phFunc)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s %s %s)", l, a.V, r), nil
	case "uop":
		if whiteList[a.V.(string)] == false {
			return "", fmt.Errorf("invalid unary operator %s", a.V)
		}
		if len(a.C) != 1 {
			return "", fmt.Errorf("invalid op number %d, should be 1", len(a.C))
		}
		l, err := a.C[0].compileToSqlInternal(ctx, phFunc)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s %s)", a.V, l), nil
	case "call":
		return "", fmt.Errorf("func call is not supportted now")
	case "col":
		c := a.V.(string)
		if !columnCheckRegexp.MatchString(c) {
			return "", fmt.Errorf("invalid column name '%s'", c)
		}
		return c, nil
	case "ph":
		ph := phFunc(ctx.currentIndex)
		ctx.currentIndex++
		return ph, nil
	default:
		return "", fmt.Errorf("invalid op type %s", a.T)
	}
}

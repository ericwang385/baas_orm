package expr

import (
	"feorm/variable"
	"strconv"
	"strings"
)

type sqlCompileCtx struct {
	currentIndex int
	variables    []string
}

func (ast *AstNode) CompileToSql(paramIndex int, placeHolderFunc func(idx int) string) (string, func(id string) ([]interface{}, error)) {
	// todo 添加缓存
	ctx := &sqlCompileCtx{currentIndex: paramIndex}
	code := ast.compileToSqlInternal(ctx, placeHolderFunc)
	return code, func(id string) ([]interface{}, error) {
		if len(ctx.variables) == 0 {
			return nil, nil
		}
		ret := make([]interface{}, len(ctx.variables))
		for i := range ret {
			var err error
			ret[i], err = variable.GetVariable(ctx.variables[i], id)
			if err != nil {
				return nil, err
			}
		}
		return ret, nil
	}
}

func (ast *AstNode) compileToSqlInternal(ctx *sqlCompileCtx, placeHolderFunc func(idx int) string) string {
	switch ast.Type {
	case ASTBinaryOP:
		return "(" + ast.Children[0].compileToSqlInternal(ctx, placeHolderFunc) + " " +
			ast.Value + " " + ast.Children[1].compileToSqlInternal(ctx, placeHolderFunc) + ")"
	case ASTColumn:
		return ast.Value
	case ASTFuncCall:
		// todo support some functions
		panic("not support function call in sql")
	case ASTUnaryOP:
		return ast.Value + " " + ast.Children[1].compileToSqlInternal(ctx, placeHolderFunc)
	case ASTValueBool:
		v, _ := strconv.ParseBool(ast.Value)
		if v {
			return "true"
		} else {
			return "false"
		}
	case ASTValueInt, ASTValueFloat:
		return ast.Value
	case ASTValueText:
		return "'" + strings.ReplaceAll(ast.Value, "'", "''") + "'"
	case ASTVariable:
		ph := placeHolderFunc(ctx.currentIndex)
		ctx.currentIndex++
		ctx.variables = append(ctx.variables, ast.Value)
		return ph
	}
	panic("")
}

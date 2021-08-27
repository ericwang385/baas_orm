package expr

import (
	"feorm/variable"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

type sqlCompileCtx struct {
	currentIndex int
	variables    []string
}

var cache map[uint32]string

//也可以换成其他的hash函数
func hash(s interface{}) uint32 {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%v", s)))
	return h.Sum32()
}

func (ast *AstNode) CompileToSql(paramIndex int, placeHolderFunc func(idx int) string) (string, func(id string) ([]interface{}, error)) {
	h := hash(ast)
	var code string
	ctx := &sqlCompileCtx{currentIndex: paramIndex}
	if c, err := cache[h]; err {
		code = ast.compileToSqlInternal(ctx, placeHolderFunc)
	} else {
		code = c
	}
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

package expr

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	code := "uid>1"
	l := newLexer(code)
	yyErrorVerbose = true
	//yyDebug=100
	yyParse(l)
	fmt.Println(l.parseResult)
	sql,_:=l.parseResult.CompileToSql(0, func(idx int) string {
		return "?"
	})
	fmt.Println(sql)
}

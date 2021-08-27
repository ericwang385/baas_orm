package expr

import (
	"fmt"
	"regexp"
	"strings"
)

var keywords = map[string]int{
	"+":     ADD,
	"-":     MINUS,
	"*":     MUL,
	"/":     DIV,
	">":     GT,
	">=":    GTE,
	"=":     EQ,
	"<":     LT,
	"<=":    LTE,
	"!=":    NEQ,
	"and":   AND,
	"or":    OR,
	"not":   NOT,
	"like":  LIKE,
	"true":  BOOL,
	"false": BOOL,
	",":     COMMA,
	"(":     LP,
	")":     RP,
	"null":  NULL,
	"@":     AT,
}

type Lexer struct {
	offset      int
	str         string
	length      int
	parseResult *AstNode
}

func newLexer(code string) *Lexer {
	return &Lexer{
		str: code,
	}
}

func (l *Lexer) Text() string {
	return l.str[l.offset-l.length : l.offset]
}

func (l *Lexer) Offset() int {
	return l.offset
}

func match(pattern, str string) int {
	re := regexp.MustCompile(pattern)
	result := re.FindStringSubmatchIndex(str)
	if len(result) > 0 {
		return result[1]
	} else {
		return 0
	}

}

func (l *Lexer) Lex(lval *yySymType) int {
	if l.Offset() == len(l.str) {
		return 0
	}
	var result = -1
	for result <= 0 {
		if l.Offset() >= len(l.str) {
			return 0
		}
		switch {
		case match(`^([ \s\n]+)`, l.str[l.offset:len(l.str)]) > 0:
			result = -1
			l.offset += match(`^[ \s\n]+`, l.str[l.offset:len(l.str)])
		case match(`^(\d+)`, l.str[l.offset:len(l.str)]) > 0:
			result = INT
			l.length = match(`^(\d+)`, l.str[l.offset:len(l.str)])
			l.offset += l.length
		case match(`^(\d+\.\d*)`, l.str[l.offset:len(l.str)]) > 0:
			result = FLOAT
			l.length = match(`^(\d+\.\d*)`, l.str[l.offset:len(l.str)])
			l.offset += l.length
		case match(`^('[^']*')`, l.str[l.offset:len(l.str)]) > 0:
			result = STR
			l.length = match(`^('[^']*')`, l.str[l.offset:len(l.str)])
			l.offset += l.length
		case match(`(?i)^(\+|-|\*|/|>|<|=|!=|<>|>=|<=|and|or|not|@|null|true|false|is\s|,|\(|\)|like)`, l.str[l.offset:len(l.str)]) > 0:
			l.length = match(`(?i)^(\+|-|\*|/|>|<|=|!=|<>|>=|<=|and|or|not|null|@|true|false|is|,|\(|\)|like)`, l.str[l.offset:len(l.str)])
			l.offset += l.length
			result = keywords[strings.ToLower(l.Text())]
		case match(`^(\w+)`, l.str[l.offset:len(l.str)]) > 0:
			l.length = match(`^(\w+)`, l.str[l.offset:len(l.str)])
			result = ID
			l.offset += l.length
		default:
			panic("unexpected token " + l.str[l.offset:len(l.str)])
		}
	}
	lval.offset = l.offset
	lval.text = l.Text()
	return result
}

func (l *Lexer) Error(s string) {
	errInfo := fmt.Sprintf("\n%s\n%s\n", l.str, strings.Repeat(" ", l.Offset())+strings.Repeat("^", len(l.Text())))
	panic(errInfo + s)
}

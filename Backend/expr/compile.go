package expr

import (
	"errors"
	"fmt"
	"strconv"
)

func Compile(code string) (p *Program, err error) {
	if code == "" {
		return nil, nil
	}
	l := newLexer(code)
	yyErrorVerbose = true
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", recover())
		}
		return
	}()
	yyParse(l)
	return l.parseResult.compile(code)
}
func CompileAst(code string) (a *AstNode, err error) {
	if code == "" {
		return nil, nil
	}
	l := newLexer(code)
	yyErrorVerbose = true
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", recover())
		}
		return
	}()
	yyParse(l)
	return l.parseResult, nil
}

func (a *AstNode) compile(code string) (*Program, error) {
	var p = &Program{
		code:    make([]opCode, 0),
		rawCode: code,
	}
	err := a.compileInternal(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (a *AstNode) compileInternal(p *Program) error {
	switch a.Type {
	case ASTValueFloat:
		f, err := strconv.ParseFloat(a.Value, 64)
		if err != nil {
			return err
		}
		p.code = append(p.code, opCode{
			code:     opLoad,
			value:    f,
			position: a.Offset,
		})
	case ASTValueInt:
		i, err := strconv.ParseInt(a.Value, 10, 64)
		if err != nil {
			return err
		}
		p.code = append(p.code, opCode{
			code:     opLoad,
			value:    i,
			position: a.Offset,
		})
	case ASTValueBool:
		b, err := strconv.ParseBool(a.Value)
		if err != nil {
			return err
		}
		p.code = append(p.code, opCode{
			code:     opLoad,
			value:    b,
			position: a.Offset,
		})
	case ASTValueText:
		p.code = append(p.code, opCode{
			code:     opLoad,
			value:    a.Value,
			position: a.Offset,
		})
	case ASTVariable:
		p.code = append(p.code, opCode{
			code:     opVar,
			value:    a.Value,
			position: a.Offset,
		})
	case ASTColumn:
		p.code = append(p.code, opCode{
			code:     opColumn,
			value:    a.Value,
			position: a.Offset,
		})
		p.columns = append(p.columns, a.Value)
	case ASTBinaryOP:
		if err := a.Children[0].compileInternal(p); err != nil {
			return err
		}
		if err := a.Children[1].compileInternal(p); err != nil {
			return err
		}
		p.code = append(p.code, opCode{
			code:     opCall,
			value:    a.Value,
			argc:     2,
			position: a.Offset,
		})
	case ASTUnaryOP:
		if err := a.Children[0].compileInternal(p); err != nil {
			return err
		}
		p.code = append(p.code, opCode{
			code:     opCall,
			value:    a.Value,
			argc:     1,
			position: a.Offset,
		})
	case ASTFuncCall:
		return errors.New("not support function now")
	default:
		return errors.New("invalid ast node")
	}
	return nil
}

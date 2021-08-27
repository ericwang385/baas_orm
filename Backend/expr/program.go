package expr

import (
	"feorm/variable"
	"fmt"
	"strings"
	"sync"
)

type op int

const (
	opLoad = iota
	opCall
	opVar
	opColumn
)

type opCode struct {
	code     op
	value    interface{}
	argc     int
	position int
}

var stackPool = sync.Pool{New: func() interface{} {
	return &Stack{
		s: make([]interface{}, 4),
	}
}}

func NewStack() *Stack {
	return stackPool.Get().(*Stack)
}

type Stack struct {
	s []interface{}
}

func (s *Stack) Close() {
	s.s = s.s[0:0]
	stackPool.Put(s)
}
func (s *Stack) Push(item interface{}) {
	s.s = append(s.s, item)
}

func (s *Stack) Pop() interface{} {
	l := len(s.s) - 1
	item := s.s[l]
	s.s = s.s[0:l]
	return item
}
func (s *Stack) PopN(n int) []interface{} {
	// 不用处理n>len(s.s)，不应该发生
	l := len(s.s) - n
	items := s.s[l:len(s.s)]
	s.s = s.s[0:l]
	return items
}

type Program struct {
	code    []opCode
	rawCode string
	columns []string
}

func (p *Program) error(err interface{}, pos int) error {
	return fmt.Errorf("%s\n%s^\n%v", p.rawCode, strings.Repeat(" ", pos), err)
}

func (p *Program) Run(uid string, data *RowData, stack *Stack) (interface{}, error) {
	if stack == nil {
		stack = NewStack()
	}
	defer stack.Close()
	for _, op := range p.code {
		switch op.code {
		case opLoad:
			stack.Push(op.value)
		case opVar:
			v, err := variable.GetVariable(op.value.(string), uid)
			if err != nil {
				return nil, p.error(err, op.position)
			}
			stack.Push(v)
		case opColumn:
			stack.Push(data.Get(op.value.(string)))
		case opCall:
			f := functionMap[op.value.(string)]
			if f == nil {
				return nil, p.error(fmt.Sprintf("function %s 未定义", op.value.(string)), op.position)
			}
			args := stack.PopN(op.argc)
			ret, err := f(args)
			if err != nil {
				return nil, p.error(err, op.position)
			}
			stack.Push(ret)
		}
	}
	return stack.Pop(), nil
}

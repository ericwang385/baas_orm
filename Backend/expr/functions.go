package expr

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

func buildBinaryTypeError(op string, v1, v2 interface{}) error {
	return fmt.Errorf("无法对 %s 和 %s 执行 %s 操作", reflect.TypeOf(v1).Name(), reflect.TypeOf(v2).Name(), op)
}

func Falsy(v interface{}) bool {
	if v == nil {
		return true
	}
	switch vv := v.(type) {
	case int64:
		return vv == 0
	case float64:
		return vv <= math.SmallestNonzeroFloat64
	case string:
		return vv == ""
	case bool:
		return !vv
	case time.Time:
		return true
	}
	panic("not support falsy of type " + reflect.TypeOf(v).Name())
}

var functionMap = map[string]func(args []interface{}) (interface{}, error){
	"+": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv + rv, nil
			case float64:
				return float64(lv) + rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case int64:
				return lv + float64(rv), nil
			case float64:
				return lv + rv, nil
			}
		case string:
			switch rv := r.(type) {
			case string:
				return lv + rv, nil

			}
		}
		return nil, buildBinaryTypeError("+", l, r)
	},
	"-": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv - rv, nil
			case float64:
				return float64(lv) - rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case int64:
				return lv - float64(rv), nil
			case float64:
				return lv - rv, nil
			}
		}
		return nil, buildBinaryTypeError("-", l, r)
	},
	"*": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv * rv, nil
			case float64:
				return float64(lv) * rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case int64:
				return lv * float64(rv), nil
			case float64:
				return lv * rv, nil
			}
		}
		return nil, buildBinaryTypeError("*", l, r)
	}, "/": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				if rv == 0 {
					return nil, errors.New("divide zero")
				}
				return lv / rv, nil
			case float64:
				return float64(lv) / rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case int64:
				return lv / float64(rv), nil
			case float64:
				return lv / rv, nil
			}
		}
		return nil, buildBinaryTypeError("/", l, r)
	},
	"=": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv == rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case float64:
				return lv == rv, nil
			}
		case string:
			switch rv := r.(type) {
			case string:
				return lv == rv, nil
			}
		case bool:
			switch rv := r.(type) {
			case bool:
				return lv == rv, nil
			}
		case time.Time:
			switch rv := r.(type) {
			case time.Time:
				return lv.Equal(rv), nil
			}
		}
		return nil, buildBinaryTypeError("=", l, r)
	},
	">": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv > rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case float64:
				return lv > rv, nil
			}
		case string:
			switch rv := r.(type) {
			case string:
				return lv > rv, nil
			}
		case time.Time:
			switch rv := r.(type) {
			case time.Time:
				return lv.After(rv), nil
			}
		}
		return nil, buildBinaryTypeError(">", l, r)
	},
	">=": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv >= rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case float64:
				return lv >= rv, nil
			}
		case string:
			switch rv := r.(type) {
			case string:
				return lv >= rv, nil
			}
		case time.Time:
			switch rv := r.(type) {
			case time.Time:
				return lv.After(rv) || lv.Equal(rv), nil
			}
		}
		return nil, buildBinaryTypeError(">=", l, r)
	},
	"<=": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv <= rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case float64:
				return lv <= rv, nil
			}
		case string:
			switch rv := r.(type) {
			case string:
				return lv <= rv, nil
			}
		case time.Time:
			switch rv := r.(type) {
			case time.Time:
				return lv.Before(rv) || lv.Equal(rv), nil
			}
		}
		return nil, buildBinaryTypeError("<=", l, r)
	},
	"<": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if l == nil || r == nil {
			return nil, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv < rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case float64:
				return lv < rv, nil
			}
		case string:
			switch rv := r.(type) {
			case string:
				return lv < rv, nil
			}
		case time.Time:
			switch rv := r.(type) {
			case time.Time:
				return lv.Before(rv), nil
			}
		}
		return nil, buildBinaryTypeError("<", l, r)
	}, "!=": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		if (l == nil && r != nil) || (l != nil && r == nil) {
			return true, nil
		}
		if l == nil && r == nil {
			return false, nil
		}
		switch lv := l.(type) {
		case int64:
			switch rv := r.(type) {
			case int64:
				return lv != rv, nil
			}
		case float64:
			switch rv := r.(type) {
			case float64:
				return lv != rv, nil
			}
		case string:
			switch rv := r.(type) {
			case string:
				return lv != rv, nil
			}
		case time.Time:
			switch rv := r.(type) {
			case time.Time:
				return !lv.Equal(rv), nil
			}
		}
		return nil, buildBinaryTypeError("!=", l, r)
	},
	"and": func(args []interface{}) (interface{}, error) {
		l := args[0]
		r := args[1]
		return (!Falsy(l)) && (!Falsy(r)), nil
	},
}

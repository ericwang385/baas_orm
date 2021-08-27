package normalize

import (
	"reflect"
	"time"
)

func NormalizeType(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case string, int64, float64, bool, time.Time:
		return v
	default:
		panic("please support normalize " + reflect.TypeOf(v).Name())
	}
}

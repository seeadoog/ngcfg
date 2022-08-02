//go:build go1.18

package values

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValueTypes interface {
	string | int | int64 | bool | float64
}

type TagValue[T ValueTypes] struct {
	val     T
	options map[string]string
}

func (ov *TagValue[T]) UnmarshalCfg(path string, val interface{}) error {
	switch v := val.(type) {
	case []string:
		if len(v) == 0 {
			return nil
		}
		val := v[0]
		vv, err := valueOfOptionValue(ov.val, val)
		if err != nil {
			return fmt.Errorf("'%s' %v", path, err)
		}
		ov.val = vv.(T)

		ov.options = map[string]string{}
		for _, s := range v[1:] {
			name, val, ok := strings.Cut(s, "=")
			if ok {
				ov.options[name] = val
			} else {
				ov.options[name] = ""
			}
		}
		return nil
	default:
		return nil
	}
}

func (ov *TagValue[T]) Val() T {
	return ov.val
}

func (ov *TagValue[T]) GetTag(key string) string {
	return ov.options[key]
}

func (ov *TagValue[T]) GetTagOf(key string) (v string, ok bool) {
	v, ok = ov.options[key]
	return
}

func (ov *TagValue[T]) GetBool(key string) bool {
	v, ok := ov.options[key]
	if ok {
		switch v {
		case "", "true", "1", "True", "TRUE":
			return true
		default:
			return false
		}
	}
	return false
}

func (ov *TagValue[T]) String() string {
	sb := strings.Builder{}
	for key, val := range ov.options {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(val)
		sb.WriteString(" ")
	}
	return fmt.Sprintf("%v %v", ov.val, sb.String())
}

func (ov *TagValue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(ov.String())
}

func valueOfOptionValue(vt interface{}, vv string) (interface{}, error) {
	switch vt.(type) {
	case string:
		return vv, nil
	case int:
		return strconv.Atoi(vv)
	case int64:
		vvv, err := strconv.Atoi(vv)
		return int64(vvv), err
	case float64:
		return strconv.ParseFloat(vv, 64)
	case bool:
		return strconv.ParseBool(vv)
	default:
		return nil, fmt.Errorf("unsupport type of TagValue :%v", reflect.TypeOf(vt))
	}
}

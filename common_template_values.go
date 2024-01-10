package ngcfg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Set struct {
	data map[string]bool
}

func (s *Set) UnmarshalCfg(path string, val interface{}) error {
	arr, ok := val.([]string)
	if !ok {
		return fmt.Errorf("%s decode value to Set err, type must be []string", path)
	}
	s.data = map[string]bool{}
	for _, v := range arr {
		s.data[v] = true
	}
	return nil
}
func (s *Set) Contains(key string) bool {
	if s == nil {
		return false
	}
	return s.data[key]
}

func (s *Set) Values() map[string]bool {
	return s.data
}

func (s *Set) String() string {
	sb := strings.Builder{}

	for k, _ := range s.data {
		sb.WriteString(k)
		sb.WriteString(";")
	}
	return sb.String()
}

/*
arssd nam=5 age=8 sname=lixiang
*/
type Options[O any] struct {
	opt *O
	raw []string
}

func (s *Options[O]) Opt() *O {
	return s.opt
}

func (s *Options[O]) UnmarshalCfg(path string, val interface{}) error {
	arr, ok := val.(BasicValue)
	if !ok {
		return fmt.Errorf("%s decode value to Options err, type must be []string", path)
	}
	s.raw = arr
	e := NewElem()
	for _, v := range arr {
		key, val, ok := strings.Cut(v, "=")
		if !ok {
			e.Set(key, BasicValue{})
		} else {
			e.Set(key, BasicValue{val})
		}
	}
	opt := new(O)
	err := e.Decode(opt)
	if err != nil {
		return fmt.Errorf("decode opt error :'%s' %w", path, err)
	}
	s.opt = opt
	return nil
}

func (s *Options[O]) String() string {
	return strings.Join(s.raw, " ")
}

type ValueTypes any

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

func (ov *TagValue[T]) LookUpTag(key string) (v string, ok bool) {
	v, ok = ov.options[key]
	return
}

func (ov *TagValue[T]) GetTagBool(key string) bool {
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

func (ov *TagValue[T]) GetTagInt(key string) (int, error) {
	v := ov.options[key]
	if v == "" {
		return 0, fmt.Errorf("'%s' empty value", key)
	}
	va, err := strconv.Atoi(key)
	if err != nil {
		return 0, fmt.Errorf("get tag as  int error , key %s ,val %v err:%w", key, v, err)
	}
	return va, nil
}

func (ov *TagValue[T]) GetTagFloat(key string) (float64, error) {
	v := ov.options[key]
	if v == "" {
		return 0, fmt.Errorf("'%s' empty value", key)
	}
	va, err := strconv.ParseFloat(key, 64)
	if err != nil {
		return 0, fmt.Errorf("get tag as  float  error , key %s ,val %v err:%w", key, v, err)
	}
	return va, nil
}

func (ov *TagValue[T]) GetTagUint(key string) (uint, error) {
	v := ov.options[key]
	if v == "" {
		return 0, fmt.Errorf("'%s' empty value", key)
	}
	va, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("get tag as  uint  error , key %s ,val %v err:%w", key, v, err)
	}
	return uint(va), nil
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

type TagValueT[Val any, Tag any] struct {
	values *TagValue[Val]
	tag    *Tag
}

func (t *TagValueT[Val, Tag]) Val() Val {
	return t.values.val
}

func (t *TagValueT[Val, Tag]) Tag() *Tag {
	return t.tag
}

func (t *TagValueT[Val, Tag]) UnmarshalCfg(path string, v any) error {
	t.values = new(TagValue[Val])
	err := t.values.UnmarshalCfg(path, v)
	if err != nil {
		return err
	}

	e := NewElem()

	for key, v2 := range t.values.options {
		switch v2 {
		case "":
			e.Set(key, []string{})
		default:
			e.Set(key, []string{v2})
		}

	}
	t.tag = new(Tag)
	return e.Decode(t.tag)
}

func (t *TagValueT[Val, Tag]) String() string {
	return t.values.String()
}

func (t *TagValueT[Val, Tag]) MarshalJSON() (b []byte, err error) {
	return t.values.MarshalJSON()
}

// expand 19.22.34.{12...13} to [19.22.34.12 19.22.34.13]
type ConsecutiveString struct {
	ips []string
}

func (c *ConsecutiveString) UnmarshalCfg(path string, val interface{}) error {
	switch v := val.(type) {
	case []string:
		ips, err := ParseRangeIpss(v)
		if err != nil {
			return err
		}
		c.ips = ips
		return nil
	case string:
		ips, err := ParseRangeIps(v)
		if err != nil {
			return err
		}
		c.ips = ips
		return nil
	default:
		return fmt.Errorf("%s :unsupport valtype:%v", path, reflect.TypeOf(val))
	}
}

func (c *ConsecutiveString) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ips)
}

func (c *ConsecutiveString) String() string {
	if c == nil {
		return ""
	}
	return "[ " + strings.Join(c.ips, " ") + "]"
}

func (c *ConsecutiveString) Strings() []string {
	if c == nil {
		return nil
	}
	return c.ips
}

// 172.21.157.[25-30]
// 172.21.157.[25,40,45]
type RangeStringArray struct {
	ips []string
}

func (c *RangeStringArray) UnmarshalCfg(path string, val interface{}) error {
	switch v := val.(type) {
	case []string:
		ips, err := parseRangeIps2(v)
		if err != nil {
			return err
		}
		c.ips = ips
		return nil
	case string:
		ips, err := parseRangeIp2(v)
		if err != nil {
			return err
		}
		c.ips = ips
		return nil
	default:
		return fmt.Errorf("%s :unsupport valtype:%v", path, reflect.TypeOf(val))
	}
}

func (c *RangeStringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ips)
}

func (c *RangeStringArray) String() string {
	if c == nil {
		return ""
	}
	return "[ " + strings.Join(c.ips, " ") + "]"
}

func (c *RangeStringArray) Strings() []string {
	if c == nil {
		return nil
	}
	return c.ips
}

var (
	rangedStringType = reflect.TypeOf(RangedString{})
)

type RangedString []string

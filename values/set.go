package values

import (
	"fmt"
	"strings"

	"github.com/seeadoog/ngcfg"
)

// Deprecated
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

// Deprecated
type Options[O any] struct {
	opt *O
	raw []string
}

func (s *Options[O]) Opt() *O {
	return s.opt
}

func (s *Options[O]) UnmarshalCfg(path string, val interface{}) error {
	arr, ok := val.(ngcfg.BasicValue)
	if !ok {
		return fmt.Errorf("%s decode value to Set err, type must be *elem", path)
	}
	s.raw = arr
	e := ngcfg.NewElem()
	for _, v := range arr {
		key, val, ok := strings.Cut(v, "=")
		if !ok {
			e.Set(key, ngcfg.BasicValue{})
		} else {
			e.Set(key, ngcfg.BasicValue{val})
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

package ngcfg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/seeadoog/linkedMap"
)

// linked map
type LSMap[V any] struct {
	l *linkedMap.Map[string, V]
}

func NewLSMap[V any]() *LSMap[V] {
	return &LSMap[V]{l: linkedMap.New[string, V]()}
}

func (l *LSMap[V]) Set(k string, v V) {
	l.l.Set(k, v)
}

func (l *LSMap[V]) Get(k string) (v V, ok bool) {
	return l.l.Get(k)
}

func (l *LSMap[V]) Delete(k string) {
	l.l.Delete(k)
}

func (l *LSMap[V]) Range(f func(k string, v V) bool) {
	l.l.Range(f)
}

func (l *LSMap[V]) String() string {
	if l == nil {
		return "<nil>"
	}
	bf := &strings.Builder{}
	bf.WriteString("[")
	l.Range(func(k string, v V) bool {
		bf.WriteString(k)
		bf.WriteString(":")
		bf.WriteString(fmt.Sprintf("%v", v))
		bf.WriteString("; ")
		return true
	})
	bf.WriteString("]")
	return bf.String()
}

func (l *LSMap[V]) UnmarshalCfg(path string, val interface{}) error {
	l.l = linkedMap.New[string, V]()
	switch v := val.(type) {
	case *Elem:
		it := v.Iterator()
		for it.HasNext() {
			e := it.Next()
			data := new(V)
			v := reflect.ValueOf(data)
			err := unmarshalObject2Struct(path+"."+e.Key, e.Val, v, false)
			if err != nil {
				return fmt.Errorf("%s %w", path, err)
			}
			l.Set(e.Key, *data)
		}
		return nil
	case nil:
		return nil
	}
	return fmt.Errorf("cannot unmarshal %s to LSMap", reflect.TypeOf(val))
}

func (l *LSMap[V]) MarshalJSON() ([]byte, error) {
	return l.l.MarshalJSON()
}

func (l *LSMap[V]) UnmarshalJSON(b []byte) error {
	l.l = linkedMap.New[string, V]()
	return l.l.UnmarshalJSON(b)
}

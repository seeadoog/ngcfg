package ngcfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Elem struct {
	data   *LinkedMap
	idx    int
	parent *Elem
}

func (e *Elem) getIdx() string {
	sidx := "index_" + strconv.Itoa(e.idx)
	e.idx++
	return sidx
}

func (e *Elem) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.data)
}

func (e *Elem) Set(k string, v interface{}) error {
	_, ok := e.data.Get(k)
	if ok {
		return fmt.Errorf("%s key has already defined", k)
	}
	//e.data[k] = v
	ele, ok := v.(*Elem)
	if ok {
		ele.parent = e
	}
	if k == "-" {
		e.data.Set(e.getIdx(), v)
		return nil
	}
	e.data.Set(k, v)

	return nil
}

func (e *Elem) GetCtxString(key string) (string, error) {
	_, ok := e.data.Get(key)
	if !ok {
		if e.parent != nil {
			return e.parent.GetCtxString(key)
		}
	}
	return e.GetString(key)
}

func (e *Elem) GetCtxStringDef(key, def string) string {
	val, err := e.GetCtxString(key)
	if err != nil {
		return def
	}
	return val
}

func (e *Elem) GetCtxArray(key string) ([]string, error) {
	_, ok := e.data.Get(key)
	if !ok {
		if e.parent != nil {
			return e.parent.GetCtxArray(key)
		}
	}
	return e.GetArray(key)
}

func (e *Elem) GetCtxBool(key string) (bool, error) {
	_, ok := e.data.Get(key)
	if !ok {
		if e.parent != nil {
			return e.parent.GetCtxBool(key)
		}
	}
	return e.GetBool(key)
}

func (e *Elem) GetCtxBoolDef(key string, def bool) bool {
	val, err := e.GetCtxBool(key)
	if err != nil {
		return def
	}
	return val
}

func (e *Elem) GetCtxNumber(key string) (float64, error) {
	_, ok := e.data.Get(key)
	if !ok {
		if e.parent != nil {
			return e.parent.GetCtxNumber(key)
		}
	}
	return e.GetNumber(key)
}

func (e *Elem) GetCtxNumberDef(key string, def float64) float64 {
	val, err := e.GetCtxNumber(key)
	if err != nil {
		return def
	}
	return val
}

func (e *Elem) GetCtxInt(key string) (int, error) {
	_, ok := e.data.Get(key)
	if !ok {
		if e.parent != nil {
			return e.parent.GetCtxInt(key)
		}
	}
	return e.GetInt(key)
}

func (e *Elem) GetCtxIntDef(key string, def int) int {
	val, err := e.GetCtxInt(key)
	if err != nil {
		return def
	}
	return val
}

func (e *Elem) GetCtxElem(key string) (*Elem, error) {
	_, ok := e.data.Get(key)
	if !ok {
		if e.parent != nil {
			return e.parent.GetCtxElem(key)
		}
	}
	return e.GetElem(key)
}

func (e *Elem) SetSub(key, sub string, v interface{}) error {
	ele, ok := e.data.Get(key)
	if !ok {
		ele = NewElem()
		if err := e.Set(key, ele); err != nil {
			return err
		}
	}
	eleO, ok := ele.(*Elem)
	if !ok {
		return fmt.Errorf("set sub key failed:parent key %s is not object", key)
	}
	return eleO.Set(sub, v)
}

func (e *Elem) RawMap() *LinkedMap {
	return e.data
}

func (e *Elem) Iterator() Iterator {
	return e.data.Iterator()
}

func (e *Elem) Get(key string) interface{} {
	v, _ := e.data.Get(key)
	return v
}

func (e *Elem) GetCtx(key string) interface{} {
	v, ok := e.data.Get(key)
	if !ok {
		if e.parent != nil {
			return e.parent.GetCtx(key)
		}
	}
	return v
}

func NewElem() *Elem {
	return &Elem{data: NewLinkedMap()}
}

//return fist elem of array
func (e *Elem) GetString(key string) (string, error) {
	v, ok := e.data.Get(key)
	if !ok {
		return "", fmt.Errorf("key %s doesn't exists", key)
	}
	switch v.(type) {
	case []string:
		ss := v.([]string)
		if len(ss) > 0 {
			return ss[0], nil
		}
		return "", nil
	}
	return "", fmt.Errorf("type of %s is object", key)
}

func (e *Elem) GetNumber(key string) (float64, error) {
	s, err := e.GetString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

func (e *Elem) GetInt(key string) (int, error) {
	f, err := e.GetNumber(key)
	if err != nil {
		return 0, err
	}
	i := int(f)

	if float64(i) != f {
		return i, fmt.Errorf("value of %s is float not int:%v", key, f)
	}
	return i, nil
}

func (e *Elem) GetBool(key string) (bool, error) {
	s, err := e.GetString(key)
	if err != nil {
		return false, err
	}
	bv, err := boolOf(s)
	if err != nil {
		return false, fmt.Errorf("key:%s,%w", key, err)
	}
	return bv, nil
}

func (e *Elem) GetArray(key string) ([]string, error) {
	arr, ok := e.data.Get(key)
	if !ok {
		return nil, fmt.Errorf("key %s doesn't exists", key)
	}
	switch arr.(type) {
	case []string:
		return arr.([]string), nil
	}
	return nil, fmt.Errorf("type of %s is not array", key)
}

func (e *Elem) GetStringDef(key string, def string) string {
	s, err := e.GetString(key)
	if err != nil {
		return def
	}
	return s
}

func (e *Elem) GetNumberDef(key string, def float64) float64 {
	v, err := e.GetNumber(key)
	if err != nil {
		return def
	}
	return v
}

func (e *Elem) GetIntDef(key string, def int) int {
	v, err := e.GetInt(key)
	if err != nil {
		return def
	}
	return v
}

func (e *Elem) GetBoolDef(key string, def bool) bool {
	v, err := e.GetBool(key)
	if err != nil {
		return def
	}
	return v
}

func (e *Elem) GetArrayDef(key string, def []string) []string {
	v, err := e.GetArray(key)
	if err != nil {
		return def
	}
	return v
}

func (e *Elem) Elem(key string) *Elem {
	if e == nil {
		return e
	}
	res, _ := e.GetElem(key)
	return res
}

func (e *Elem) GetElem(key string) (*Elem, error) {
	v, ok := e.data.Get(key)
	if !ok {
		return nil, fmt.Errorf("key %s does not exist", key)
	}
	switch v.(type) {
	case *Elem:
		return v.(*Elem), nil
	}
	return nil, fmt.Errorf("type of %s is not elem", key)
}

func (e *Elem) AsStringArray() ([]string, error) {
	it := e.Iterator()
	res := make([]string, 0)
	for it.HasNext() {
		elem := it.Next()
		switch elem.Val.(type) {
		case []string:
			res = append(res, elem.Val.([]string)...)
		default:
			return nil, fmt.Errorf("type is not array:%s", reflect.TypeOf(elem.Val).String())
		}
	}
	return res, nil
}

func (e *Elem) AsArray() ([][]string, error) {
	it := e.Iterator()
	res := make([][]string, 0)
	for it.HasNext() {
		elem := it.Next()
		switch elem.Val.(type) {
		case []string:
			res = append(res, elem.Val.([]string))
		default:
			return nil, fmt.Errorf("type is not array:%s", reflect.TypeOf(elem.Val).String())
		}
	}
	return res, nil
}

func (e *Elem) MarshalCfg(n int) ([]byte, error) {
	bf := bytes.Buffer{}
	it := e.Iterator()
	for it.HasNext() {
		o := it.Next()
		om, ok := o.Val.(Marshaller)
		if ok {
			b, err := om.MarshalCfg(n + 1)
			if err != nil {
				return nil, err
			}
			bf.WriteString(nt(n))
			bf.WriteString(o.Key)
			bf.WriteString(" {\r\n")
			bf.Write(b)
			bf.WriteString(nt(n))
			bf.WriteString("}\r\n")
		} else {
			bf.WriteString(fmt.Sprintf("%s%s %s\r\n", nt(n), o.Key, vts(o.Val)))
		}
	}
	return bf.Bytes(), nil
}

func vts(i interface{}) string {
	b, _ := i.([]string)
	sb := strings.Builder{}
	for _, s := range b {
		sb.WriteString(fmt.Sprintf("%s ", s))
	}
	return sb.String()
}

func nt(n int) string {
	b := strings.Builder{}
	for i := 0; i < n; i++ {
		b.WriteString("   ")
	}
	return b.String()
}

type Marshaller interface {
	MarshalCfg(n int) ([]byte, error)
}

func boolOf(s string) (bool, error) {
	switch s {
	case "true", "1", "", "on", "yes", "ok":
		return true, nil
	case "false", "0", "off", "no", "never":
		return false, nil
	}
	return false, fmt.Errorf("invalid bool value of %s", s)
}

//func (e *Elem)ToArray()[]interface{}{
//	res:=make([]interface{}, 0,len(e.data.data))
//	it:=e.RawMap().MapItem()
//	for it != nil{
//		res = append(res,it.Val)
//	}
//}

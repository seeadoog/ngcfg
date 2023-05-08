package ngcfg

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// UnmarshalFromMap 将map 中的值序列化到 struct 中
func UnmarshalFromElem(in *Elem, template interface{}) error {
	v := reflect.ValueOf(template)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		panic("template value is nil or not pointer")
	}
	return unmarshalObject2Struct("", in, v, false)
}

func UnmarshalFromFile(file string, tlp interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return UnmarshalFromBytes(data, tlp)
}

var (
	elemType = reflect.TypeOf(&Elem{})
)

type Unmarshaller interface {
	UnmarshalCfg(path string, val interface{}) error
}

var (
	unmarshalType = reflect.TypeOf(new(Unmarshaller)).Elem()

	customTagValidator = map[string]func(path string, v interface{}, tagV string) error{}
)

func RegisterCustomTagValidator(tagname string, f func(path string, v interface{}, tagV string) error) {
	customTagValidator[tagname] = f
}

type BinarySize int64

var (
	durationType = reflect.TypeOf(time.Duration(0))

	bytesSizeType = reflect.TypeOf(BinarySize(0))
)

func unmarshalObject2Struct(path string, in interface{}, v reflect.Value, usectx bool) error {
	if in == nil {
		return nil
	}

	if v.Kind() != reflect.Ptr && !v.CanSet() {
		return nil
	}

	if v.Type() == elemType {
		inv := reflect.ValueOf(in)
		if inv.Type() == elemType {
			v.Set(inv)
			return nil
		}
		return fmt.Errorf("%s is not *Elem ,but :%v", path, inv.Type())
	}
	if v.Type().Implements(unmarshalType) {
		switch v.Kind() {
		case reflect.Ptr:
			if v.IsNil() {
				elemType := v.Type().Elem()
				newV := reflect.New(elemType)
				err := (newV.Interface().(Unmarshaller)).UnmarshalCfg(path, in)
				if err != nil {
					return fmt.Errorf("'%s' %w", path, err)
				}
				v.Set(newV)
				return nil
			}
			fallthrough
		default:
			err := v.Interface().(Unmarshaller).UnmarshalCfg(path, in)
			if err != nil {
				return fmt.Errorf("'%s' %w", path, err)
			}
			return nil
		}
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			vt := v.Type()
			elemType := vt.Elem()
			var nv reflect.Value
			switch elemType.Kind() {
			default:
				nv = reflect.New(elemType)
			}
			err := unmarshalObject2Struct(path, in, nv.Elem(), usectx)
			if err != nil {
				return err
			}
			v.Set(nv)
			return nil
		}
		return unmarshalObject2Struct(path, in, v.Elem(), usectx)
	case reflect.Slice:
		//arr, ok := in.(*Elem)
		t := v.Type()

		switch arr := in.(type) {
		case []string:
			elemType := t.Elem()
			slice := reflect.MakeSlice(t, 0, len(arr))
			for idx, strv := range arr {
				elemVal := reflect.New(elemType)
				err := unmarshalObject2Struct(path+fmt.Sprintf("[%v]", idx), strv, elemVal, usectx)
				if err != nil {
					return err
				}
				slice = reflect.Append(slice, elemVal.Elem())
			}
			v.Set(slice)
		case *Elem:
			elemType := t.Elem()
			slice := reflect.MakeSlice(t, 0, arr.data.Len())

			it := arr.Iterator()
			for it.HasNext() {
				e := it.Next()
				elemVal := reflect.New(elemType)
				err := unmarshalObject2Struct(path+"."+e.Key, e.Val, elemVal, usectx)
				if err != nil {
					return err
				}
				slice = reflect.Append(slice, elemVal.Elem())
			}
			v.Set(slice)
		default:
			return fmt.Errorf("type of %s should be slice, but:%v", path, t)
		}

		//for _, v := range arr {
		//	elemVal := reflect.New(elemType)
		//	err := unmarshalObject2Struct(path, v, elemVal)
		//	if err != nil {
		//		return err
		//	}
		//	slice = reflect.Append(slice, elemVal.Elem())
		//}
		return nil
	case reflect.String:
		vv, err := StringValueOf(in)
		if err != nil {
			return fmt.Errorf("type of %s should be string, but:%v", path, err)
		}
		v.SetString(vv)
		return nil
	case reflect.Map:
		vmap, ok := in.(*Elem)
		if !ok {
			return fmt.Errorf("type of %s should be object, but %v", path, reflect.TypeOf(in))
		}
		t := v.Type()
		elemT := t.Elem()
		newV := v
		if v.IsNil() {
			newV = reflect.MakeMap(v.Type())
		}
		it := vmap.Iterator()
		for it.HasNext() {
			e := it.Next()
			key := e.Key
			val := e.Val
			elemV := reflect.New(elemT)
			err := unmarshalObject2Struct(path+"."+key, val, elemV, usectx)
			if err != nil {
				return err
			}
			newV.SetMapIndex(reflect.ValueOf(key), elemV.Elem())
		}
		v.Set(newV)
		return nil
	case reflect.Struct:
		t := v.Type()

		vmap, ok := in.(*Elem)
		if !ok {
			return fmt.Errorf("type of %s should be object", path)
		}
		for i := 0; i < t.NumField(); i++ {
			fieldT := t.Field(i)
			name, _, _ := strings.Cut(fieldT.Tag.Get("json"), ",")
			if name == "" {
				name = fieldT.Name
			}
			if fieldT.Anonymous {
				err := unmarshalObject2Struct(path+"."+name, in, v.Field(i), usectx)
				if err != nil {
					return err
				}
				continue
			}
			var elemV interface{}
			if usectx {
				elemV = vmap.GetCtx(name)
			} else {
				elemV = vmap.Get(name)
			}
			var err error
			if elemV == nil {

				envKey := fieldT.Tag.Get("env")
				if envKey != "" {
					envVal := os.Getenv(envKey)
					if envVal != "" {
						err = unmarshalObject2Struct(path+"."+name, envVal, v.Field(i), usectx)
						if err != nil {
							return err
						}
						continue
					}
				}

				def := fieldT.Tag.Get("default")
				if def != "" {
					err = unmarshalObject2Struct(path+"."+name, def, v.Field(i), usectx)
					if err != nil {
						return err
					}
					continue
				}
				if isTrue(fieldT.Tag.Get("required")) {
					return fmt.Errorf("miss field ,'%s' is required", path+"."+name)
				}
				continue
			}

			for key, f := range customTagValidator {
				cv := fieldT.Tag.Get(key)
				if cv != "" {
					err = f(path+"."+name, elemV, cv)
					if err != nil {
						return fmt.Errorf("'%s' %v", path+"."+name, err)
					}
				}
			}
			err = unmarshalObject2Struct(path+"."+name, elemV, v.Field(i), usectx)
			if err != nil {
				return err
			}

		}
		return nil
	case reflect.Interface:
		inVal := reflect.ValueOf(in)
		if inVal.Type().Implements(v.Type()) {
			v.Set(inVal)
		}
		return nil
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:

		//时间类型
		if v.Type() == durationType {
			val, err := StringValueOf(in)
			if err != nil {
				return fmt.Errorf("cannot parse '%s' as duration value: '%v' err:%w", path, in, err)
			}

			d, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("cannot parse '%s' as duration value: '%v' err:%w", path, in, err)
			}
			v.SetInt(int64(d))
			return nil
		}

		if v.Type() == bytesSizeType {
			val, err := StringValueOf(in)
			if err != nil {
				return fmt.Errorf("cannot parse '%s' as BinarySize value: '%v' err:%w", path, in, err)
			}

			d, err := parseByteSize(val)
			if err != nil {
				return fmt.Errorf("cannot parse '%s' as BinarySize value: '%v' err:%w", path, in, err)
			}
			v.SetInt(int64(d))

			return nil
		}

		intV, err := intValueOf(in)
		if err != nil {
			return err
		}
		v.SetInt(intV)
		return nil
	case reflect.Bool:
		boolV, err := boolValueOf(in)
		if err != nil {
			return fmt.Errorf("%s error:%w", path, err)
		}
		v.SetBool(boolV)
		return nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		intV, err := intValueOf(in)
		if err != nil {
			return err
		}
		v.SetUint(uint64(intV))
		return nil
	case reflect.Float64, reflect.Float32:
		floatV, err := floatValueOf(in)
		if err != nil {
			return err
		}
		v.SetFloat(floatV)
		return nil
	//case reflect.Array:
	//
	//	arr, ok := in.([]interface{})
	//	//t := v.Type()
	//	if !ok {
	//		return fmt.Errorf("type of %s should be slice", path)
	//	}
	//
	//	arType := reflect.ArrayOf(v.Len(), v.Type().Elem())
	//	arrv := reflect.New(arType)
	//	pointer := arrv.Pointer()
	//	eleSize := v.Type().Elem().Size()
	//	if v.Len() < len(arr) {
	//		return fmt.Errorf("length of %s is %d . but target value length is %d", path, v.Len(), len(arr))
	//	}
	//	for i, vv := range arr {
	//		elemV := reflect.New(v.Type().Elem())
	//		err := unmarshalObject2Struct(path, vv, elemV)
	//		if err != nil {
	//			return err
	//		}
	//		memCopy(pointer+uintptr(i)*eleSize, elemV.Pointer(), eleSize)
	//	}
	//	v.Set(arrv.Elem())
	default:
		return fmt.Errorf("'%s' not support :%v", path, v.Kind().String())
	}
	//return nil
}

func intValueOf(v interface{}) (int64, error) {
	strv, err := StringValueOf(v)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strv, 10, 64)
}

func boolValueOf(v interface{}) (bool, error) {
	strv, err := StringValueOf(v)
	if err != nil {
		return false, err
	}
	switch strv {
	case "on", "true", "1", "yes", "ok", "":
		return true, nil
	case "off", "false", "0", "never", "no":
		return false, nil
	default:
		return strconv.ParseBool(strv)
	}
}

func floatValueOf(v interface{}) (float64, error) {
	strv, err := StringValueOf(v)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strv, 64)
}

func bytesOf(p uintptr, len uintptr) []byte {
	h := &reflect.SliceHeader{
		Data: p,
		Len:  int(len),
		Cap:  int(len),
	}
	return *(*[]byte)(unsafe.Pointer(h))
}

func memCopy(dst, src uintptr, len uintptr) {
	db := bytesOf(dst, len)
	sb := bytesOf(src, len)
	copy(db, sb)
}

func StringValueOf(v interface{}) (string, error) {
	switch res := v.(type) {
	case string:
		return res, nil
	case []string:
		return strings.Join(res, " "), nil

	default:
		return "", fmt.Errorf("value is not string:%v %v", v, reflect.TypeOf(v))
	}
}

func isTrue(s string) bool {
	switch s {
	case "1", "true", "True", "t":
		return true
	}
	return false
}

func StringUnmarshal(f func()) {

}

/*

asmsdsd {

	name string maxLength=15 minLength=22 'enums=1,2,3,4,5,6'

	status int max=15 min=16

	acmds {

	}
}

*/

const (
	kb = 1024
	mb = 1024 * kb
	gb = 1024 * mb
	tb = 1024 * gb
)

func parseByteSize(s string) (int, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("invalid byte size:%v", s)
	}
	n, err := strconv.ParseFloat(s[:len(s)-1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid byte size:%v", s)
	}

	switch s[len(s)-1] {
	case 'b', 'B':
		return int(n), nil
	case 'k', 'K':
		return int(n * kb), nil
	case 'm', 'M':
		return int(n * mb), nil
	case 'g', 'G':
		return int(n * gb), nil
	case 't', 'T':
		return int(n * tb), nil
	default:
		return 0, fmt.Errorf("invalid byte size:%v", s)
	}
}

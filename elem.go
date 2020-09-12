package ngcfg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type Elem struct {
	data *LinkedMap
	idx int
}

func (e *Elem)getIdx()string{
	sidx:="index_"+strconv.Itoa(e.idx)
	e.idx ++
	return sidx
}


func (e *Elem)MarshalJSON()([]byte,error){
	return json.Marshal(e.data.data)
}

func (e *Elem)Set(k string,v interface{})error{
	_,ok:=e.data.Get(k)
	if ok{
		return fmt.Errorf("%s key has already defined",k)
	}
	//e.data[k] = v
	if k == "-"{
		e.data.Set(e.getIdx(),v)
		return nil
	}
	e.data.Set(k,v)
	return nil
}

func (e *Elem)RawMap()*LinkedMap{
	return e.data
}

func (e *Elem)Iterator()Iterator{
	return e.data.Iterator()
}

func (e *Elem)Get(key string)interface{}{
	v,_:= e.data.Get(key)
	return v
}
func NewElem()*Elem{
	return &Elem{data: NewLinkedMap()}
}
//return fist elem of array
func (e *Elem)GetString(key string)(string,error){
	v,ok:=e.data.Get(key)
	if !ok {
		return "", fmt.Errorf("key %s doesn't exists",key)
	}
	switch v.(type) {
	case []string:
		ss:=v.([]string)
		if len(ss) >0{
			return ss[0],nil
		}
		return "",nil
	}
	return "",fmt.Errorf("type of %s is object",key)
}

func (e *Elem)GetNumber(key string)(float64,error){
	s,err:=e.GetString(key)
	if err != nil{
		return 0, err
	}
	return strconv.ParseFloat(s,64)
}


func (e *Elem)GetInt(key string)(int,error){
	f,err:=e.GetNumber(key)
	if err != nil{
		return 0,err
	}
	i:=int(f)

	if float64(i) != f{
		return i,fmt.Errorf("value of %s is float not int:%v",key,f)
	}
	return i,nil
}

func (e *Elem)GetBool(key string)(bool,error){
	s,err:=e.GetString(key)
	if err != nil{
		return false, err
	}
	 bv,err:=boolOf(s)
	 if err != nil{
	 	return false,fmt.Errorf("key:%s,%w",key,err)
	 }
	 return bv,nil
}

func (e *Elem)GetArray(key string)([]string,error){
	arr,ok:=e.data.Get(key)
	if !ok{
		return nil,fmt.Errorf("key %s doesn't exists",key)
	}
	switch arr.(type) {
	case []string:
		return arr.([]string),nil
	}
	return nil,fmt.Errorf("type of %s is not array",key)
}

func (e *Elem)GetStringDef(key string,def string)string{
	s,err:=e.GetString(key)
	if err != nil{
		return def
	}
	return s
}

func (e *Elem)GetNumberDef(key string,def float64)float64{
	v,err:=e.GetNumber(key)
	if err != nil{
		return def
	}
	return v
}

func (e *Elem)GetIntDef(key string,def int)int{
	v,err:=e.GetInt(key)
	if err != nil{
		return def
	}
	return v
}

func (e *Elem)GetBoolDef(key string,def bool)bool{
	v,err:=e.GetBool(key)
	if err != nil{
		return def
	}
	return v
}

func (e *Elem)GetArrayDef(key string,def []string)[]string{
	v,err:=e.GetArray(key)
	if err != nil{
		return def
	}
	return v
}

func (e *Elem)GetElem(key string)(*Elem,error){
	v,ok:=e.data.Get(key)
	if !ok{
		return nil,fmt.Errorf("key %s does not exist",key)
	}
	switch v.(type) {
	case *Elem:
		return v.(*Elem),nil
	}
	return nil,fmt.Errorf("type of %s is not elem",key)
}

func (e *Elem)AsArray()([]string,error){
	it:=e.Iterator()
	res:=make([]string,0)
	for it.HasNext(){
		elem:=it.Next()
		switch elem.Val.(type) {
		case []string:
			res = append(res,elem.Val.([]string)...)
		default:
			return nil,fmt.Errorf("type is not array:%s",reflect.TypeOf(elem.Val).String())
		}
	}
	return res,nil
}


func boolOf(s string)(bool,error){
	switch s {
	case "true","1","","on","yes","ok":
		return true,nil
	case "false","0","off","no","never":
		return false,nil
	}
	return false,fmt.Errorf("invalid bool value of %s",s)
}

//func (e *Elem)ToArray()[]interface{}{
//	res:=make([]interface{}, 0,len(e.data.data))
//	it:=e.RawMap().MapItem()
//	for it != nil{
//		res = append(res,it.Val)
//	}
//}
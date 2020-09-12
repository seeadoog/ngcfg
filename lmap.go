package ngcfg

import "encoding/json"

// 有序map
type MapElem struct {
	Key  string
	Val  interface{}
	next *MapElem
	pre  *MapElem
	l *LinkedMap
}

func (e *MapElem)MarshalJSON()([]byte,error){
	return json.Marshal(e.Val)
}

func (m *MapElem)Next()*MapElem{
	if m.next == m.l.back{
		return nil
	}
	return m.next
}

func NewLinkedMap()*LinkedMap{
	m:=&LinkedMap{
		data: map[string]*MapElem{},
		front: nil,
		back:  nil,
	}
	m.front = new(MapElem)
	m.back = new(MapElem)
	m.front.l = m
	m.back.l = m
	m.front.next = m.back
	m.back.pre = m.front

	return m
}

type LinkedMap struct {
	data map[string]*MapElem
	front *MapElem
	back *MapElem
	iterNode *MapElem

}



func (m *LinkedMap)Len()int{
	return len(m.data)
}
//front 1->2->3->back
func (m *LinkedMap)pushBack(e *MapElem){
	pb:=m.back.pre
	pb.next = e
	e.next = m.back
	m.back.pre = e
	e.pre = pb
}


func (m *LinkedMap)Set(key string,val interface{}){
	e,ok:=m.data[key]
	if ok{
		e.Val = val
		return
	}
	e = &MapElem{Val:val,l:m,Key:key}
	m.pushBack(e)
	m.data[key]=e
}

func (m *LinkedMap)Get(key string)(interface{},bool){
	v,ok:=m.data[key]
	if !ok{
		return nil,false
	}
	return v.Val,ok
}

func (m *LinkedMap)MapItem()*MapElem{
	return m.front.next
}


func (m *LinkedMap)Iterator()Iterator{
	  m.iterNode = m.front.next
      return m
}

func (m *LinkedMap) HasNext() bool {
	 if m.iterNode ==m.back || m.iterNode == nil{
	 	return false
	 }
	 return true
}

func (m *LinkedMap) Next() *MapElem {
	v:=m.iterNode
	m.iterNode = m.iterNode.next
	return v
}

type Iterator interface {
	HasNext()bool
	Next()*MapElem
}
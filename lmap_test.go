package ngcfg

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewLinkedMap(t *testing.T) {
	a:=map[string]string{
		"ss":"胜多负少的\n",
	}
	b,_:=json.Marshal(a)
	fmt.Println(len(b),string(b))
	m:=NewLinkedMap()

	m.Set("1",1)
	m.Set("2",1)
	m.Set("3",1)
	m.Set("4",1)
	m.Set("5",1)
	m.Set("6",1)
	m.Set("7",1)
	m.Set("8",1)

	for it:=m.Iterator();it.HasNext();{
		e:=it.Next()
		fmt.Println(e.Key)
	}
}
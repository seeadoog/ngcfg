package ngcfg

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	e := NewElem()
	e.Set("name", []string{"xxxx"})
	e.Set("age", []string{"5"})
	ef := NewElem()
	ef.Set("ase", []string{"true"})
	e.Set("o", ef)
	bs, _ := toJsonRaw(e)
	fmt.Println(string(bs))

}

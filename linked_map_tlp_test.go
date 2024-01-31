package ngcfg

import (
	"encoding/json"
	"fmt"
	"testing"
)

type cfgg struct {
	Cmd *LSMap[*LSMap[A]] `json:"cmd"`
}

type A int

func TestLSMap(t *testing.T) {
	c := new(cfgg)
	err := UnmarshalFromBytes([]byte(`
cmd {
	1 {
		1 1
		2 2 
	}
	2 {
		1 1
		2 2 
	}
}
`), c)
	if err != nil {
		panic(err)
	}

	bs, _ := json.Marshal(c)

	fmt.Println(string(bs))
	fmt.Println(c)
}
func TestUM(t *testing.T) {
	c := new(cfgg)
	err := json.Unmarshal([]byte(`{"cmd":{"1":{"1":1,"2":2},"2":{"1":1,"2":2}}}`), c)
	if err != nil {
		panic(err)
	}

	fmt.Println(c.Cmd)
}

type Vad struct {
	A   string `vad:"max=5,min=3"`
	B   int    `vad:""`
	Sub string `vad:"len=3"`
	D   string `vad:"required"`
}

func Test_Name(t *testing.T) {

}

func Test_Decode(t *testing.T) {
	var a []int
	err := DecodeValue(&a, []string{"5", "5"})
	if err != nil {
		panic(err)
	}
	fmt.Println(a)
}

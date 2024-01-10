package ngcfg

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

type Ccstring string

func (c *Ccstring) UnmarshalCfg(path string, v interface{}) error {
	*c = "axxxxxxx"
	return nil
}

func TestUnmarshalFromElem(t *testing.T) {

	b := []byte(`
common {

}
server {
	listen :3333
	ports 333333 3333
}
kvs {
	- {
		a b 
		c d
	}
}

kgs {
	- {
		name fff
		age 5
	}

}
d {
	hbase 333
	redis ggggg
}
f {
	---name ffff
	---ggg ghhh
	---kbmncfg {
			ddd ggg
	}
}
e {
	dd dd
}
g ttt
ips http://172.21.157.[13,14,88-98]:8080

`)
	c := &Config{}

	e, err := Parse(b)
	if err != nil {
		panic(err)
	}
	//fmt.Println(UnmarshalFromBytes(b, c))
	fmt.Println(e)
	err = UnmarshalFromElem(e, c)
	if err != nil {
		panic(err)
	}

	bs, _ := json.Marshal(c)
	fmt.Println(string(bs))
	fmt.Println(reflect.TypeOf(c.D))

	fmt.Println(c.F)
	fmt.Println(c.E)
	fmt.Println(reflect.TypeOf((Ccstring("x"))).Implements(unmarshalType))

}

func TestMarsha(t *testing.T) {
	bs := `
a 5
name str
cccc "
\n
"
ab {

}

abc 
`
	e, err := Parse([]byte(bs))
	if err != nil {
		panic(err)
	}
	cc, _ := e.MarshalCfg(0)

	fmt.Println(string(cc))

}

func TestUnmarshalRendByEnvs(t *testing.T) {
	bs := `
a 5
name str
path {{.GOOS}}
gos {{.Gos}}
`
	type cfg struct {
		A    string `json:"a"`
		Name string `json:"name"`
		Path string `json:"path"`
		Gos  string `json:"gos"`
	}
	var c = new(cfg)
	UnmarshalWithRendByEnvs([]byte(bs), map[string]string{
		"Gos":  "haha",
		"GOOS": "linux",
	}, c)
	fmt.Println(c)
}

func Test_Str(t *testing.T) {
	fmt.Println(strings.Cut("abcds,22", ","))
	fmt.Println(reflect.TypeOf(BasicValue{}))

}

type target struct {
	Weight string `json:"weight,omitempty"`
}

type upstream struct {
	Addr        map[string]target `json:"addr,omitempty"`
	HashHeaders []string          `json:"hash_headers,omitempty"`
	HashQueries []string          `json:"hash_queries,omitempty"`
	HashOn      string            `json:"hash_on,omitempty"`
}

type nginx struct {
	Http struct {
		Upstream *LSMap[*upstream] `json:"upstream,omitempty"`
		Server   map[string]Server `json:"server,omitempty"`
	} `json:"http,omitempty"`
}

func Test_String3(t *testing.T) {
	var a *Elem

	err := UnmarshalFromFile("test.cfga", &a)
	if err != nil {
		panic(err)
	}

	fmt.Println(a)

	j := json.NewEncoder(os.Stdout)
	j.SetIndent("", "\t")
	j.SetEscapeHTML(false)
	j.Encode(a)

	fmt.Println(int('\t'), int(' '))
}

type jsono struct {
	Name string `json:"name,omitempty,global"`
}

func Test_jo(t *testing.T) {
	n := jsono{
		Name: "xx",
	}
	ns, err := json.Marshal(n)
	fmt.Println(string(ns), err)

	n2 := new(jsono)

	err = json.Unmarshal(ns, n2)
	fmt.Println(n2, err)

}

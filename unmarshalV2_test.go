package ngcfg

import (
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
	name ffff
	ggg ghhh
	kbmncfg {
			ddd ggg
	}
}
e {
	dd dd
}
g ttt
ips "http://10.1.87.{50...59}:8001/ats/traffics"
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

	fmt.Println(reflect.TypeOf(c.D))

	fmt.Println(c.F)
	fmt.Println(c.E)
	fmt.Println(reflect.TypeOf((Ccstring("x"))).Implements(unmarshalType))

	cc, _ := e.MarshalCfg(0)

	fmt.Println(string(cc))

}

func TestMarsha(t *testing.T) {
	bs := `
a 5
name str

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
}

func Test_ENV(t *testing.T) {
	type Inner struct {
		Name string `env:"NAME"`
	}
	type cfg struct {
		Name string `env:"NAME"`
		In   Inner  `json:"in"`
	}
	os.Setenv("NAME", "hello world")
	c := new(cfg)
	err := UnmarshalFromBytes([]byte(`
	name dd
	in {
		
	}
	`), c)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
}

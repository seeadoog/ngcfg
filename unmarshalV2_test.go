package ngcfg

import (
	"fmt"
	"reflect"
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
	fmt.Println(c.Ips)
	fmt.Println(reflect.TypeOf((Ccstring("x"))).Implements(unmarshalType))
}
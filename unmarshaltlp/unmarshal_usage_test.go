package unmarshaltlp

import (
	"fmt"
	"github.com/seeadoog/ngcfg"
	"testing"
)

type cfg struct {
	Addr  *ConsecutiveString `json:"addr"`
	Size  *ByteSize          `json:"size"`
	Route map[string]struct {
		ProxyPass string     `json:"proxy_pass"`
		Cmd       [][]string `json:"cmd"`
	} `json:"route"`
}

func TestConsecutiveIps_UnmarshalCfg(t *testing.T) {
	c := &cfg{}
	err := ngcfg.UnmarshalFromBytes([]byte(`
addr '1.2.3.4.{5...6}:{7...8}:{12...20}'
size 2.5m

route /api1 {
	proxy_pass http://1.2.3,45 dd
	cmd {
		- set_header a b 
		- set_header c d 
		- set_resp_header http_pp_acc $get_header("dd")
		
	}
}

route /api2{

}

`), c)
	if err != nil {
		panic(err)
	}

	fmt.Println(c.Addr.Strings())
	fmt.Println(c.Size.String())
	fmt.Println(c.Size.Size())
	fmt.Println(c.Route)
	for _, s := range c.Route {
		for _, strings := range s.Cmd {
			for _, s2 := range strings {
				fmt.Println(s2)
			}
			fmt.Println("-----")
		}
	}
}

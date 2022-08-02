package values

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/seeadoog/ngcfg"
	"reflect"
	"testing"
)

type cfg struct {
	Addr  *ConsecutiveString `json:"addr"`
	Size  *ByteSize          `json:"size" default:"1k"`
	Route map[string]struct {
		ProxyPass string     `json:"proxy_pass"`
		Cmd       [][]string `json:"cmd"`
	} `json:"route"`
	Cops      *TagValue[int] `json:"cops"`
	Dbs       *configMap     `json:"dbs"`
	Timeout   *Timeduration  `json:"timeout" default:"15m"`
	Jsonblock string         `json:"jsonblock"`
}

type configFactory interface {
	Type() string
	Config() config
	New(cfg interface{}) (interface{}, error)
}

type configMap struct {
	cfg map[string]config
}

var (
	factoryMap = map[string]configFactory{}
)

type config interface {
	Type() string
}

func (c *configMap) UnmarshalCfg(path string, vv interface{}) error {
	c.cfg = map[string]config{}
	switch v2 := vv.(type) {

	case *ngcfg.Elem:
		it := v2.Iterator()
		for it.HasNext() {
			v3 := it.Next()

			switch v := v3.Val.(type) {
			case *ngcfg.Elem:
				t := v.GetStringDef("type", "")
				fact := factoryMap[t]
				if fact == nil {
					return fmt.Errorf("unknown factory:'%v' ", t)
				}

				confInst := fact.Config()
				err := v.Decode(confInst)
				if err != nil {
					return err
				}
				c.cfg[v3.Key] = confInst
				return nil
			}
		}
	}

	return fmt.Errorf("invalid cfg")
}

type redisConfig struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type redisFactory struct {
}

func (r *redisFactory) Type() string {
	return "redis"
}

func (r *redisFactory) Config() config {
	return &redisConfig{}
}

func (r *redisFactory) New(cfg interface{}) (interface{}, error) {

	return nil, nil
}

func (r *redisConfig) Type() string {
	return "redis"
}

var (
	_ = func() int {
		fmt.Println("rdis resi ")
		factoryMap["redis"] = &redisFactory{}

		return 0
	}()
)

func TestConsecutiveIps_UnmarshalCfg(t *testing.T) {
	c := &cfg{}
	err := ngcfg.UnmarshalFromBytes([]byte(`
addr '1.2.3.4.{5...6}:{7...8}:{12...20}' '10.1.98.21:{9889...9896}'
size 2.5m

route /api1 {
	proxy_pass http://1.2.3,45 dd
	cmd {
		- set_header a b 
		- set_header c d 
		- set_resp_header http_pp_acc $get_header("dd")
		- set_header c = d 
		- set_response_header request-id $request-id
		- set_request_header request-id rand_id()
	}


}

route /api2{

}

cops 14 name=5 age=7

dbs{
	asc {
		type redis
		name 12123
		addr 10.1.87.590
	}
}
timeout1 100ms

jsonblock '
{
	"haha":"dsts"
}

'

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
	fmt.Println(c.Cops)
	fmt.Println(c.Cops.GetTag("name"))
	fmt.Println(c.Cops.Val())

	var a ngcfg.Unmarshaller = &configMap{}

	ta := reflect.TypeOf(a)
	fmt.Println(ta)

	fmt.Println(c.Dbs.cfg["asc"])
	fmt.Println(c.Timeout)

	bs, _ := json.Marshal(c)

	bf := bytes.NewBuffer(nil)
	json.Indent(bf, bs, "", "    ")

	fmt.Println(bf.String())

}

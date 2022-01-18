package ngcfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

/**
aa bb
cc dd
 */
var cfg = []byte(` 
server {  # server config 
	host 127.0.0.1  # server listen host
	port 8080  # server listen port
	tps aaa bbb ccc \
		ddd eee fff \
		ggg hhh eee   # sf
     bingo df


    ans{
		name 123 
 		dfg 123 456 \
			123 123
		du "hw sdf ok thanks"
	}

    handlers{
		auth_by_lua_block "
			if ctx.param.app_id=='4cc5779a'
			then
				ctx.exit(403,'file doesnot exits')
			end
		"
		
		log_by_lua_block "
			ctx.log.info('log .... ',\"hello\")
		"	
		
		content_by_lua "
			
			ctx.writer.write(200,{
				code = 0,
				data = {
					id = ctx.response.consumer_id
				}
			})
		"
		lua "dsfdsf afdf afd"
	}



}
workerProcess $OS(cpuNum)
http proxy{
	
}	

http admin{
	
}

`)
func Test_parse(t *testing.T) {
	e,err:=parse(cfg)
	fmt.Println(e,err)
	ser:=e.Get("server").(*Elem)
	fmt.Println(ser.Get("host"))
	fmt.Println(ser.Get("port"))
	fmt.Println(ser.Get("tps"))
	fmt.Println(ser.GetBool("bingo"))
	ans:=ser.Get("ans").(*Elem)
	fmt.Println(ans.GetNumber("name"))
	fmt.Println(ans.Get("dfg"))
	fmt.Println(ans.GetString("du"))
	handlers:=ser.Get("handlers").(*Elem)
	fmt.Println(handlers.GetString("auth_by_lua_block"))
	fmt.Println(handlers.GetString("log_by_lua_block"))
	fmt.Println(handlers.GetString("content_by_lua"))
}


type Config struct {
	Common struct{

	} `json:"common"`
	Server struct{
		Listen string `json:"listen"`
		Options map[string]string `json:"options"`
		Ports []string `json:"ports"`
	} `json:"server"`

	Kvs []map[string]string `json:"kvs"`

	Kgs []*Kgs `json:"kgs"`
	E *Elem `json:"e"`
}

type Kgs struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func TestParseCfg(t *testing.T){
	b,err:=ioutil.ReadFile("test.cfg")
	if err != nil{
		panic(err)
	}
	c:=&Config{}
	fmt.Println(UnmarshalFromBytes(b,c))
	fmt.Println(c)
}

type Server struct {
	Proto string `json:"proto"`
	Listen []string `json:"listen"`
	AccessByLua string `json:"access_by_lua"`
}

type Upstream struct {
	Hosts []string `json:"hosts"`
	Targets []string `json:"targets"`
}

type Mysql struct {
	Addr string `json:"addr"`
	Password string `json:"password"`
}
type Redis struct {
	Addr string `json:"addr"`
	Password string `json:"password"`
}

type Storage struct {
	Redis Redis `json:"redis"`
	Mysql Mysql `json:"mysql"`
}
type NginxServer struct {
	CfgJson string `json:"cfg_json"`
	WorkerProcess int `json:"worker_process"`
	Server *Server `json:"server"`
	Upstreams []Upstream `json:"upstreams"`
	Storage map[string]Redis `json:"storage"`
	Args map[string]string `json:"args"`
	E Elem `json:"e"`
	Cmds []string `json:"cmds"`
	Schema string `json:"schema"`
	Ids []int `json:"ids"`
	Onj interface{} `json:"onj"`
}

func TestDemo(t *testing.T){
	c:=&NginxServer{}
	cfg:=`
worker_process 5  #进程数量   
onj{
	aaa t
}
server{
    proto   http   # protocols
    # listen addrs 
    listen  0.0.0.0:8000 0.0.0.0:8001 0.0.0.0:8002 \ 
            0.0.0.0:8003 0.0.0.0:8004 0.0.0.0:8005  # ffff
    
    access_by_lua "
        ngx.log.info('access',\"user\")

    "

}

upstreams{ # 如果upstream 模版是数组，那么server 就会被当作数组元素处理，忽略key，但是key 不能重复, 或者 key 为- 会自动生成索引id
    - {
        hosts www.test.com www.test.cn
        targets 192.168.23.12:9004 192.168.23.12:9003 "a.b.c ee"
    }
    - {
        hosts www.test2.com www.test2.cn
        targets 192.168.23.12:9004 192.168.23.12:9003
    }
	- {

    }
}

storage{
    mysql{
        addr 192.33.22.22
        password 123456
    }
    
    redis{
        addr 192.33.22.22
        password 123456
	}

}

storage mysql1{
	addr 192.33.22.22
    password 123456
}

storage redis2{
	addr 192.33.22.22
}
password 12345689999

server http{
	
}

server tcp{

}

"cfg_json" '{"name":"string"}'

args {

	"sdfsdf"  sdfsdf
	"#ffsd sfd" sdfsff
}
e{
	name strr
	age  33
	child{
		name ca
		age 5
	}
}

mysqls{
	1{
		addr 22222
		password xxxxxx
	}
	
}

schema '
{
	"type":"object",
	"properties":""
}'

cmds{
	- aaa
	- bbb
	- 'gggg 滚滚滚'
	- 'ls -lh as a '
	
}

resource redis{
	
}

resource mysql{

}

resource rmq{
	
}

ids{
	- 1 2 3
	- 4 5 6
}

`
	if err:=UnmarshalFromBytesCtx([]byte(cfg),c);err != nil{
		panic(err)
	}

	b,_:=json.Marshal(c)

	fmt.Println(string(b))
	fmt.Println(c.E.GetBool("name"))
	//it:=c.E.Iterator()
	e,err:=c.E.GetElem("child")
	if err != nil{
		panic(err)
	}
	fmt.Println(e.GetString("age"))
	fmt.Println(c.E.Get("abds"))
	fmt.Println(c)

}

func TestParseCtx(t *testing.T){
	cfg:=[]byte(`

workers 1
timeout 5
server api1{
	workers 5
}

server api2{
	workers 10
}

api /v2/create{
	method post
	timeout 50
}

api /v2/delete{
	method delete
}


method get 

api /v2/update{
	- ffffff
    -  ggggg
	- ffffff ggggg '{ws:{}}'
}

services{
	ddd dddd
}

http{
	ssl on
	apiv1{
		ssl off
	}
}

http sss{
	ssl off 
	scripts{
	
	}
}

cf_config{
	d.f {
		spdNwwfpsc 5
		nnslpll  5
		gpu_id 0
	}

# oh my god
# 
#
}

`)
	_,err:=Parse(cfg)
	if err != nil{
		panic(err)
	}
	//b,_:=json.Marshal(e)
	//
	//it:=e.Iterator()
	//for it.HasNext(){
	//	e:=it.Next()
	//	fmt.Println(e.Key,e.Val)
	//}
	////
	//fmt.Println(string(b))
	//s,err:=e.GetElem("server")
	//fmt.Println(err)
	//api1,_:=s.GetElem("api1")
	//fmt.Println(api1.GetCtxString("workers"))
	//fmt.Println(e.Elem("api").Elem("/v2/update").AsStringArray())
	//fmt.Println(e.Elem("http").GetBool("ssl"))
	//fmt.Println(e.Elem("http").Elem("wjge").GetCtxBool("ssl"))

	if err != nil{
		panic(err)
	}
}


type NginxConf struct {
	WorkerProcess int
	Http map[string]struct{

	} `json:"http"`
	Tcp map[string]struct{

	} `json:"tcp"`

}

func TestFile(t *testing.T){
	b,err:=ioutil.ReadFile(`test.cfga`)
	if err != nil{
		panic(err)
	}
	e,err:=parse(b)
	if err != nil{
		panic(err)
	}

	ss,_:=json.Marshal(e)
	fmt.Println(string(ss))
	fmt.Println(e.Elem("datas").AsArray())
}

// 1 core   3000000 goroutines per seconds
func TestDefault(t *testing.T){
	type Base struct {
		GG string `json:"gg"`
		EE string `json:"ee"`
	}

	type Cfg struct {
		Base
		Name string `json:"name" required:"true"`
		Age string `json:"age" default:"5"`
		Swa struct{
			Nae string `json:"nae" default:"556" required:"true"`
		} `json:"swa" required:"true"`
	}

	cfgs:=`
gg 688
name 5
swa{
nae 5
}
`
	v:=&Cfg{}
	err:=UnmarshalFromBytes([]byte(cfgs),v)
	if err != nil{
		panic(err)
	}
	fmt.Println(v.GG)

	//tp:=reflect.TypeOf(*v)
	//for i:=0;i<tp.NumField();i++{
	//	ft:=tp.Field(i)
	//	fmt.Println(ft.Anonymous,ft.Name)
	//}
	fmt.Println(v)
}

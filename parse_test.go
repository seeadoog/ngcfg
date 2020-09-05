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
		
	}
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
	Storage *Storage `json:"storage"`
	Args map[string]string `json:"args"`
	E Elem `json:"e"`
	Cmds []string `json:"cmds"`
	Schema string `json:"schema"`
}

func TestDemo(t *testing.T){
	c:=&NginxServer{}
	cfg:=`
worker_process 5  #进程数量       
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




`
	if err:=UnmarshalFromBytes([]byte(cfg),c);err != nil{
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
## ngcfg
一款类似nginx 格式配置文件的 golang 解析工具。相比与各种配置配置文件有着独特的优势
- toml ：风格简洁，但是对结构化配置表达不友好，
- json： 结构化强， 但是无法注释，很少用来做配置文件
- xml：结构化能力强，但是语法太啰嗦。编写容易出错
- yaml：语法简洁，结构化能力强，并且支持注释，但是格式容易写错

## 优势与特点
-  结构化表达能力强，各种结构都能表达出来
-  配置简洁，不易出错
-  支持配置换行，应对配置文较多时的情况
-  支持配置块，块里面可以放置任何内容
-  支持注释
-  支持struct模板解析
-  支持配置继承，配置字段缺失时，可以去父区块寻找相应的字段并获取其值（需要使用UnmarshalFromBytesCtx）
-  支持从环境变量绑定配置
#### 配置示例
```
server{
    proto   http   # protocols
    # listen addrs 
    listen  0.0.0.0:8000 0.0.0.0:8001 0.0.0.0:8002 \
            0.0.0.0:8003 0.0.0.0:8004 0.0.0.0:8005
    listen_tcp :8932 ssl backlog=10240 reuseport
    
    access_by_lua "
        ngx.log.info('access',\"user\")
        if ngx.err then
            ngx.exit(500,{message='haha'})
        end
    "

}
#这个地方表示数组时，upstream 中的每个元素的key 必须不一样。
upstreams{
    server1{
        hosts www.test.com www.test.cn
        targets 192.168.23.12:9004 192.168.23.12:9003 "a.b.c ee"
    }
    server2{
        hosts www.test2.com www.test2.cn
        targets 192.168.23.12:9004 192.168.23.12:9003
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


```
#### 使用示例

```go
cfg:=`
user{
    name ssss
}
listen :80 ssl 

`

type User struct{
    Name string `json:"name" default:"jhon"  env:"NAME"`
}
type ListenOpt struct{
    Ssl bool `json:"ssl"`
    Reuseport bool `json:"reuseport"`
    Backlog int `json:"backlog" default:"10240"`
}
type Config struct{
    User User `json:"user" required:"true"`
    Listen values.TagValueT[string,]
}

cfg:=&Config{}
UnmarshalFromBytes([]byte(cfg),cfg)


```


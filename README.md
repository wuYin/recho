## Recho

Recho 是一个封装 [Echo framework](https://github.com/labstack/echo) 路由、中间件与常用 HTTP 响应函数的 toolbox



## Quick Start

#### 下载依赖

```
go get github.com/labstack/echo
go get github.com/naoina/toml
go get github.com/wuYin/recho
```


#### 创建路由文件 `routes.toml`

```toml
[routes]
"GET:/user"="handlers.User.GetUserInfo"

[validators]
"/"=["validators.RespMiddleware.Process"]
"/user"=["validators.User.CheckSession", "!validators.User.UnusedChecker"]	# !开头的会跳过
```



#### 创建运行文件 `main.go`

```go
package main

import (
	"recho/utils"
	"recho/handlers"
	"recho/validators"
)

func main() {
	s := utils.InitEnv("./routes.toml")
	s.RegisterHandler(&handlers.User{})               // 业务逻辑处理器
	s.RegisterValidator(&validators.User{})           // 业务逻辑验证中间件
	s.RegisterValidator(&validators.RespMiddleware{}) // 响应中间件
	s.Start(":2333")
}
```



#### 运行后请求： [127.0.0.1:2333/user](http://127.0.0.1:2333/user)

```
go run main.go
```



#### 服务端处理请求：

![](http://p2j5s8fmr.bkt.clouddn.com/new-recho-run.png)

 

#### 客户端接收响应：

 ![](http://p2j5s8fmr.bkt.clouddn.com/resp_succ2.png)



## Features

- [x] 封装路由与中间件到一个配置文件
- [x] 封装常用的 HTTP 响应函数



## Structures

```
➜  recho git:(master) tree -L 2
.
├── handlers
│   └── User.go 	# User 业务处理
├── validators
│   ├── Resp.go 	# 响应中间件
│   └── User.go 	# User 业务验证
├── utils
│   ├── codes.go	# 状态码与状态信息
│   ├── http.go 	# 封装 HTTP 响应函数	
│   └── utils.go	# 封装配置文件
├── main.go
└── routes.toml 	# 路由与中间件配置文件
```

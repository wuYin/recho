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
"/user"=["validators.User.CheckSession", "!validators.User.UnusedChecker"]
```



#### 创建运行文件 `main.go`

```go
package main

import (
	"recho/handlers"
	"recho/utils"
	"recho/validators"
)

func main() {
	s := utils.InitEnv("./routes.toml")
	s.RegisterHandler(&handlers.User{})
	s.RegisterValidator(&validators.User{})
	s.Start(":2333")
}
```



#### 运行：

```
go run main.go
```

请求：[127.0.0.1:2333/user](http://127.0.0.1:2333/user)



#### 调用成功

![](http://p2j5s8fmr.bkt.clouddn.com/new-recho-run.png)

 

## Structures

```
➜  recho git:(master) tree -L 2
.
├── handlers
│   └── User.go 	# 业务处理
├── validators
│   └── User.go 	# 验证处理
├── utils
│   └── utils.go    # 封装细节	
├── main.go 		# 服务运行文件
└── routes.toml 	# 路由与中间件配置文件
```





## TODO

- [x] 封装路由与中间件到一个配置文件
- [ ] 封装常用的 HTTP 响应函数
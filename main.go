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

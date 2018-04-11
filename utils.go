package recho

import (
	"github.com/labstack/echo"
	"reflect"
)

// 处理器
type handler struct {
	handleName string        // 名字
	httpMethod string        // 请求方法
	handleFunc reflect.Value // 函数实现
}

// 验证器（中间件）
type validator struct {
	routePrefix string                  // 验证的路由前缀
	handleName  string                  // 处理验证的函数
	skipRoutes  map[string]*interface{} // 无需验证的路由
	handleFunc  echo.MiddlewareFunc     //处理验证的函数
}

// 封装后的 Echo Server
type RechoServer struct {
	server            *echo.Echo
	route2Handler     map[string]*handler   // 一条路由 对 一个处理器
	route2Validators  map[string]*validator // 一条路由 对 多个验证器
	handler2Routes    map[string][]string   // 一个处理器 对 多条路由
	handler2Validator map[string]*validator // 一个处理器 对 多个验证器
}

// toml 路由配置项
type Conf struct {
	Routes     map[string]string   // 路由配置
	Validators map[string][]string // 验证器配置
}

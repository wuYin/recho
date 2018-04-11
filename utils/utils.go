package utils

import (
	"github.com/labstack/echo"
	"reflect"
	"log"
	"fmt"
	"github.com/naoina/toml"
	"io/ioutil"
	"strings"
	"net/http"
	"sort"
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

//
// 读取路由配置项，初始化环境
//
func InitEnv(confPath string) {

	// 读取配置项
	data, err := ioutil.ReadFile(confPath)
	checkErr(err, "Can't open toml conf file.")
	var conf Conf
	err = toml.Unmarshal(data, &conf)

	// 初始化 Recho Server
	// 映射路由与处理器
	route2Handler, handler2Routes := mapRouteAndHandlers(conf)
	pr(route2Handler, handler2Routes)
}

//
// 映射路由与处理器
//
func mapRouteAndHandlers(conf Conf) (route2Handler map[string]*handler, handler2Routes map[string][]string) {
	route2Handler = make(map[string]*handler, 1)
	handler2Routes = make(map[string][]string, 1)
	for r, h := range conf.Routes {
		splits := strings.SplitN(r, ":", 2) // 分割路由，如 "POST:/user"
		method := http.MethodGet
		if len(splits) > 1 {
			method = strings.ToUpper(splits[0])
			r = splits[1]
		}
		// 建立 route 与 handler 的一对一关系
		route2Handler[r] = &handler{
			handleName: h,
			httpMethod: method,
		}

		// 建立 handler 与 routes 的一对多关系
		routes, ok := handler2Routes[h]
		if !ok {
			routes = make([]string, 0, 1)
		}
		routes = append(routes, r)
		sort.Strings(routes)
		handler2Routes[h] = routes
	}

	return
}

//
// 错误检查
//
func checkErr(err error, info string) {
	if err != nil {
		log.Fatalln(info, err)
	}
}

//
// 调试函数
//
func pr(vs ...interface{}) {
	for _, v := range vs {
		fmt.Printf("%v\n", v)
	}
}

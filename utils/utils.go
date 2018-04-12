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
	route2Handler     map[string]*handler     // 一条路由 对 一个处理器
	route2Validators  map[string][]*validator // 一条路由 对 多个验证器
	handler2Routes    map[string][]string     // 一个处理器 对 多条路由
	handler2Validator map[string]*validator   // 一个处理器 对 多个验证器
	validateRoutes    []string                // 需要做验证的路由
}

// toml 路由配置项
type Conf struct {
	Routes     map[string]string   // 路由配置
	Validators map[string][]string // 验证器配置
}

//
// 读取路由配置项，初始化环境
//
func InitEnv(confPath string) *RechoServer {

	// 读取配置项
	data, err := ioutil.ReadFile(confPath)
	checkErr(err, "Can't open toml conf file.")
	var conf Conf
	err = toml.Unmarshal(data, &conf)

	// 初始化 Recho Server
	// 映射路由与处理器
	route2Handler, handler2Routes := mapRouteAndHandlers(conf)
	// 映射路由与验证器
	validateRoutes, route2Validators, handler2Validator := mapRouteAndValidators(conf)

	return &RechoServer{
		echo.New(),
		route2Handler,
		route2Validators,
		handler2Routes,
		handler2Validator,
		validateRoutes,
	}
}

//
// 将处理器反射到包中的函数
//
func (s *RechoServer) RegisterHandler(handlerPath interface{}) *RechoServer {
	rHVal := reflect.ValueOf(handlerPath)
	rHType := rHVal.Elem().Type()
	rHPath := rHType.String() // package.struct.HandlerFunc

	// 要处理的所有路由
	hs := make([]string, 0, 1)
	for h := range s.handler2Routes {
		hs = append(hs, h)
	}
	sort.Strings(hs)

	used := false
	// 遍历所有 handler 下的所有 route
	for _, h := range hs {
		routes := s.handler2Routes[h]
		for _, route := range routes {
			// handlerPath 函数要处理的路由
			if strings.HasPrefix(h, rHPath) {
				handleFuncName := strings.TrimPrefix(strings.TrimPrefix(h, rHPath), ".")
				handleFunc := rHVal.MethodByName(handleFuncName)
				if handleFunc.Kind() == reflect.Invalid || handleFunc.IsNil() {
					log.Panicf("[ERROR]: HandleFunc %s Not Exist In %s", handleFuncName, rHPath)
				}

				s.route2Handler[route].handleFunc = handleFunc
				used = true
				log.Printf("[INFO]: Register Succeed: %s -> %s.%s", route, rHPath, handleFuncName)
			}
		}
	}

	if !used {
		log.Printf("[WARN]: HandlerFunc Not Used: %s", rHPath)
	}
	return s
}

//
// 将验证器反射到包中的函数
//
func (s *RechoServer) RegisterValidator(vPath interface{}) *RechoServer {
	vValue := reflect.ValueOf(vPath)
	vType := vValue.Elem().Type()
	vName := vType.String()

	// 所有待验证的路由
	routes := make([]string, 0, 1)
	for r := range s.route2Validators {
		routes = append(routes, r)
	}
	sort.Strings(routes)

	used := false
	// 遍历处理所有需要验证的路由
	for _, r := range routes {
		// 一条路由对应多个验证器
		vs := s.route2Validators[r]
		for _, v := range vs {
			// 注册的处理器处理当前路由
			if strings.HasPrefix(v.handleName, vName) {
				validateFuncName := strings.TrimPrefix(strings.TrimPrefix(v.handleName, vName), ".")
				validateFunc := vValue.MethodByName(validateFuncName)

				if validateFunc.Kind() == reflect.Invalid || validateFunc.IsNil() {
					log.Panicf("[ERROR]: ValidateFunc %s Not Exist In %s", validateFuncName, vName)
				} else {
					// 检查类型
					ok := validateFunc.Type().ConvertibleTo(reflect.TypeOf((func(echo.HandlerFunc) echo.HandlerFunc)(nil)))
					if !ok {
						log.Panicf("[ERROR]: ValidateFunc %s Not MiddlewareFunc", validateFuncName)
					}

					// 建立映射关系
					v.handleFunc = validateFunc.Interface().(func(echo.HandlerFunc) echo.HandlerFunc)
					used = true
				}

				log.Printf("[INFO]: Register Succeed: %s -> %s.%s", r, vName, validateFuncName)
			}
		}
	}

	if !used {
		log.Printf("[WARN]: %s Not Used", vName)
	}

	return s
}

//
// 启动 RechoServer
//
func (s *RechoServer) Start(port string) error {
	// 检查所有验证函数是否都注入
	for r, vs := range s.route2Validators {
		for _, v := range vs {
			if nil == v.handleFunc {
				log.Fatalf("[ERROR]: %s -> %s Is Not Injected", r, v.handleName)
			}
		}
	}

	// 取出所有路由
	rs := make([]string, 0, 1)
	for r := range s.route2Handler {
		rs = append(rs, r)
	}
	sort.Strings(rs)

	// 将处理函数注册到 echo.Server
	for _, r := range rs {
		h := s.route2Handler[r]
		m := strings.ToUpper(h.httpMethod)
		if h.handleFunc.Kind() == reflect.Invalid {
			if m != "FILE" && m != "STATIC" {
				log.Fatalf("[ERROR]: %s -> %s Is Not Injected", r, h.handleName)
			}
		}

		handleFunc := func(ctx echo.Context) error {
			context := reflect.ValueOf(ctx)
			h.handleFunc.Call([]reflect.Value{context})
			return nil
		}

		usedVs := make([]echo.MiddlewareFunc, 0, 1)
		for _, prefix := range s.validateRoutes {
			if strings.HasPrefix(r, prefix) {
				vs, ok := s.route2Validators[r]
				if ok {
				FLAG:
					for _, v := range vs {
						for skipPrefix := range v.skipRoutes {
							if strings.HasPrefix(r, skipPrefix) {
								log.Printf("[INFO]: Vlidator Skipped: %s -/-> %s", r, v.handleName)
								continue FLAG
							}
						}
						usedVs = append(usedVs, v.handleFunc)
					}
				}
			}
		}

		// 发布路由所需的处理器和验证器
		switch m {
		case http.MethodGet:
			s.server.GET(r, handleFunc, usedVs...)
		case http.MethodPost:
			s.server.POST(r, handleFunc, usedVs...)
		case http.MethodHead:
			s.server.HEAD(r, handleFunc, usedVs...)
		case "FILE":
			s.server.File(r, h.handleName)
		case "STATIC":
			s.server.Static(r, h.handleName)
		}
	}

	log.Fatalln(s.server.Start(port))
	return nil
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
// 映射路由与处理器
//
func mapRouteAndValidators(conf Conf) ([]string, map[string][]*validator, map[string]*validator) {
	validateRoutes := make([]string, 0, len(conf.Validators))
	route2Validators := make(map[string][]*validator, 1)
	handler2Validator := make(map[string]*validator, 1)
	// 遍历验证器组
	for route, handlers := range conf.Validators {
		// 建立 route 与 validators 的一对多关系
		validators, ok := route2Validators[route]
		if !ok {
			validators = make([]*validator, 0, len(handlers))
			validateRoutes = append(validateRoutes, route)
		}

		// 建立 handler 与 validator 一对一的关系
		for _, h := range handlers {
			realH := strings.TrimPrefix(h, "!")
			v, ok := handler2Validator[realH]
			if !ok {
				v = &validator{
					handleName: realH,
					skipRoutes: make(map[string]*interface{}, 1),
				}
				handler2Validator[h] = v
			}

			if strings.HasPrefix(h, "!") {
				v.skipRoutes[route] = nil // ! 开头的验证器忽略
			}

			validators = append(validators, v)
		}

		route2Validators[route] = validators
	}
	sort.Strings(validateRoutes)
	return validateRoutes, route2Validators, handler2Validator
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
func Pr(vs ...interface{}) {
	for _, v := range vs {
		fmt.Printf("%v\n", v)
	}
}

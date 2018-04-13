package utils

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//
// 响应状态
//
var RESP_SUCC = Status{0, "请求成功"}
var RESP_FAIL = Status{-1, "请求失败"}

//
// 服务器相关
//
// 发生异常
var RESP_S_INTERNAL_ERROR = Status{500, "服务器内部错误，请稍后再试"}

//
// 客户端相关
//
// 参数错误
var RESP_C_PARAMS_INVALID = Status{400, "参数错误"}

// 根据业务需求
// 定义更多自己的状态码与描述信息
// ...

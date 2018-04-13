package utils

import "github.com/labstack/echo"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const RESP_DATA_KEY = "resp_data"

//
// 请求成功
//
func SendServerSucc(ctx echo.Context, data interface{}) error {
	resp := Response{RESP_SUCC.Code, RESP_SUCC.Message, data}
	ctx.Set(RESP_DATA_KEY, resp)
	return nil
}

//
// 客户端参数错误
//
func SendParamsInvalid(ctx echo.Context, data interface{}) error {
	resp := Response{RESP_C_PARAMS_INVALID.Code, RESP_C_PARAMS_INVALID.Message, nil}
	ctx.Set(RESP_DATA_KEY, resp)
	return nil
}

//
// 服务器内部错误
//
func SendServerError(ctx echo.Context, data interface{}) error {
	resp := Response{RESP_S_INTERNAL_ERROR.Code, RESP_S_INTERNAL_ERROR.Message, nil}
	ctx.Set(RESP_DATA_KEY, resp)
	return nil
}

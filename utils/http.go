package utils

import "github.com/labstack/echo"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const RESP_DATA_KEY = "resp_data"

//
// 200 OK
//
func SendServerSucc(ctx echo.Context, data interface{}) error {
	resp := Response{SERVER_OK.Code, SERVER_OK.Message, data}
	ctx.Set(RESP_DATA_KEY, resp)
	return nil
}

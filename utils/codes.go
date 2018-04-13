package utils

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var SERVER_OK = Status{200, "请求成功"}

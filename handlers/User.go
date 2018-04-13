package handlers

import (
	"github.com/labstack/echo"
	"recho/utils"
	"fmt"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *User) GetUserInfo(ctx echo.Context) {
	// 一般开发 API 的三个步骤
	// 1. 检查请求参数
	// 2. 处理业务逻辑
	// 3. 发送响应

	// 参数检查
	//utils.SendParamsInvalid(ctx, nil)	// 若参数验证失败

	// 业务处理
	//utils.SendServerError(ctx, nil) // 若服务器处理出错

	fmt.Println("调用 handlers.User.GetUserInfo 处理业务逻辑 :)")
	pike := User{"Robert C. Pike", 62}
	// 发送响应
	utils.SendServerSucc(ctx, pike)
}

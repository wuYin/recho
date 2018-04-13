package handlers

import (
	"github.com/labstack/echo"
	"fmt"
	"recho/utils"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *User) GetUserInfo(ctx echo.Context) {
	fmt.Println("调用 handlers.User.GetUserInfo :)")
	pike := User{"Robert C. Pike", 62}
	utils.SendServerSucc(ctx, pike)
}

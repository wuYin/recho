package handlers

import (
	"github.com/labstack/echo"
	"fmt"
)

type User struct {
}

func (u *User) GetUserInfo(ctx echo.Context) {
	fmt.Println("调用 handlers.User.GetUserInfo :)")
}

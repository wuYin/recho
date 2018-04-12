package validators

import (
	"github.com/labstack/echo"
	"fmt"
)

type User struct {
}

func (u *User) CheckSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// 从 ctx 中取出 session 进行验证
		fmt.Println("通过验证 :)")
		next(ctx)
		return nil
	}
}

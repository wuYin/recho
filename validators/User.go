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
		fmt.Println("调用 validators.User.CheckSession 验证器做验证 :)")
		next(ctx)
		return nil
	}
}

func (u *User) UnusedChecker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fmt.Println("不会被使用的验证器 :)")
		next(ctx)
		return nil
	}
}

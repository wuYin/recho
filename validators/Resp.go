package validators

import (
	"github.com/labstack/echo"
	"net/http"
	"fmt"
)

type RespMiddleware struct {
}

const RESP_DATA_KEY = "resp_data"

func (r *RespMiddleware) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// 使用 defer 在 return 前执行的特性
		// 实现 next() 处理业务逻辑，将响应数据存在 RESP_DATA_KEY 中
		defer func() {
			data := ctx.Get(RESP_DATA_KEY)
			fmt.Printf("data: %v\n", data)
			ctx.JSON(http.StatusOK, data)
		}()
		return next(ctx)
	}
}

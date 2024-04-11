package middleware

import (
	"github.com/labstack/echo/v4"
)

// ABTestInfo 中间件将 abtest_info 从 cookies 中解析出来并带入 request header
func ABTestInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cookie, err := c.Cookie("abtest_info"); err == nil {
				c.Request().Header.Set("Abtest-Info", cookie.Value)
			}
			return next(c)
		}
	}
}

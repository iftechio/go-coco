package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthProbe() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/health" {
				return c.JSON(http.StatusOK, map[string]bool{"success": true})
			}
			return next(c)
		}
	}
}

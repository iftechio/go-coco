package middleware

import (
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Sentry 将挂载一个 hub 至 request context 上，之后可以通过 sentry.GetHubFromContext 获取 hub
func Sentry() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hub := sentry.CurrentHub().Clone()
			hub.Scope().SetRequest(c.Request())
			req := c.Request()
			req = req.WithContext(sentry.SetHubOnContext(c.Request().Context(), hub))
			c.SetRequest(req)
			defer func() {
				if v := recover(); v != nil {
					err, ok := v.(error)
					if !ok {
						err = errors.Errorf("%#v", err)
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

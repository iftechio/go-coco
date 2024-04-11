package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSentry(t *testing.T) {
	e := echo.New()
	e.Use(Sentry())
	e.GET("/error", func(context echo.Context) error {
		assert.NotNil(t, sentry.GetHubFromContext(context.Request().Context()))
		return context.NoContent(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	resp := httptest.NewRecorder()
	e.ServeHTTP(resp, req)
}

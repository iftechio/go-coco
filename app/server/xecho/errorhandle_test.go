package xecho

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/xieziyu/go-coco/app/server/custom"
	"github.com/xieziyu/go-coco/app/server/xecho/middleware"
)

func Test_handleError(t *testing.T) {
	e := New()
	e.HTTPErrorHandler = e.handleError
	e.GET("/ping", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.GET("/empty", func(c echo.Context) error {
		return nil
	})
	e.GET("/validate", func(c echo.Context) error {
		var param struct {
			Name string `json:"name" validate:"number"`
		}
		err := c.Bind(&param)
		if err != nil {
			return err
		}
		return c.Validate(param)
	})
	e.GET("/sleep", func(c echo.Context) error {
		var param struct {
			Sec int `json:"sec"`
		}
		if err := c.Bind(&param); err != nil {
			return err
		}
		time.Sleep(time.Duration(param.Sec) * time.Second)
		return c.NoContent(http.StatusOK)
	})
	e.GET("/customerror", func(c echo.Context) error {
		return &custom.Error{Status: http.StatusNoContent}
	})
	e.GET("/error", func(c echo.Context) error {
		return fmt.Errorf("")
	})
	t.Run("empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/empty", nil)
		resp := httptest.NewRecorder()
		e.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("http error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/xxx", nil)
		resp := httptest.NewRecorder()
		e.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
	t.Run("http error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/validate", strings.NewReader(`{"name":1}`))
		req.Header.Add("content-type", "application/json")
		resp := httptest.NewRecorder()
		e.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("validate error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/validate", strings.NewReader(`{"name":"abc"}`))
		req.Header.Add("content-type", "application/json")
		resp := httptest.NewRecorder()
		e.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("custom error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customerror", nil)
		req.Header.Add("content-type", "application/json")
		resp := httptest.NewRecorder()
		e.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})
	t.Run("uncatch", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		req.Header.Add("content-type", "application/json")
		resp := httptest.NewRecorder()
		e.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
	t.Run("sentry", func(t *testing.T) {
		errorCount := 0

		// launch mock sentry server
		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			errorCount++
		}))
		defer srv.Close()
		dsn, _ := url.Parse(srv.URL)
		dsn.User = url.User("root")
		err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn.String() + "/123",
		})
		assert.NoError(t, err)
		sentry.Flush(time.Second)
		assert.Equal(t, 0, errorCount)

		e := New()
		e.HTTPErrorHandler = e.handleError
		e.Use(middleware.Sentry())
		e.GET("/error", func(context echo.Context) error {
			return fmt.Errorf("nothing")
		})
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		resp := httptest.NewRecorder()
		e.ServeHTTP(resp, req)
		sentry.Flush(time.Second)
		assert.Equal(t, 1, errorCount)
	})
}

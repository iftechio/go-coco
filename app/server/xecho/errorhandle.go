package xecho

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/iftechio/go-coco/app/server/custom"
	xLogger "github.com/iftechio/go-coco/utils/logger"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func asHTTPError(err error) (*echo.HTTPError, bool) {
	httpErr := &echo.HTTPError{}
	if errors.As(err, &httpErr) {
		return httpErr, true
	}
	return nil, false
}

func asCustomError(err error) (*custom.Error, bool) {
	customErr := &custom.Error{}
	if ok := errors.As(err, &customErr); ok {
		return customErr, true
	}
	return nil, false
}

func asValidateError(err error) (validator.ValidationErrors, bool) {
	var validateErr validator.ValidationErrors
	if ok := errors.As(err, &validateErr); ok {
		return validateErr, true
	}
	return validateErr, false
}

func (s *Server) handleError(err error, c echo.Context) {
	ctx := c.Request().Context()
	var message string
	var toast string
	var status int
	var errText string
	var errData interface{}
	if errors.Is(err, context.Canceled) {
		errText = err.Error()
		xLogger.FromContext(ctx).Warning("request canceled")
	} else if e, ok := asHTTPError(err); ok {
		// echo http error
		if internalErr, ok := asHTTPError(e.Internal); ok {
			e = internalErr
		}
		status = e.Code
		if e.Message != nil {
			message = fmt.Sprintf("%v", e.Message)
		}
	} else if e, ok := asCustomError(err); ok {
		// custom error
		if e.Status != 0 {
			status = e.Status
		}
		if e.Data != nil {
			errData = e.Data
		}
		message = e.Message
	} else if e, ok := asValidateError(err); ok {
		// validator error
		status = http.StatusBadRequest
		errText = e.Error()
	} else {
		xLogger.F(ctx).Error(err)
		// unexpected error
		// should use middleware.Sentry to store sentry.Hub in context beforehand.
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.CaptureException(err)
		} else {
			xLogger.FromContext(ctx).Warning("no attached sentry hub to capture exception")
		}
		errText = err.Error()
	}
	if customErr, ok := asCustomError(err); ok {
		// custom toast
		toast = customErr.Toast
		if customErr.Err != nil {
			if hub := sentry.GetHubFromContext(ctx); hub != nil {
				hub.CaptureException(customErr.Err)
			}
		}
	}
	if !c.Response().Committed {
		if status == 0 {
			status = http.StatusInternalServerError
		}
		if message == "" {
			message = http.StatusText(status)
		}
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(status)
		} else {
			resp := custom.ResponseCommon{
				Message: message,
				Data:    errData,
			}
			if (c.Echo().Debug ||
				strings.HasPrefix(c.Request().URL.Path, "/internals/") ||
				strings.HasPrefix(c.Request().URL.Path, "/management/")) && toast == "" {
				toast = errText
			}
			if toast != "" {
				resp.Error = toast
			} else if status == http.StatusInternalServerError {
				resp.Error = "系统错误"
			}
			err = c.JSON(status, resp)
		}
		if err != nil {
			xLogger.FromContext(c.Request().Context()).Error(err)
		}
	}
}

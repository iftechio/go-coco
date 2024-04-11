package xecho

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	xLogger "github.com/xieziyu/go-coco/utils/logger"
)

type Server struct {
	*echo.Echo
}

type customValidator struct {
	validator *validator.Validate
}

func (v *customValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func New() *Server {
	e := &Server{
		Echo: echo.New(),
	}
	e.Logger = &logger{}
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &customValidator{validator: validator.New()}
	e.HTTPErrorHandler = e.handleError
	return e
}

// Start overrides echo's Start func for graceful shutdown
func (s *Server) Start(address string) error {
	xLogger.Info("server started on " + address)
	if err := s.Echo.Start(address); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

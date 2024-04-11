package custom

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ error = (*Error)(nil)

type Error struct {
	Message  string      `json:"message,omitempty"`
	Status   int         `json:"-"`
	GRPCCode codes.Code  `json:"-"`
	Toast    string      `json:"toast,omitempty"`
	Code     string      `json:"code,omitempty"`
	Err      error       `json:"-"`
	Data     interface{} `json:"data,omitempty"`
}

func (c *Error) Error() string {
	if c == nil {
		return ""
	}
	return c.Message
}

// clone 浅拷贝一份 Error
func (c *Error) clone() *Error {
	if c == nil {
		return nil
	}
	t := *c
	e := t
	return &e
}

func (c *Error) WithToast(s string) *Error {
	e := c.clone()
	e.Toast = s
	return e
}

func (c *Error) WithError(err error) *Error {
	e := c.clone()
	e.Err = err
	return e
}

func (c *Error) Unwrap() error {
	return c.Err
}

func (c *Error) GRPCStatus() *status.Status {
	if c == nil {
		return nil
	}
	gcode := c.GRPCCode
	if gcode == 0 {
		gcode = grpcCodeFromHTTPStatus(c.Status)
	}
	s := status.New(gcode, c.Error())

	return s
}

func (c *Error) WithData(d interface{}) *Error {
	e := c.clone()
	e.Data = d
	return e
}

var (
	ErrNoValidUser = &Error{
		Message: "No Valid User",
		Status:  http.StatusUnauthorized,
	}
	ErrUserLoginRequired = &Error{
		Message: "User Login Required",
		Status:  http.StatusForbidden,
		Toast:   "登录后即可操作",
		Code:    "E101",
	}
	ErrUserBanned = &Error{
		Message: "User Banned",
		Status:  http.StatusForbidden,
	}
	ErrResourceNotFound = &Error{
		Message: "Resource Not Found",
		Status:  http.StatusNotFound,
	}
	ErrBadRequest = &Error{
		Message: "Bad Request",
		Status:  http.StatusBadRequest,
		Toast:   "参数错误",
	}
	ErrPaymentError = &Error{
		Message: "Payment Service Error",
		Status:  http.StatusBadRequest,
		Toast:   "支付失败",
	}
	ErrIAPError = &Error{
		Message: "IAP Error",
		Status:  http.StatusBadRequest,
	}
	ErrNotBindPhoneNumber = &Error{
		Message: "Need Binding Phone Number",
		Status:  http.StatusForbidden,
		Toast:   "请先绑定手机号哦",
	}
	ErrSponsorRequired = &Error{
		Message: "Sponsor Required",
		Status:  http.StatusForbidden,
		Toast:   "请先开通即刻会员",
	}
	ErrNotFound = &Error{
		Message: "Not Found",
		Status:  http.StatusNotFound,
	}
)

// GRPCCodeFromHTTPStatus convert http code to grpc code
// https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md
func grpcCodeFromHTTPStatus(s int) codes.Code {
	switch s {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	default:
		return codes.Unknown
	}
}

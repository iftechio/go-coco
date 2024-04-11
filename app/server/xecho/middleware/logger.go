package middleware

import (
	"github.com/labstack/echo/v4/middleware"
)

type LoggerConfig struct {
	Skipper           middleware.Skipper
	RecordRequestBody bool
	Enable            bool
}

type LoggerConfigBuilder func(*LoggerConfig)

// 开启/关闭 Access Log
func LoggerEnable(v bool) LoggerConfigBuilder {
	return func(c *LoggerConfig) {
		c.Enable = v
	}
}

// 开启/关闭记录 Request Body
func LoggerRecordRequestBody(v bool) LoggerConfigBuilder {
	return func(c *LoggerConfig) {
		c.RecordRequestBody = v
	}
}

// 配置路由 Skipper
func LoggerSkipper(v middleware.Skipper) LoggerConfigBuilder {
	return func(c *LoggerConfig) {
		c.Skipper = v
	}
}

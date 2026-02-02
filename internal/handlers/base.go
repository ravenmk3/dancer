package handlers

import (
	"github.com/labstack/echo/v4"
)

// Response 统一响应结构
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c echo.Context, data interface{}) error {
	return c.JSON(200, Response{
		Code:    "success",
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c echo.Context, message string, data interface{}) error {
	return c.JSON(200, Response{
		Code:    "success",
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c echo.Context, code string, message string, statusCode int) error {
	return c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c echo.Context, code string, message string, statusCode int, data interface{}) error {
	return c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// 错误码常量
const (
	CodeSuccess            = "success"
	CodeError              = "error"
	CodeInvalidInput       = "invalid_input"
	CodeInvalidCredentials = "invalid_credentials"
	CodeUnauthorized       = "unauthorized"
	CodeForbidden          = "forbidden"
	CodeUserNotFound       = "user_not_found"
	CodeUserExists         = "user_exists"
	CodeRecordNotFound     = "record_not_found"
	CodeRecordExists       = "record_exists"
	CodeInternalError      = "internal_error"
)

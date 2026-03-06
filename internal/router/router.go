package router

import (
	"errors"
	"net/http"

	"dancer/internal/auth"
	apperrors "dancer/internal/errors"
	"dancer/internal/handlers"
	"dancer/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Response 统一响应结构
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func New(
	userHandler *handlers.UserHandler,
	zoneHandler *handlers.ZoneHandler,
	domainHandler *handlers.DomainHandler,
	healthHandler *handlers.HealthHandler,
) *echo.Echo {
	e := echo.New()
	e.HideBanner = true // 隐藏 Echo 默认 Banner

	// 设置全局错误处理器
	e.HTTPErrorHandler = customHTTPErrorHandler

	// 中间件
	e.Use(CustomLogger()) // 自定义访问日志中间件
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	// API 路由组
	api := e.Group("/api")

	// 健康检查（公开端点，支持 GET 和 POST）
	api.GET("/health", healthHandler.Check)
	api.POST("/health", healthHandler.Check)

	// 公开路由
	authGroup := api.Group("/auth")
	authGroup.POST("/login", userHandler.Login)
	authGroup.POST("/refresh", userHandler.RefreshToken, auth.JWTMiddleware())

	// 需要认证的路由
	me := api.Group("/me", auth.JWTMiddleware())
	me.POST("", userHandler.GetCurrentUser)
	me.POST("/change-password", userHandler.ChangePassword)

	// 用户管理（需要管理员权限）
	user := api.Group("/user", auth.JWTMiddleware(), auth.RequireAdmin())
	user.POST("/list", userHandler.ListUsers)
	user.POST("/create", userHandler.CreateUser)
	user.POST("/update", userHandler.UpdateUser)
	user.POST("/delete", userHandler.DeleteUser)

	// DNS Zone 管理（需要管理员权限）
	zones := api.Group("/dns/zones", auth.JWTMiddleware(), auth.RequireAdmin())
	zones.POST("/list", zoneHandler.ListZones)
	zones.POST("/get", zoneHandler.GetZone)
	zones.POST("/create", zoneHandler.CreateZone)
	zones.POST("/update", zoneHandler.UpdateZone)
	zones.POST("/delete", zoneHandler.DeleteZone)

	// DNS Domain 管理（需要认证）
	domains := api.Group("/dns/domains", auth.JWTMiddleware())
	domains.POST("/list", domainHandler.ListDomains)
	domains.POST("/get", domainHandler.GetDomain)
	domains.POST("/create", domainHandler.CreateDomain)
	domains.POST("/update", domainHandler.UpdateDomain)
	domains.POST("/delete", domainHandler.DeleteDomain)

	return e
}

// customHTTPErrorHandler 自定义全局错误处理器
func customHTTPErrorHandler(err error, c echo.Context) {
	// 如果响应已经写入，直接返回
	if c.Response().Committed {
		return
	}

	// 尝试从 echo.HTTPError 中提取内部错误
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		if innerErr, ok := httpErr.Message.(error); ok {
			err = innerErr
		} else {
			// Message 不是 error 类型，直接返回 HTTP 错误
			c.JSON(httpErr.Code, Response{
				Code:    "http_error",
				Message: httpErr.Error(),
			})
			return
		}
	}

	// 检查是否是业务错误
	var bizErr apperrors.BusinessError
	if errors.As(err, &bizErr) {
		c.JSON(bizErr.Status(), Response{
			Code:    bizErr.Code(),
			Message: bizErr.Error(),
		})
		return
	}

	// 未知错误
	logger.Log.WithError(err).Error("Unhandled error")
	c.JSON(http.StatusInternalServerError, Response{
		Code:    "internal_error",
		Message: "internal server error",
	})
}

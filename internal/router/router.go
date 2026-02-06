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
	authGroup.POST("/refresh", userHandler.RefreshToken)

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

	// 处理 echo 的 HTTPError
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.Code, Response{
			Code:    "http_error",
			Message: httpErr.Error(),
		})
		return
	}

	// 业务错误映射
	switch {
	// etcd 不可用
	case errors.Is(err, apperrors.ErrEtcdUnavailable):
		c.JSON(http.StatusServiceUnavailable, Response{
			Code:    "service_unavailable",
			Message: "etcd service temporarily unavailable, please retry later",
		})

	// 用户相关错误
	case errors.Is(err, apperrors.ErrUserNotFound):
		c.JSON(http.StatusNotFound, Response{
			Code:    "user_not_found",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrUserExists):
		c.JSON(http.StatusConflict, Response{
			Code:    "user_exists",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrCannotDeleteDefaultAdmin):
		c.JSON(http.StatusForbidden, Response{
			Code:    "cannot_delete_default_admin",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, Response{
			Code:    "invalid_credentials",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrWrongPassword):
		c.JSON(http.StatusBadRequest, Response{
			Code:    "wrong_password",
			Message: err.Error(),
		})

	// DNS 记录相关错误
	case errors.Is(err, apperrors.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, Response{
			Code:    "record_not_found",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrRecordExists):
		c.JSON(http.StatusConflict, Response{
			Code:    "record_exists",
			Message: err.Error(),
		})

	// Zone 相关错误
	case errors.Is(err, apperrors.ErrZoneNotFound):
		c.JSON(http.StatusNotFound, Response{
			Code:    "zone_not_found",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrZoneExists):
		c.JSON(http.StatusConflict, Response{
			Code:    "zone_exists",
			Message: err.Error(),
		})

	// Domain 相关错误
	case errors.Is(err, apperrors.ErrDomainNotFound):
		c.JSON(http.StatusNotFound, Response{
			Code:    "domain_not_found",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrDomainExists):
		c.JSON(http.StatusConflict, Response{
			Code:    "domain_exists",
			Message: err.Error(),
		})

	// 认证授权错误
	case errors.Is(err, apperrors.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, Response{
			Code:    "invalid_token",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrTokenExpired):
		c.JSON(http.StatusUnauthorized, Response{
			Code:    "token_expired",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, Response{
			Code:    "unauthorized",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrForbidden):
		c.JSON(http.StatusForbidden, Response{
			Code:    "forbidden",
			Message: err.Error(),
		})
	case errors.Is(err, apperrors.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, Response{
			Code:    "invalid_input",
			Message: err.Error(),
		})

	// 未知错误
	default:
		logger.Log.WithError(err).Error("Unhandled error")
		c.JSON(http.StatusInternalServerError, Response{
			Code:    "internal_error",
			Message: "internal server error",
		})
	}
}

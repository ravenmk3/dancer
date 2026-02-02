package router

import (
	"net/http"
	"time"

	"dancer/internal/auth"
	"dancer/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(
	userHandler *handlers.UserHandler,
	dnsHandler *handlers.DNSHandler,
) *echo.Echo {
	e := echo.New()

	// 中间件
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	// 健康检查
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// API 路由组
	api := e.Group("/api")

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

	// DNS 记录管理（需要认证）
	dns := api.Group("/dns/records", auth.JWTMiddleware())
	dns.POST("/list", dnsHandler.ListRecords)
	dns.POST("/create", dnsHandler.CreateRecord)
	dns.POST("/update", dnsHandler.UpdateRecord)
	dns.POST("/delete", dnsHandler.DeleteRecord)

	return e
}

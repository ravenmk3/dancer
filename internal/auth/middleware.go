package auth

import (
	"net/http"
	"strings"

	"dancer/internal/errors"
	"dancer/internal/models"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware JWT认证中间件
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.ErrUnauthorized)
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, errors.ErrInvalidToken)
			}

			claims, err := ValidateToken(parts[1])
			if err != nil {
				if err == errors.ErrTokenExpired {
					return echo.NewHTTPError(http.StatusUnauthorized, errors.ErrTokenExpired)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, errors.ErrInvalidToken)
			}

			// 将用户信息存入上下文
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("user_type", claims.UserType)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

// RequireAdmin 管理员权限检查中间件
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userType := c.Get("user_type")
			if userType == nil || userType.(string) != string(models.UserTypeAdmin) {
				return echo.NewHTTPError(http.StatusForbidden, errors.ErrForbidden)
			}
			return next(c)
		}
	}
}

// GetCurrentUser 从上下文获取当前用户信息
func GetCurrentUser(c echo.Context) *models.CurrentUser {
	return &models.CurrentUser{
		ID:       c.Get("user_id").(string),
		Username: c.Get("username").(string),
		UserType: models.UserType(c.Get("user_type").(string)),
	}
}

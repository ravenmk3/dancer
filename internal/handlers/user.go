package handlers

import (
	"net/http"

	"dancer/internal/auth"
	"dancer/internal/errors"
	"dancer/internal/logger"
	"dancer/internal/models"
	"dancer/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *services.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// Login 用户登录
func (h *UserHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, CodeInvalidInput, "invalid request body", http.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
	}

	token, user, err := h.userService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		if err == errors.ErrInvalidCredentials {
			return Error(c, CodeInvalidCredentials, err.Error(), http.StatusUnauthorized)
		}
		logger.Log.WithError(err).Error("Login failed")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return Success(c, &models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// RefreshToken 刷新Token
func (h *UserHandler) RefreshToken(c echo.Context) error {
	currentUser := auth.GetCurrentUser(c)

	token, err := auth.GenerateToken(currentUser.ID, currentUser.Username, string(currentUser.UserType))
	if err != nil {
		logger.Log.WithError(err).Error("Failed to refresh token")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return Success(c, map[string]string{"token": token})
}

// GetCurrentUser 获取当前用户信息
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	currentUser := auth.GetCurrentUser(c)

	user, err := h.userService.GetCurrentUser(c.Request().Context(), currentUser.ID)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return Error(c, CodeUserNotFound, err.Error(), http.StatusNotFound)
		}
		logger.Log.WithError(err).Error("Failed to get current user")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return Success(c, user)
}

// ChangePassword 修改当前用户密码
func (h *UserHandler) ChangePassword(c echo.Context) error {
	var req models.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, CodeInvalidInput, "invalid request body", http.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
	}

	currentUser := auth.GetCurrentUser(c)

	if err := h.userService.ChangePassword(c.Request().Context(), currentUser.ID, req.OldPassword, req.NewPassword); err != nil {
		if err == errors.ErrWrongPassword {
			return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
		}
		if err == errors.ErrUserNotFound {
			return Error(c, CodeUserNotFound, err.Error(), http.StatusNotFound)
		}
		logger.Log.WithError(err).Error("Failed to change password")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return SuccessWithMessage(c, "password changed successfully", nil)
}

// ListUsers 列出所有用户（Admin）
func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userService.ListUsers(c.Request().Context())
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list users")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return Success(c, &models.UserListResponse{Users: users})
}

// CreateUser 创建用户（Admin）
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, CodeInvalidInput, "invalid request body", http.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
	}

	user, err := h.userService.CreateUser(c.Request().Context(), &req)
	if err != nil {
		if err == errors.ErrUserExists {
			return Error(c, CodeUserExists, err.Error(), http.StatusConflict)
		}
		logger.Log.WithError(err).Error("Failed to create user")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return Success(c, user)
}

// UpdateUser 更新用户（Admin）
func (h *UserHandler) UpdateUser(c echo.Context) error {
	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, CodeInvalidInput, "invalid request body", http.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
	}

	if err := h.userService.UpdateUser(c.Request().Context(), &req); err != nil {
		if err == errors.ErrUserNotFound {
			return Error(c, CodeUserNotFound, err.Error(), http.StatusNotFound)
		}
		if err == errors.ErrUserExists {
			return Error(c, CodeUserExists, err.Error(), http.StatusConflict)
		}
		logger.Log.WithError(err).Error("Failed to update user")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return SuccessWithMessage(c, "user updated successfully", nil)
}

// DeleteUser 删除用户（Admin）
func (h *UserHandler) DeleteUser(c echo.Context) error {
	var req models.DeleteUserRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, CodeInvalidInput, "invalid request body", http.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
	}

	if err := h.userService.DeleteUser(c.Request().Context(), req.ID); err != nil {
		if err == errors.ErrUserNotFound {
			return Error(c, CodeUserNotFound, err.Error(), http.StatusNotFound)
		}
		if err.Error() == "cannot delete default admin user" {
			return Error(c, CodeForbidden, err.Error(), http.StatusForbidden)
		}
		logger.Log.WithError(err).Error("Failed to delete user")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return SuccessWithMessage(c, "user deleted successfully", nil)
}

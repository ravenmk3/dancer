package handlers

import (
	"dancer/internal/auth"
	apperrors "dancer/internal/errors"
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

// toUserDTO 将 User 转换为 UserDTO（排除 password 字段）
func toUserDTO(user *models.User) *models.UserDTO {
	return &models.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Login 用户登录
func (h *UserHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidInput
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.ErrInvalidInput
	}

	token, _, err := h.userService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(200, &models.LoginResponse{
		Token: token,
	})
}

// RefreshToken 刷新Token
func (h *UserHandler) RefreshToken(c echo.Context) error {
	currentUser := auth.GetCurrentUser(c)

	token, err := auth.GenerateToken(currentUser.ID, currentUser.Username, string(currentUser.UserType))
	if err != nil {
		logger.Log.WithError(err).Error("Failed to refresh token")
		return err
	}

	return c.JSON(200, &models.Response{
		Code:    "success",
		Message: "success",
		Data:    &models.LoginResponse{Token: token},
	})
}

// GetCurrentUser 获取当前用户信息
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	currentUser := auth.GetCurrentUser(c)

	user, err := h.userService.GetCurrentUser(c.Request().Context(), currentUser.ID)
	if err != nil {
		return err
	}

	return c.JSON(200, toUserDTO(user))
}

// ChangePassword 修改当前用户密码
func (h *UserHandler) ChangePassword(c echo.Context) error {
	var req models.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidInput
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.ErrInvalidInput
	}

	currentUser := auth.GetCurrentUser(c)

	if err := h.userService.ChangePassword(c.Request().Context(), currentUser.ID, req.OldPassword, req.NewPassword); err != nil {
		return err
	}

	return c.JSON(200, &models.Response{
		Code:    "success",
		Message: "password changed successfully",
	})
}

// ListUsers 列出所有用户（Admin）
func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userService.ListUsers(c.Request().Context())
	if err != nil {
		return err
	}

	// 转换用户列表，排除 password 字段
	responses := make([]*models.UserDTO, len(users))
	for i, user := range users {
		responses[i] = toUserDTO(user)
	}

	return c.JSON(200, &models.UserListDTO{Users: responses})
}

// CreateUser 创建用户（Admin）
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidInput
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.ErrInvalidInput
	}

	user, err := h.userService.CreateUser(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(200, toUserDTO(user))
}

// UpdateUser 更新用户（Admin）
func (h *UserHandler) UpdateUser(c echo.Context) error {
	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidInput
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.ErrInvalidInput
	}

	if err := h.userService.UpdateUser(c.Request().Context(), &req); err != nil {
		return err
	}

	return c.JSON(200, &models.Response{
		Code:    "success",
		Message: "user updated successfully",
	})
}

// DeleteUser 删除用户（Admin）
func (h *UserHandler) DeleteUser(c echo.Context) error {
	var req models.DeleteUserRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidInput
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.ErrInvalidInput
	}

	if err := h.userService.DeleteUser(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return c.JSON(200, &models.Response{
		Code:    "success",
		Message: "user deleted successfully",
	})
}

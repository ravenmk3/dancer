package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dancer/internal/auth"
	apperrors "dancer/internal/errors"
	"dancer/internal/logger"
	"dancer/internal/models"
	"dancer/internal/storage/etcd"
)

type UserService struct {
	userStorage *etcd.UserStorage
}

func NewUserService(userStorage *etcd.UserStorage) *UserService {
	return &UserService{userStorage: userStorage}
}

// InitDefaultAdmin 初始化默认管理员账户
func (s *UserService) InitDefaultAdmin(ctx context.Context) error {
	// 检查是否已有用户
	count, err := s.userStorage.CountUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	// 如果已有用户，不创建默认管理员
	if count > 0 {
		logger.Log.Info("Users already exist, skipping default admin creation")
		return nil
	}

	// 创建默认管理员
	hashedPassword, err := auth.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("failed to hash default admin password: %w", err)
	}

	admin := &models.User{
		ID:        "10000",
		Username:  "admin",
		Password:  hashedPassword,
		UserType:  models.UserTypeAdmin,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := s.userStorage.CreateUser(ctx, admin); err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	logger.Log.Info("Default admin user created successfully (ID: 10000)")
	return nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, password string) (string, *models.User, error) {
	user, err := s.userStorage.GetUserByUsername(ctx, username)
	if err != nil {
		if err == apperrors.ErrUserNotFound {
			return "", nil, apperrors.ErrInvalidCredentials
		}
		return "", nil, err
	}

	if !auth.CheckPassword(password, user.Password) {
		return "", nil, apperrors.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, user.Username, string(user.UserType))
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser(ctx context.Context, userID string) (*models.User, error) {
	return s.userStorage.GetUser(ctx, userID)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.userStorage.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	if !auth.CheckPassword(oldPassword, user.Password) {
		return apperrors.ErrWrongPassword
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	user.UpdatedAt = time.Now().Unix()

	return s.userStorage.UpdateUser(ctx, user)
}

// ListUsers 列出所有用户
func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.userStorage.ListUsers(ctx)
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// 检查用户名是否已存在
	_, err := s.userStorage.GetUserByUsername(ctx, req.Username)
	if err == nil {
		return nil, apperrors.ErrUserExists
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:        fmt.Sprintf("%d", time.Now().UnixMilli()),
		Username:  req.Username,
		Password:  hashedPassword,
		UserType:  req.UserType,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := s.userStorage.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, req *models.UpdateUserRequest) error {
	user, err := s.userStorage.GetUser(ctx, req.ID)
	if err != nil {
		return err
	}

	// 如果修改了用户名，检查是否已存在
	if req.Username != "" && req.Username != user.Username {
		existing, _ := s.userStorage.GetUserByUsername(ctx, req.Username)
		if existing != nil {
			return apperrors.ErrUserExists
		}
		user.Username = req.Username
	}

	// 如果修改了密码
	if req.Password != "" {
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	// 如果修改了用户类型
	if req.UserType != "" {
		user.UserType = req.UserType
	}

	user.UpdatedAt = time.Now().Unix()

	return s.userStorage.UpdateUser(ctx, user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	// 检查是否是默认管理员
	if userID == "10000" {
		return errors.New("cannot delete default admin user")
	}

	return s.userStorage.DeleteUser(ctx, userID)
}

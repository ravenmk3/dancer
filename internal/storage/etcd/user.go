package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dancer/internal/errors"
	"dancer/internal/models"
	"dancer/internal/storage"
	"go.etcd.io/etcd/client/v3"
)

// UserStorage 用户存储操作
type UserStorage struct {
	client *Client
}

// defaultWaitTimeout 默认等待超时时间
const defaultWaitTimeout = 5 * time.Second

func NewUserStorage(client *Client) *UserStorage {
	return &UserStorage{client: client}
}

// checkConnection 检查 etcd 连接，使用默认超时
func (s *UserStorage) checkConnection() error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}
	return nil
}

// GetUser 根据ID获取用户
func (s *UserStorage) GetUser(ctx context.Context, id string) (*models.User, error) {
	if err := s.checkConnection(); err != nil {
		return nil, err
	}

	key := storage.UserKeyPrefix + id
	resp, err := s.client.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, errors.ErrUserNotFound
	}

	var user models.User
	if err := json.Unmarshal(resp.Kvs[0].Value, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户（使用范围查询）
func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	if err := s.checkConnection(); err != nil {
		return nil, err
	}

	// 获取所有用户并筛选
	resp, err := s.client.client.Get(ctx, storage.UserKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range resp.Kvs {
		var user models.User
		if err := json.Unmarshal(kv.Value, &user); err != nil {
			continue
		}
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, errors.ErrUserNotFound
}

// ListUsers 列出所有用户
func (s *UserStorage) ListUsers(ctx context.Context) ([]*models.User, error) {
	if err := s.checkConnection(); err != nil {
		return nil, err
	}

	resp, err := s.client.client.Get(ctx, storage.UserKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var user models.User
		if err := json.Unmarshal(kv.Value, &user); err != nil {
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

// CreateUser 创建用户
func (s *UserStorage) CreateUser(ctx context.Context, user *models.User) error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	key := storage.UserKeyPrefix + user.ID

	// 检查用户是否已存在（通过用户名）
	_, err := s.GetUserByUsername(ctx, user.Username)
	if err == nil {
		return errors.ErrUserExists
	}

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	_, err = s.client.client.Put(ctx, key, string(data))
	return err
}

// UpdateUser 更新用户
func (s *UserStorage) UpdateUser(ctx context.Context, user *models.User) error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	key := storage.UserKeyPrefix + user.ID

	// 检查用户是否存在
	_, err := s.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	_, err = s.client.client.Put(ctx, key, string(data))
	return err
}

// DeleteUser 删除用户
func (s *UserStorage) DeleteUser(ctx context.Context, id string) error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	key := storage.UserKeyPrefix + id
	_, err := s.client.client.Delete(ctx, key)
	return err
}

// CountUsers 统计用户数量
func (s *UserStorage) CountUsers(ctx context.Context) (int64, error) {
	if err := s.checkConnection(); err != nil {
		return 0, err
	}

	resp, err := s.client.client.Get(ctx, storage.UserKeyPrefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, err
	}
	return resp.Count, nil
}

// IsUserExists 检查用户是否存在（通过ID）
func (s *UserStorage) IsUserExists(ctx context.Context, id string) bool {
	_, err := s.GetUser(ctx, id)
	return err == nil
}

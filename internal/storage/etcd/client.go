package etcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"dancer/internal/config"
	"dancer/internal/logger"
	"go.etcd.io/etcd/client/v3"
)

// ConnectionState 连接状态
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
)

// Client etcd 客户端包装器，支持自动重连
type Client struct {
	client    *clientv3.Client
	config    *config.Config
	state     ConnectionState
	stateMu   sync.RWMutex
	stopCh    chan struct{}
	connectCh chan struct{}
}

// NewClient 创建 etcd 客户端（异步初始化，允许启动时无连接）
func NewClient(cfg *config.Config) (*Client, error) {
	c := &Client{
		config:    cfg,
		state:     StateDisconnected,
		stopCh:    make(chan struct{}),
		connectCh: make(chan struct{}, 1),
	}

	// 异步尝试首次连接
	go c.initialConnection()

	// 启动后台连接 goroutine
	go c.connectLoop()

	// 启动健康检查
	go c.healthCheckLoop()

	return c, nil
}

// initialConnection 首次连接尝试
func (c *Client) initialConnection() {
	if err := c.tryConnect(); err != nil {
		logger.Log.WithError(err).Warn("Initial etcd connection failed, will retry in background")
	} else {
		logger.Log.Info("Initial etcd connection established")
	}
}

// getDialTimeout 获取连接超时时间
func (c *Client) getDialTimeout() time.Duration {
	if c.config.Etcd.DialTimeout > 0 {
		return time.Duration(c.config.Etcd.DialTimeout) * time.Second
	}
	return 5 * time.Second
}

// getConnectInterval 获取初始连接间隔
func (c *Client) getConnectInterval() time.Duration {
	if c.config.Etcd.ReconnectInterval > 0 {
		return time.Duration(c.config.Etcd.ReconnectInterval) * time.Second
	}
	return 5 * time.Second
}

// getMaxConnectInterval 获取最大连接间隔
func (c *Client) getMaxConnectInterval() time.Duration {
	if c.config.Etcd.MaxReconnectInterval > 0 {
		return time.Duration(c.config.Etcd.MaxReconnectInterval) * time.Second
	}
	return 30 * time.Second
}

// getHealthCheckInterval 获取健康检查间隔
func (c *Client) getHealthCheckInterval() time.Duration {
	if c.config.Etcd.HealthCheckInterval > 0 {
		return time.Duration(c.config.Etcd.HealthCheckInterval) * time.Second
	}
	return 30 * time.Second
}

// tryConnect 尝试连接 etcd
func (c *Client) tryConnect() error {
	c.setState(StateConnecting)

	etcdCfg := clientv3.Config{
		Endpoints:   c.config.Etcd.Endpoints,
		DialTimeout: c.getDialTimeout(),
	}

	if c.config.Etcd.Username != "" && c.config.Etcd.Password != "" {
		etcdCfg.Username = c.config.Etcd.Username
		etcdCfg.Password = c.config.Etcd.Password
	}

	client, err := clientv3.New(etcdCfg)
	if err != nil {
		c.setState(StateDisconnected)
		return fmt.Errorf("failed to create etcd client: %w", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), c.getDialTimeout())
	defer cancel()

	if _, err := client.Status(ctx, c.config.Etcd.Endpoints[0]); err != nil {
		client.Close()
		c.setState(StateDisconnected)
		return fmt.Errorf("failed to connect to etcd: %w", err)
	}

	c.client = client
	c.setState(StateConnected)
	return nil
}

// connectLoop 后台连接循环（指数退避）
func (c *Client) connectLoop() {
	baseInterval := c.getConnectInterval()
	maxInterval := c.getMaxConnectInterval()
	currentInterval := baseInterval

	for {
		select {
		case <-c.stopCh:
			return
		case <-c.connectCh:
			if c.IsConnected() {
				continue
			}

			logger.Log.Infof("Attempting to connect to etcd (interval: %v)", currentInterval)

			if err := c.tryConnect(); err != nil {
				logger.Log.WithError(err).Error("Etcd connection failed")
				// 指数退避
				currentInterval *= 2
				if currentInterval > maxInterval {
					currentInterval = maxInterval
				}
				// 调度下一次连接
				go func() {
					time.Sleep(currentInterval)
					c.triggerConnect()
				}()
			} else {
				logger.Log.Info("Etcd connection successful")
				currentInterval = baseInterval // 重置间隔
			}
		}
	}
}

// healthCheckLoop 健康检查循环
func (c *Client) healthCheckLoop() {
	ticker := time.NewTicker(c.getHealthCheckInterval())
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			if !c.IsConnected() {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), c.getDialTimeout())
			_, err := c.client.Status(ctx, c.config.Etcd.Endpoints[0])
			cancel()

			if err != nil {
				logger.Log.WithError(err).Error("Etcd health check failed")
				c.setState(StateDisconnected)
				c.client.Close()
				c.client = nil
				c.triggerConnect()
			}
		}
	}
}

// triggerConnect 触发连接
func (c *Client) triggerConnect() {
	select {
	case c.connectCh <- struct{}{}:
	default:
	}
}

// WaitForConnection 等待连接就绪
// timeout: 0 表示不等待，立即检查
//
//	>0 表示最多等待指定时间
func (c *Client) WaitForConnection(timeout time.Duration) error {
	if c.IsConnected() {
		return nil
	}

	if timeout == 0 {
		return fmt.Errorf("etcd not connected")
	}

	// 触发连接
	c.triggerConnect()

	// 等待连接或超时
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for etcd connection")
		case <-ticker.C:
			if c.IsConnected() {
				return nil
			}
		}
	}
}

// IsConnected 检查是否已连接
func (c *Client) IsConnected() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state == StateConnected && c.client != nil
}

// GetState 获取当前连接状态
func (c *Client) GetState() ConnectionState {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state
}

func (c *Client) setState(state ConnectionState) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.state = state
}

// GetClient 获取底层 etcd 客户端
// 注意：调用前应先检查 IsConnected() 或使用 WaitForConnection()
func (c *Client) GetClient() *clientv3.Client {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.client
}

// Close 关闭客户端
func (c *Client) Close() error {
	close(c.stopCh)

	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// SetupTLS 配置TLS（可选）
func SetupTLS(certFile, keyFile, caFile string) (*tls.Config, error) {
	return &tls.Config{}, nil
}

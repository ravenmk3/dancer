package etcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"dancer/internal/config"
	"go.etcd.io/etcd/client/v3"
)

type Client struct {
	client *clientv3.Client
}

func NewClient(cfg *config.Config) (*Client, error) {
	etcdCfg := clientv3.Config{
		Endpoints:   cfg.Etcd.Endpoints,
		DialTimeout: 5 * time.Second,
	}

	if cfg.Etcd.Username != "" && cfg.Etcd.Password != "" {
		etcdCfg.Username = cfg.Etcd.Username
		etcdCfg.Password = cfg.Etcd.Password
	}

	client, err := clientv3.New(etcdCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Status(ctx, cfg.Etcd.Endpoints[0]); err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) GetClient() *clientv3.Client {
	return c.client
}

// SetupTLS 配置TLS（可选）
func SetupTLS(certFile, keyFile, caFile string) (*tls.Config, error) {
	// 简化实现，实际需要加载证书
	return &tls.Config{}, nil
}

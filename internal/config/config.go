package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值
	if cfg.App.Host == "" {
		cfg.App.Host = "0.0.0.0"
	}
	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}
	if cfg.App.Env == "" {
		cfg.App.Env = "development"
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "your-256-bit-secret-change-in-production"
	}
	if cfg.JWT.Expiry == 0 {
		cfg.JWT.Expiry = 86400 // 24小时
	}
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}
	if cfg.Logger.FilePath == "" {
		cfg.Logger.FilePath = "logs/dancer.log"
	}
	if cfg.Logger.MaxSize == 0 {
		cfg.Logger.MaxSize = 100
	}
	if cfg.Logger.MaxBackup == 0 {
		cfg.Logger.MaxBackup = 7
	}
	if cfg.Logger.MaxAge == 0 {
		cfg.Logger.MaxAge = 7
	}
	if cfg.Etcd.CorednsPrefix == "" {
		cfg.Etcd.CorednsPrefix = "/skydns"
	}

	GlobalConfig = &cfg
	return nil
}

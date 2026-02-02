package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"dancer/internal/config"
	"dancer/internal/handlers"
	"dancer/internal/logger"
	"dancer/internal/router"
	"dancer/internal/services"
	"dancer/internal/storage/etcd"
)

func main() {
	// 命令行参数
	configPath := flag.String("config", "config.toml", "配置文件路径")
	flag.Parse()

	// 加载配置
	if err := config.Load(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}
	cfg := config.GetConfig()

	// 初始化日志
	if err := logger.Init(
		cfg.Logger.Level,
		cfg.Logger.FilePath,
		cfg.Logger.MaxSize,
		cfg.Logger.MaxBackup,
		cfg.Logger.MaxAge,
	); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Log.Info("Starting Dancer DNS Management Tool")

	// 初始化 etcd 客户端
	etcdClient, err := etcd.NewClient(cfg)
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to connect to etcd")
	}
	defer etcdClient.Close()
	logger.Log.Info("Connected to etcd")

	// 初始化存储层
	userStorage := etcd.NewUserStorage(etcdClient)
	dnsStorage := etcd.NewDNSStorage(etcdClient)

	// 初始化服务层
	userService := services.NewUserService(userStorage)
	dnsService := services.NewDNSService(dnsStorage)

	// 初始化默认管理员
	ctx := context.Background()
	if err := userService.InitDefaultAdmin(ctx); err != nil {
		logger.Log.WithError(err).Fatal("Failed to initialize default admin")
	}

	// 初始化处理器
	userHandler := handlers.NewUserHandler(userService)
	dnsHandler := handlers.NewDNSHandler(dnsService)

	// 初始化路由
	e := router.New(userHandler, dnsHandler)

	// 启动服务器
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
		logger.Log.Infof("Server starting on %s", addr)
		if err := e.Start(addr); err != nil {
			logger.Log.WithError(err).Info("Server stopped")
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")
	if err := e.Shutdown(ctx); err != nil {
		logger.Log.WithError(err).Error("Server forced to shutdown")
	}
	logger.Log.Info("Server exited")
}

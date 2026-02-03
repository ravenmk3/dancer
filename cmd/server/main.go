package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dancer/internal/config"
	"dancer/internal/handlers"
	"dancer/internal/logger"
	"dancer/internal/router"
	"dancer/internal/services"
	"dancer/internal/storage/etcd"
)

func main() {
	// 打印青色 ASCII Logo（最先显示）
	fmt.Println("\033[36m")
	fmt.Println("    ██████╗  █████╗ ███╗   ██╗ ██████╗███████╗██████╗ ")
	fmt.Println("    ██╔══██╗██╔══██╗████╗  ██║██╔════╝██╔════╝██╔══██╗")
	fmt.Println("    ██║  ██║███████║██╔██╗ ██║██║     █████╗  ██████╔╝")
	fmt.Println("    ██║  ██║██╔══██║██║╚██╗██║██║     ██╔══╝  ██╔══██╗")
	fmt.Println("    ██████╔╝██║  ██║██║ ╚████║╚██████╗███████╗██║  ██║")
	fmt.Println("    ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝╚══════╝╚═╝  ╚═╝")
	fmt.Println("\033[0m")
	fmt.Println("         DNS Management Tool")
	fmt.Println()

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

	// 初始化 etcd 客户端（允许启动时无连接）
	etcdClient, err := etcd.NewClient(cfg)
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to initialize etcd client")
	}
	defer etcdClient.Close()

	if etcdClient.IsConnected() {
		logger.Log.Info("Connected to etcd")
	} else {
		logger.Log.Warn("Etcd not connected, will retry in background")
	}

	// 初始化存储层
	userStorage := etcd.NewUserStorage(etcdClient)
	dnsStorage := etcd.NewDNSStorage(etcdClient)

	// 初始化服务层
	userService := services.NewUserService(userStorage)
	dnsService := services.NewDNSService(dnsStorage)

	// 初始化默认管理员（在后台 goroutine 中执行，避免阻塞启动）
	go func() {
		// 等待 etcd 连接就绪
		if err := etcdClient.WaitForConnection(30 * time.Second); err != nil {
			logger.Log.WithError(err).Error("Failed to wait for etcd connection, skipping default admin initialization")
			return
		}

		ctx := context.Background()
		if err := userService.InitDefaultAdmin(ctx); err != nil {
			logger.Log.WithError(err).Error("Failed to initialize default admin")
		}
	}()

	// 初始化处理器
	userHandler := handlers.NewUserHandler(userService)
	dnsHandler := handlers.NewDNSHandler(dnsService)
	healthHandler := handlers.NewHealthHandler(etcdClient)

	// 初始化路由
	e := router.New(userHandler, dnsHandler, healthHandler)

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
	shutdownCtx := context.Background()
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Log.WithError(err).Error("Server forced to shutdown")
	}
	logger.Log.Info("Server exited")
}

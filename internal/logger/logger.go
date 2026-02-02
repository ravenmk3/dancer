package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func Init(level, filePath string, maxSize, maxBackup, maxAge int) error {
	Log = logrus.New()

	// 设置日志级别
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	Log.SetLevel(logLevel)

	// 确保日志目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 配置文件轮转
	fileHook := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize, // MB
		MaxBackups: maxBackup,
		MaxAge:     maxAge, // days
		Compress:   true,
	}

	// 控制台输出使用彩色格式
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 同时输出到文件和控制台
	Log.SetOutput(os.Stdout)

	// 添加文件钩子
	Log.AddHook(&fileHookAdapter{writer: fileHook})

	Log.Info("Logger initialized successfully")
	return nil
}

type fileHookAdapter struct {
	writer *lumberjack.Logger
}

func (h *fileHookAdapter) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *fileHookAdapter) Fire(entry *logrus.Entry) error {
	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	line, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(line)
	return err
}

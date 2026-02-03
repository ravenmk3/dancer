package router

import (
	"fmt"
	"time"

	"dancer/internal/logger"
	"github.com/labstack/echo/v4"
)

// CustomLogger 自定义访问日志中间件
// 使用 logrus 输出 DEBUG 级别的访问日志
// 格式: DEBU[2026-02-03 23:26:42] 127.0.0.1 | GET /api/health | 200 | 0ms | 0B/43B
func CustomLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// 执行请求
			err := next(c)

			// 计算延迟（转换为毫秒）
			latencyMs := time.Since(start).Milliseconds()

			// 获取请求信息
			req := c.Request()
			res := c.Response()

			// 格式化访问日志
			// 格式: {remote_ip} | {method} {uri} | {status} | {latency}ms | {bytes_in}B/{bytes_out}B
			logger.Log.Debug(fmt.Sprintf("%s | %s %s | %d | %dms | %dB/%dB",
				c.RealIP(),
				req.Method,
				req.RequestURI,
				res.Status,
				latencyMs,
				req.ContentLength,
				res.Size,
			))

			return err
		}
	}
}

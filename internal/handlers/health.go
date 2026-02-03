package handlers

import (
	"net/http"

	"dancer/internal/storage/etcd"
	"github.com/labstack/echo/v4"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	etcdClient *etcd.Client
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(etcdClient *etcd.Client) *HealthHandler {
	return &HealthHandler{
		etcdClient: etcdClient,
	}
}

// Check 健康检查
// 所有 components 为 up 时返回 200，任一 component 为 down 时返回 503
func (h *HealthHandler) Check(c echo.Context) error {
	// 检查 etcd 状态
	etcdStatus := "up"
	if !h.etcdClient.IsConnected() {
		etcdStatus = "down"
	}

	// 计算总状态
	overallStatus := "up"
	httpStatus := http.StatusOK

	if etcdStatus == "down" {
		overallStatus = "down"
		httpStatus = http.StatusServiceUnavailable
	}

	return c.JSON(httpStatus, map[string]interface{}{
		"status": overallStatus,
		"components": map[string]string{
			"etcd": etcdStatus,
		},
	})
}

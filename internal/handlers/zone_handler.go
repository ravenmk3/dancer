package handlers

import (
	"dancer/internal/logger"
	"dancer/internal/models"
	"dancer/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ZoneHandler Zone HTTP 处理器
type ZoneHandler struct {
	zoneService *services.ZoneService
	validate    *validator.Validate
}

func NewZoneHandler(zoneService *services.ZoneService) *ZoneHandler {
	return &ZoneHandler{
		zoneService: zoneService,
		validate:    validator.New(),
	}
}

// ListZones 列出所有 Zone
func (h *ZoneHandler) ListZones(c echo.Context) error {
	zones, err := h.zoneService.ListZones(c.Request().Context())
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list zones")
		return err
	}

	return c.JSON(200, &models.ZoneListResponse{Zones: zones})
}

// GetZone 获取 Zone 详情
func (h *ZoneHandler) GetZone(c echo.Context) error {
	var req models.GetZoneRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	zone, err := h.zoneService.GetZone(c.Request().Context(), req.Zone)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get zone")
		return err
	}

	return c.JSON(200, &models.ZoneResponse{Zone: zone})
}

// CreateZone 创建 Zone
func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req models.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	zone, err := h.zoneService.CreateZone(c.Request().Context(), &req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create zone")
		return err
	}

	return c.JSON(200, &models.ZoneResponse{Zone: zone})
}

// UpdateZone 更新 Zone
func (h *ZoneHandler) UpdateZone(c echo.Context) error {
	var req models.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	zone, err := h.zoneService.UpdateZone(c.Request().Context(), &req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to update zone")
		return err
	}

	return c.JSON(200, &models.ZoneResponse{Zone: zone})
}

// DeleteZone 删除 Zone
func (h *ZoneHandler) DeleteZone(c echo.Context) error {
	var req models.DeleteZoneRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	if err := h.zoneService.DeleteZone(c.Request().Context(), &req); err != nil {
		logger.Log.WithError(err).Error("Failed to delete zone")
		return err
	}

	return c.JSON(200, map[string]interface{}{
		"code":    "success",
		"message": "Zone deleted successfully",
	})
}

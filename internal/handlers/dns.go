package handlers

import (
	"dancer/internal/logger"
	"dancer/internal/models"
	"dancer/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type DNSHandler struct {
	dnsService *services.DNSService
	validate   *validator.Validate
}

func NewDNSHandler(dnsService *services.DNSService) *DNSHandler {
	return &DNSHandler{
		dnsService: dnsService,
		validate:   validator.New(),
	}
}

// ListRecords 列出DNS记录
func (h *DNSHandler) ListRecords(c echo.Context) error {
	var req models.ListDNSRequest
	if err := c.Bind(&req); err != nil {
		// 如果没有请求体，继续处理（domain为空）
	}

	records, err := h.dnsService.ListRecords(c.Request().Context(), req.Domain)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list DNS records")
		return err
	}

	return c.JSON(200, &models.DNSListResponse{Records: records})
}

// CreateRecord 创建DNS记录
func (h *DNSHandler) CreateRecord(c echo.Context) error {
	var req models.CreateDNSRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	userID := c.Get("user_id").(string)
	record, err := h.dnsService.CreateRecord(c.Request().Context(), &req, userID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create DNS record")
		return err
	}

	return c.JSON(200, &models.DNSResponse{Record: record})
}

// UpdateRecord 更新DNS记录
func (h *DNSHandler) UpdateRecord(c echo.Context) error {
	var req models.UpdateDNSRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	userID := c.Get("user_id").(string)
	record, err := h.dnsService.UpdateRecord(c.Request().Context(), &req, userID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to update DNS record")
		return err
	}

	return c.JSON(200, &models.DNSResponse{Record: record})
}

// DeleteRecord 删除DNS记录
func (h *DNSHandler) DeleteRecord(c echo.Context) error {
	var req models.DeleteDNSRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	if err := h.dnsService.DeleteRecord(c.Request().Context(), req.Key); err != nil {
		logger.Log.WithError(err).Error("Failed to delete DNS record")
		return err
	}

	return c.JSON(200, map[string]interface{}{
		"code":    "success",
		"message": "DNS record deleted successfully",
	})
}

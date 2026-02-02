package handlers

import (
	"net/http"

	"dancer/internal/errors"
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
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return Success(c, &models.DNSListResponse{Records: records})
}

// CreateRecord 创建DNS记录
func (h *DNSHandler) CreateRecord(c echo.Context) error {
	var req models.CreateDNSRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, CodeInvalidInput, "invalid request body", http.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
	}

	userID := c.Get("user_id").(string)
	record, err := h.dnsService.CreateRecord(c.Request().Context(), &req, userID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create DNS record")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return Success(c, &models.DNSResponse{Record: record})
}

// UpdateRecord 更新DNS记录
func (h *DNSHandler) UpdateRecord(c echo.Context) error {
	var req models.UpdateDNSRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, CodeInvalidInput, "invalid request body", http.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
	}

	userID := c.Get("user_id").(string)
	record, err := h.dnsService.UpdateRecord(c.Request().Context(), &req, userID)
	if err != nil {
		if err == errors.ErrRecordNotFound {
			return Error(c, CodeRecordNotFound, err.Error(), http.StatusNotFound)
		}
		logger.Log.WithError(err).Error("Failed to update DNS record")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return Success(c, &models.DNSResponse{Record: record})
}

// DeleteRecord 删除DNS记录
func (h *DNSHandler) DeleteRecord(c echo.Context) error {
	var req models.DeleteDNSRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, CodeInvalidInput, "invalid request body", http.StatusBadRequest)
	}

	if err := h.validate.Struct(req); err != nil {
		return Error(c, CodeInvalidInput, err.Error(), http.StatusBadRequest)
	}

	if err := h.dnsService.DeleteRecord(c.Request().Context(), req.Key); err != nil {
		if err == errors.ErrRecordNotFound {
			return Error(c, CodeRecordNotFound, err.Error(), http.StatusNotFound)
		}
		logger.Log.WithError(err).Error("Failed to delete DNS record")
		return Error(c, CodeInternalError, "internal server error", http.StatusInternalServerError)
	}

	return SuccessWithMessage(c, "DNS record deleted successfully", nil)
}

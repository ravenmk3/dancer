package handlers

import (
	"dancer/internal/logger"
	"dancer/internal/models"
	"dancer/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// DomainHandler Domain HTTP 处理器
type DomainHandler struct {
	domainService *services.DomainService
	validate      *validator.Validate
}

func NewDomainHandler(domainService *services.DomainService) *DomainHandler {
	return &DomainHandler{
		domainService: domainService,
		validate:      validator.New(),
	}
}

// ListDomains 列出 Zone 下所有 Domain
func (h *DomainHandler) ListDomains(c echo.Context) error {
	var req models.ListDomainsRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	domains, err := h.domainService.ListDomains(c.Request().Context(), &req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list domains")
		return err
	}

	return c.JSON(200, &models.DomainListResponse{Domains: domains})
}

// GetDomain 获取 Domain 详情
func (h *DomainHandler) GetDomain(c echo.Context) error {
	var req models.GetDomainRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	domain, err := h.domainService.GetDomain(c.Request().Context(), &req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get domain")
		return err
	}

	return c.JSON(200, &models.DomainResponse{Domain: domain})
}

// CreateDomain 创建 Domain
func (h *DomainHandler) CreateDomain(c echo.Context) error {
	var req models.CreateDomainRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	domain, err := h.domainService.CreateDomain(c.Request().Context(), &req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create domain")
		return err
	}

	return c.JSON(200, &models.DomainResponse{Domain: domain})
}

// UpdateDomain 更新 Domain
func (h *DomainHandler) UpdateDomain(c echo.Context) error {
	var req models.UpdateDomainRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	domain, err := h.domainService.UpdateDomain(c.Request().Context(), &req)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to update domain")
		return err
	}

	return c.JSON(200, &models.DomainResponse{Domain: domain})
}

// DeleteDomain 删除 Domain
func (h *DomainHandler) DeleteDomain(c echo.Context) error {
	var req models.DeleteDomainRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.validate.Struct(req); err != nil {
		return err
	}

	if err := h.domainService.DeleteDomain(c.Request().Context(), &req); err != nil {
		logger.Log.WithError(err).Error("Failed to delete domain")
		return err
	}

	return c.JSON(200, &models.Response{
		Code:    "success",
		Message: "Domain deleted successfully",
	})
}

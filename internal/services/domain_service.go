package services

import (
	"context"
	"time"

	"dancer/internal/errors"
	"dancer/internal/models"
	"dancer/internal/storage/etcd"
)

// DomainService Domain 业务逻辑
type DomainService struct {
	zoneStorage   *etcd.ZoneStorage
	domainStorage *etcd.DomainStorage
}

func NewDomainService(zoneStorage *etcd.ZoneStorage, domainStorage *etcd.DomainStorage) *DomainService {
	return &DomainService{
		zoneStorage:   zoneStorage,
		domainStorage: domainStorage,
	}
}

// ListDomains 列出 Zone 下所有 Domain
func (s *DomainService) ListDomains(ctx context.Context, req *models.ListDomainsRequest) ([]*models.Domain, error) {
	// 检查 Zone 是否存在
	_, err := s.zoneStorage.GetZone(ctx, req.Zone)
	if err != nil {
		return nil, err
	}

	return s.domainStorage.ListDomainsByZone(ctx, req.Zone)
}

// GetDomain 获取 Domain 详情
func (s *DomainService) GetDomain(ctx context.Context, req *models.GetDomainRequest) (*models.Domain, error) {
	// 检查 Zone 是否存在
	_, err := s.zoneStorage.GetZone(ctx, req.Zone)
	if err != nil {
		return nil, err
	}

	return s.domainStorage.GetDomain(ctx, req.Zone, req.Domain)
}

// CreateDomain 创建 Domain
func (s *DomainService) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*models.Domain, error) {
	// 检查 Zone 是否存在
	_, err := s.zoneStorage.GetZone(ctx, req.Zone)
	if err != nil {
		return nil, errors.ErrZoneNotFound
	}

	// 检查是否已存在
	exists, err := s.domainStorage.DomainExists(ctx, req.Zone, req.Domain)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrDomainExists
	}

	domain := &models.Domain{
		Zone:   req.Zone,
		Domain: req.Domain,
		IPs:    req.IPs,
		TTL:    req.TTL,
	}

	if err := s.domainStorage.CreateDomain(ctx, domain); err != nil {
		return nil, err
	}

	// 更新 Zone 记录数
	count, _ := s.domainStorage.GetDomainCountByZone(ctx, req.Zone)
	s.zoneStorage.UpdateZoneRecordCount(ctx, req.Zone, count)

	return domain, nil
}

// UpdateDomain 更新 Domain
func (s *DomainService) UpdateDomain(ctx context.Context, req *models.UpdateDomainRequest) (*models.Domain, error) {
	// 检查 Zone 是否存在
	_, err := s.zoneStorage.GetZone(ctx, req.Zone)
	if err != nil {
		return nil, errors.ErrZoneNotFound
	}

	// 获取现有记录
	existing, err := s.domainStorage.GetDomain(ctx, req.Zone, req.Domain)
	if err != nil {
		return nil, err
	}

	// 更新字段
	existing.IPs = req.IPs
	if req.TTL > 0 {
		existing.TTL = req.TTL
	}
	existing.UpdatedAt = time.Now().Unix()

	if err := s.domainStorage.UpdateDomain(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteDomain 删除 Domain
func (s *DomainService) DeleteDomain(ctx context.Context, req *models.DeleteDomainRequest) error {
	// 检查 Zone 是否存在
	_, err := s.zoneStorage.GetZone(ctx, req.Zone)
	if err != nil {
		return errors.ErrZoneNotFound
	}

	// 检查 Domain 是否存在
	_, err = s.domainStorage.GetDomain(ctx, req.Zone, req.Domain)
	if err != nil {
		return err
	}

	// 删除 Domain
	if err := s.domainStorage.DeleteDomain(ctx, req.Zone, req.Domain); err != nil {
		return err
	}

	// 更新 Zone 记录数
	count, _ := s.domainStorage.GetDomainCountByZone(ctx, req.Zone)
	s.zoneStorage.UpdateZoneRecordCount(ctx, req.Zone, count)

	return nil
}

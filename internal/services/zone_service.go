package services

import (
	"context"
	"time"

	"dancer/internal/errors"
	"dancer/internal/models"
	"dancer/internal/storage/etcd"
)

// ZoneService Zone 业务逻辑
type ZoneService struct {
	zoneStorage   *etcd.ZoneStorage
	domainStorage *etcd.DomainStorage
}

func NewZoneService(zoneStorage *etcd.ZoneStorage, domainStorage *etcd.DomainStorage) *ZoneService {
	return &ZoneService{
		zoneStorage:   zoneStorage,
		domainStorage: domainStorage,
	}
}

// ListZones 列出所有 Zone
func (s *ZoneService) ListZones(ctx context.Context) ([]*models.Zone, error) {
	return s.zoneStorage.ListZones(ctx)
}

// GetZone 获取 Zone 详情
func (s *ZoneService) GetZone(ctx context.Context, zone string) (*models.Zone, error) {
	return s.zoneStorage.GetZone(ctx, zone)
}

// CreateZone 创建 Zone
func (s *ZoneService) CreateZone(ctx context.Context, req *models.CreateZoneRequest) (*models.Zone, error) {
	// 检查是否已存在
	exists, err := s.zoneStorage.ZoneExists(ctx, req.Zone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrZoneExists
	}

	zone := &models.Zone{
		Zone:        req.Zone,
		RecordCount: 0,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	if err := s.zoneStorage.CreateZone(ctx, zone); err != nil {
		return nil, err
	}

	return zone, nil
}

// UpdateZone 更新 Zone
func (s *ZoneService) UpdateZone(ctx context.Context, req *models.UpdateZoneRequest) (*models.Zone, error) {
	// 检查是否存在
	zone, err := s.zoneStorage.GetZone(ctx, req.Zone)
	if err != nil {
		return nil, err
	}

	zone.UpdatedAt = time.Now().Unix()

	if err := s.zoneStorage.UpdateZone(ctx, zone); err != nil {
		return nil, err
	}

	return zone, nil
}

// DeleteZone 删除 Zone（级联删除所有 Domain）
func (s *ZoneService) DeleteZone(ctx context.Context, req *models.DeleteZoneRequest) error {
	// 检查是否存在
	_, err := s.zoneStorage.GetZone(ctx, req.Zone)
	if err != nil {
		return err
	}

	// 删除 Zone 及所有 Domain
	return s.zoneStorage.DeleteZone(ctx, req.Zone)
}

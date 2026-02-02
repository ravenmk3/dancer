package services

import (
	"context"
	"time"

	"dancer/internal/models"
	"dancer/internal/storage/etcd"
)

type DNSService struct {
	dnsStorage *etcd.DNSStorage
}

func NewDNSService(dnsStorage *etcd.DNSStorage) *DNSService {
	return &DNSService{dnsStorage: dnsStorage}
}

// ListRecords 列出DNS记录
func (s *DNSService) ListRecords(ctx context.Context, domain string) ([]*models.DNSRecord, error) {
	if domain == "" {
		return s.dnsStorage.ListAllRecords(ctx)
	}
	return s.dnsStorage.ListRecords(ctx, domain)
}

// CreateRecord 创建DNS记录
func (s *DNSService) CreateRecord(ctx context.Context, req *models.CreateDNSRequest, userID string) (*models.DNSRecord, error) {
	record := &models.DNSRecord{
		Domain: req.Domain,
		IP:     req.IP,
		TTL:    req.TTL,
	}

	if err := s.dnsStorage.CreateRecord(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

// UpdateRecord 更新DNS记录
func (s *DNSService) UpdateRecord(ctx context.Context, req *models.UpdateDNSRequest, userID string) (*models.DNSRecord, error) {
	existing, err := s.dnsStorage.GetRecord(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Domain != "" {
		existing.Domain = req.Domain
	}
	if req.IP != "" {
		existing.IP = req.IP
	}
	if req.TTL > 0 {
		existing.TTL = req.TTL
	}
	existing.UpdatedAt = time.Now().Unix()

	if err := s.dnsStorage.UpdateRecord(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteRecord 删除DNS记录
func (s *DNSService) DeleteRecord(ctx context.Context, key string) error {
	// 检查记录是否存在
	_, err := s.dnsStorage.GetRecord(ctx, key)
	if err != nil {
		return err
	}

	return s.dnsStorage.DeleteRecord(ctx, key)
}

// GetRecord 获取单个DNS记录
func (s *DNSService) GetRecord(ctx context.Context, key string) (*models.DNSRecord, error) {
	return s.dnsStorage.GetRecord(ctx, key)
}

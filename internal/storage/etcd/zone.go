package etcd

import (
	"context"
	"encoding/json"

	"dancer/internal/errors"
	"dancer/internal/models"
	"dancer/internal/storage"
	"go.etcd.io/etcd/client/v3"
)

// ZoneStorage Zone 存储操作
type ZoneStorage struct {
	client *Client
}

func NewZoneStorage(client *Client) *ZoneStorage {
	return &ZoneStorage{client: client}
}

// ListZones 列出所有 Zone
func (s *ZoneStorage) ListZones(ctx context.Context) ([]*models.Zone, error) {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return nil, errors.ErrEtcdUnavailable
	}

	resp, err := s.client.client.Get(ctx, storage.ZoneKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	zones := make([]*models.Zone, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var zone models.Zone
		if err := json.Unmarshal(kv.Value, &zone); err != nil {
			continue
		}
		zones = append(zones, &zone)
	}

	return zones, nil
}

// GetZone 获取 Zone 详情
func (s *ZoneStorage) GetZone(ctx context.Context, zone string) (*models.Zone, error) {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return nil, errors.ErrEtcdUnavailable
	}

	key := s.zoneKey(zone)
	resp, err := s.client.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, errors.ErrZoneNotFound
	}

	var z models.Zone
	if err := json.Unmarshal(resp.Kvs[0].Value, &z); err != nil {
		return nil, err
	}

	return &z, nil
}

// CreateZone 创建 Zone
func (s *ZoneStorage) CreateZone(ctx context.Context, zone *models.Zone) error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}

	key := s.zoneKey(zone.Zone)
	data, err := json.Marshal(zone)
	if err != nil {
		return err
	}

	_, err = s.client.client.Put(ctx, key, string(data))
	return err
}

// UpdateZone 更新 Zone
func (s *ZoneStorage) UpdateZone(ctx context.Context, zone *models.Zone) error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}

	key := s.zoneKey(zone.Zone)
	data, err := json.Marshal(zone)
	if err != nil {
		return err
	}

	_, err = s.client.client.Put(ctx, key, string(data))
	return err
}

// DeleteZone 删除 Zone（级联删除该 Zone 下所有 Domain）
func (s *ZoneStorage) DeleteZone(ctx context.Context, zone string) error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}

	// 删除 Zone 本身
	key := s.zoneKey(zone)
	_, err := s.client.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	// 级联删除该 Zone 下的所有 Domain
	domainPrefix := storage.DomainKeyPrefix + zone + "/"
	_, err = s.client.client.Delete(ctx, domainPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

// zoneKey 生成 Zone 的 etcd key
func (s *ZoneStorage) zoneKey(zone string) string {
	return storage.ZoneKeyPrefix + zone
}

// ZoneExists 检查 Zone 是否存在
func (s *ZoneStorage) ZoneExists(ctx context.Context, zone string) (bool, error) {
	_, err := s.GetZone(ctx, zone)
	if err != nil {
		if err == errors.ErrZoneNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UpdateZoneRecordCount 更新 Zone 的记录数量
func (s *ZoneStorage) UpdateZoneRecordCount(ctx context.Context, zone string, count int) error {
	z, err := s.GetZone(ctx, zone)
	if err != nil {
		return err
	}
	z.RecordCount = count
	return s.UpdateZone(ctx, z)
}

// IncrementZoneRecordCount 增加 Zone 的记录数量
func (s *ZoneStorage) IncrementZoneRecordCount(ctx context.Context, zone string, delta int) error {
	z, err := s.GetZone(ctx, zone)
	if err != nil {
		return err
	}
	z.RecordCount += delta
	return s.UpdateZone(ctx, z)
}

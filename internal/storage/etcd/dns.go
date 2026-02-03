package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"dancer/internal/errors"
	"dancer/internal/models"
	"dancer/internal/storage"
	"go.etcd.io/etcd/client/v3"
)

// DNSStorage DNS存储操作
type DNSStorage struct {
	client *Client
}

func NewDNSStorage(client *Client) *DNSStorage {
	return &DNSStorage{client: client}
}

// checkConnection 检查 etcd 连接，使用默认超时
func (s *DNSStorage) checkConnection() error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}
	return nil
}

// GenerateKey 生成DNS记录的etcd key
// domain: example.com -> /coredns/com/example/
// index: x1, x2, x3... 用于同一域名的多个记录
func (s *DNSStorage) GenerateKey(domain string, index int) string {
	parts := strings.Split(domain, ".")
	// 反转域名部分
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return storage.DNSKeyPrefix + strings.Join(parts, "/") + "/x" + strconv.Itoa(index)
}

// GetRecord 根据key获取DNS记录
func (s *DNSStorage) GetRecord(ctx context.Context, key string) (*models.DNSRecord, error) {
	if err := s.checkConnection(); err != nil {
		return nil, err
	}

	resp, err := s.client.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, errors.ErrRecordNotFound
	}

	var record models.DNSRecord
	if err := json.Unmarshal(resp.Kvs[0].Value, &record); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dns record: %w", err)
	}

	return &record, nil
}

// ListRecords 列出指定域名的所有记录
func (s *DNSStorage) ListRecords(ctx context.Context, domain string) ([]*models.DNSRecord, error) {
	if err := s.checkConnection(); err != nil {
		return nil, err
	}

	prefix := s.getDomainPrefix(domain)
	resp, err := s.client.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	records := make([]*models.DNSRecord, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var record models.DNSRecord
		if err := json.Unmarshal(kv.Value, &record); err != nil {
			continue
		}
		records = append(records, &record)
	}

	return records, nil
}

// ListAllRecords 列出所有DNS记录
func (s *DNSStorage) ListAllRecords(ctx context.Context) ([]*models.DNSRecord, error) {
	if err := s.checkConnection(); err != nil {
		return nil, err
	}

	resp, err := s.client.client.Get(ctx, storage.DNSKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	records := make([]*models.DNSRecord, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var record models.DNSRecord
		if err := json.Unmarshal(kv.Value, &record); err != nil {
			continue
		}
		records = append(records, &record)
	}

	return records, nil
}

// CreateRecord 创建DNS记录
func (s *DNSStorage) CreateRecord(ctx context.Context, record *models.DNSRecord) error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	// 查找该域名的下一个可用索引
	records, err := s.ListRecords(ctx, record.Domain)
	if err != nil {
		return err
	}

	index := len(records) + 1
	record.Key = s.GenerateKey(record.Domain, index)
	record.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	record.CreatedAt = time.Now().Unix()
	record.UpdatedAt = record.CreatedAt

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal dns record: %w", err)
	}

	_, err = s.client.client.Put(ctx, record.Key, string(data))
	return err
}

// UpdateRecord 更新DNS记录
func (s *DNSStorage) UpdateRecord(ctx context.Context, record *models.DNSRecord) error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	// 检查记录是否存在
	existing, err := s.GetRecord(ctx, record.Key)
	if err != nil {
		return err
	}

	// 保留创建时间和ID
	record.ID = existing.ID
	record.CreatedAt = existing.CreatedAt
	record.UpdatedAt = time.Now().Unix()

	// 如果域名变更，需要重新生成key
	if record.Domain != existing.Domain {
		// 删除旧记录
		if err := s.DeleteRecord(ctx, record.Key); err != nil {
			return err
		}
		// 创建新记录（会重新生成key）
		return s.CreateRecord(ctx, record)
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal dns record: %w", err)
	}

	_, err = s.client.client.Put(ctx, record.Key, string(data))
	return err
}

// DeleteRecord 删除DNS记录
func (s *DNSStorage) DeleteRecord(ctx context.Context, key string) error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	_, err := s.client.client.Delete(ctx, key)
	return err
}

// getDomainPrefix 获取域名的前缀
func (s *DNSStorage) getDomainPrefix(domain string) string {
	parts := strings.Split(domain, ".")
	// 反转域名部分
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return storage.DNSKeyPrefix + strings.Join(parts, "/") + "/"
}

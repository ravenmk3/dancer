package etcd

import (
	"context"
	"encoding/json"
	"path"
	"strconv"
	"strings"
	"time"

	"dancer/internal/config"
	"dancer/internal/errors"
	"dancer/internal/models"
	"dancer/internal/storage"
	"go.etcd.io/etcd/client/v3"
)

// DomainStorage Domain 存储操作
type DomainStorage struct {
	client *Client
	config *config.Config
}

func NewDomainStorage(client *Client, cfg *config.Config) *DomainStorage {
	return &DomainStorage{
		client: client,
		config: cfg,
	}
}

// getCoreDNSPrefix 获取 CoreDNS etcd 前缀
func (s *DomainStorage) getCoreDNSPrefix() string {
	prefix := s.config.Etcd.CorednsPrefix
	if prefix == "" {
		prefix = "/skydns/"
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	return prefix
}

// ListDomainsByZone 列出 Zone 下所有 Domain
func (s *DomainStorage) ListDomainsByZone(ctx context.Context, zone string) ([]*models.Domain, error) {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return nil, errors.ErrEtcdUnavailable
	}

	prefix := s.domainPrefix(zone)
	resp, err := s.client.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	domains := make([]*models.Domain, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var domain models.Domain
		if err := json.Unmarshal(kv.Value, &domain); err != nil {
			continue
		}
		domains = append(domains, &domain)
	}

	return domains, nil
}

// GetDomain 获取 Domain 详情
func (s *DomainStorage) GetDomain(ctx context.Context, zone, domain string) (*models.Domain, error) {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return nil, errors.ErrEtcdUnavailable
	}

	key := s.domainKey(zone, domain)
	resp, err := s.client.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, errors.ErrDomainNotFound
	}

	var d models.Domain
	if err := json.Unmarshal(resp.Kvs[0].Value, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

// CreateDomain 创建 Domain
func (s *DomainStorage) CreateDomain(ctx context.Context, domain *models.Domain) error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}

	// 检查是否已存在
	existing, err := s.GetDomain(ctx, domain.Zone, domain.Domain)
	if err == nil && existing != nil {
		return errors.ErrDomainExists
	}

	// 设置元数据
	now := time.Now().Unix()
	domain.Name = domain.Domain + "." + domain.Zone
	domain.RecordCount = len(domain.IPs)
	domain.CreatedAt = now
	domain.UpdatedAt = now

	// 保存 Domain 元数据
	key := s.domainKey(domain.Zone, domain.Domain)
	data, err := json.Marshal(domain)
	if err != nil {
		return err
	}

	_, err = s.client.client.Put(ctx, key, string(data))
	if err != nil {
		return err
	}

	// 同步到 CoreDNS
	return s.syncToCoreDNS(ctx, domain)
}

// UpdateDomain 更新 Domain
func (s *DomainStorage) UpdateDomain(ctx context.Context, domain *models.Domain) error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}

	// 获取现有记录
	existing, err := s.GetDomain(ctx, domain.Zone, domain.Domain)
	if err != nil {
		return err
	}

	// 更新时间戳
	domain.Name = domain.Domain + "." + domain.Zone
	domain.RecordCount = len(domain.IPs)
	domain.CreatedAt = existing.CreatedAt
	domain.UpdatedAt = time.Now().Unix()

	// 如果 TTL 未设置，使用原值
	if domain.TTL == 0 {
		domain.TTL = existing.TTL
	}

	// 保存 Domain 元数据
	key := s.domainKey(domain.Zone, domain.Domain)
	data, err := json.Marshal(domain)
	if err != nil {
		return err
	}

	_, err = s.client.client.Put(ctx, key, string(data))
	if err != nil {
		return err
	}

	// 同步到 CoreDNS
	return s.syncToCoreDNS(ctx, domain)
}

// DeleteDomain 删除 Domain
func (s *DomainStorage) DeleteDomain(ctx context.Context, zone, domain string) error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}

	// 删除 Domain 元数据
	key := s.domainKey(zone, domain)
	_, err := s.client.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	// 删除 CoreDNS 记录
	return s.deleteCoreDNSRecords(ctx, zone, domain)
}

// DeleteDomainsByZone 删除 Zone 下所有 Domain（级联删除）
func (s *DomainStorage) DeleteDomainsByZone(ctx context.Context, zone string) error {
	if err := s.client.WaitForConnection(defaultWaitTimeout); err != nil {
		return errors.ErrEtcdUnavailable
	}

	// 获取该 Zone 下所有 Domain
	domains, err := s.ListDomainsByZone(ctx, zone)
	if err != nil {
		return err
	}

	// 删除每个 Domain 的 CoreDNS 记录
	for _, domain := range domains {
		if err := s.deleteCoreDNSRecords(ctx, zone, domain.Domain); err != nil {
			return err
		}
	}

	// 删除所有 Domain 元数据
	prefix := s.domainPrefix(zone)
	_, err = s.client.client.Delete(ctx, prefix, clientv3.WithPrefix())
	return err
}

// domainKey 生成 Domain 的 etcd key
func (s *DomainStorage) domainKey(zone, domain string) string {
	return storage.DomainKeyPrefix + zone + "/" + domain
}

// domainPrefix 生成 Domain 前缀
func (s *DomainStorage) domainPrefix(zone string) string {
	return storage.DomainKeyPrefix + zone + "/"
}

// generateCoreDNSKey 生成 CoreDNS 的 etcd key
func (s *DomainStorage) generateCoreDNSKey(zone, domain, index string) string {
	// 反转 zone：example.com -> com/example
	reversed := reverseZone(zone)
	prefix := s.getCoreDNSPrefix()
	return path.Join(prefix, reversed, domain, "x"+index)
}

// reverseZone 反转域名层级
func reverseZone(zone string) string {
	parts := strings.Split(zone, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return path.Join(parts...)
}

// syncToCoreDNS 同步 Domain 到 CoreDNS
func (s *DomainStorage) syncToCoreDNS(ctx context.Context, domain *models.Domain) error {
	// 获取现有的 CoreDNS 记录
	existingKeys, err := s.getCoreDNSRecordKeys(ctx, domain.Zone, domain.Domain)
	if err != nil {
		return err
	}

	// 计算需要添加和删除的记录
	desiredIPs := make(map[string]bool)
	for _, ip := range domain.IPs {
		desiredIPs[ip] = true
	}

	// 找出需要删除的记录
	for key, ip := range existingKeys {
		if !desiredIPs[ip] {
			// 删除不再需要的记录
			_, err := s.client.client.Delete(ctx, key)
			if err != nil {
				return err
			}
			delete(existingKeys, key)
		}
	}

	// 找出需要添加的记录
	existingIPs := make(map[string]bool)
	for _, ip := range existingKeys {
		existingIPs[ip] = true
	}

	// 添加新记录或更新现有记录
	for i, ip := range domain.IPs {
		if !existingIPs[ip] {
			// 新记录，需要找到下一个可用索引
			index := strconv.Itoa(i + 1)
			for {
				key := s.generateCoreDNSKey(domain.Zone, domain.Domain, index)
				_, exists := existingKeys[key]
				if !exists {
					// 检查这个 key 是否被其他记录占用
					resp, err := s.client.client.Get(ctx, key)
					if err != nil {
						return err
					}
					if len(resp.Kvs) == 0 {
						break
					}
				}
				// 尝试下一个索引
				idx, _ := strconv.Atoi(index)
				index = strconv.Itoa(idx + 1)
			}

			key := s.generateCoreDNSKey(domain.Zone, domain.Domain, index)
			record := map[string]interface{}{
				"host": ip,
				"ttl":  domain.TTL,
			}
			data, err := json.Marshal(record)
			if err != nil {
				return err
			}
			_, err = s.client.client.Put(ctx, key, string(data))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// getCoreDNSRecordKeys 获取 CoreDNS 记录的 keys
func (s *DomainStorage) getCoreDNSRecordKeys(ctx context.Context, zone, domain string) (map[string]string, error) {
	prefix := s.getCoreDNSPrefix()
	reversed := reverseZone(zone)
	keyPrefix := path.Join(prefix, reversed, domain) + "/"

	resp, err := s.client.client.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, kv := range resp.Kvs {
		var record map[string]interface{}
		if err := json.Unmarshal(kv.Value, &record); err != nil {
			continue
		}
		if host, ok := record["host"].(string); ok {
			result[string(kv.Key)] = host
		}
	}

	return result, nil
}

// deleteCoreDNSRecords 删除 CoreDNS 记录
func (s *DomainStorage) deleteCoreDNSRecords(ctx context.Context, zone, domain string) error {
	prefix := s.getCoreDNSPrefix()
	reversed := reverseZone(zone)
	keyPrefix := path.Join(prefix, reversed, domain) + "/"

	_, err := s.client.client.Delete(ctx, keyPrefix, clientv3.WithPrefix())
	return err
}

// DomainExists 检查 Domain 是否存在
func (s *DomainStorage) DomainExists(ctx context.Context, zone, domain string) (bool, error) {
	_, err := s.GetDomain(ctx, zone, domain)
	if err != nil {
		if err == errors.ErrDomainNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetDomainCountByZone 获取 Zone 下的 Domain 数量
func (s *DomainStorage) GetDomainCountByZone(ctx context.Context, zone string) (int, error) {
	domains, err := s.ListDomainsByZone(ctx, zone)
	if err != nil {
		return 0, err
	}
	return len(domains), nil
}

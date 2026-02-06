package models

// Domain 完整域名模型
type Domain struct {
	Zone        string   `json:"zone"`         // 所属 zone，如 example.com
	Domain      string   `json:"domain"`       // 子域名部分，如 www
	Name        string   `json:"name"`         // 完整域名，如 www.example.com
	IPs         []string `json:"ips"`          // IP 地址列表
	TTL         int      `json:"ttl"`          // TTL (秒)
	RecordCount int      `json:"record_count"` // IP 记录数量
	CreatedAt   int64    `json:"created_at"`   // 创建时间戳
	UpdatedAt   int64    `json:"updated_at"`   // 更新时间戳
}

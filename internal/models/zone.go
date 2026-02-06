package models

// Zone 二级域名（Zone）模型
type Zone struct {
	Zone        string `json:"zone"`         // 二级域名，如 example.com
	RecordCount int    `json:"record_count"` // 该 zone 下的域名数量
	CreatedAt   int64  `json:"created_at"`   // 创建时间戳
	UpdatedAt   int64  `json:"updated_at"`   // 更新时间戳
}

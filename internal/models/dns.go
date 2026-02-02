package models

type DNSRecord struct {
	Key       string `json:"key"`
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	IP        string `json:"ip"`
	TTL       int64  `json:"ttl"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

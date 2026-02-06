package models

// 请求 DTO

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=32"`
	Password string   `json:"password" validate:"required,min=6"`
	UserType UserType `json:"user_type" validate:"required,oneof=admin normal"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID       string   `json:"id" validate:"required"`
	Username string   `json:"username" validate:"omitempty,min=3,max=32"`
	Password string   `json:"password" validate:"omitempty,min=6"`
	UserType UserType `json:"user_type" validate:"omitempty,oneof=admin normal"`
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// Zone 相关请求

// ListZonesRequest 列出所有 Zone 请求
type ListZonesRequest struct{}

// GetZoneRequest 获取 Zone 详情请求
type GetZoneRequest struct {
	Zone string `json:"zone" validate:"required,fqdn"`
}

// CreateZoneRequest 创建 Zone 请求
type CreateZoneRequest struct {
	Zone string `json:"zone" validate:"required,fqdn"`
}

// UpdateZoneRequest 更新 Zone 请求
type UpdateZoneRequest struct {
	Zone string `json:"zone" validate:"required,fqdn"`
}

// DeleteZoneRequest 删除 Zone 请求
type DeleteZoneRequest struct {
	Zone string `json:"zone" validate:"required,fqdn"`
}

// Domain 相关请求

// ListDomainsRequest 列出 Zone 下所有 Domain 请求
type ListDomainsRequest struct {
	Zone string `json:"zone" validate:"required,fqdn"`
}

// GetDomainRequest 获取 Domain 详情请求
type GetDomainRequest struct {
	Zone   string `json:"zone" validate:"required,fqdn"`
	Domain string `json:"domain" validate:"required"`
}

// CreateDomainRequest 创建 Domain 请求
type CreateDomainRequest struct {
	Zone   string   `json:"zone" validate:"required,fqdn"`
	Domain string   `json:"domain" validate:"required"`
	IPs    []string `json:"ips" validate:"required,dive,ip"`
	TTL    int      `json:"ttl" validate:"required,min=1"`
}

// UpdateDomainRequest 更新 Domain 请求
type UpdateDomainRequest struct {
	Zone   string   `json:"zone" validate:"required,fqdn"`
	Domain string   `json:"domain" validate:"required"`
	IPs    []string `json:"ips" validate:"required,dive,ip"`
	TTL    int      `json:"ttl" validate:"omitempty,min=1"`
}

// DeleteDomainRequest 删除 Domain 请求
type DeleteDomainRequest struct {
	Zone   string `json:"zone" validate:"required,fqdn"`
	Domain string `json:"domain" validate:"required"`
}

// 响应 DTO

// Response 统一响应结构
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []*User `json:"users"`
}

// ZoneListResponse Zone 列表响应
type ZoneListResponse struct {
	Zones []*Zone `json:"zones"`
}

// ZoneResponse Zone 详情响应
type ZoneResponse struct {
	Zone *Zone `json:"zone"`
}

// DomainListResponse Domain 列表响应
type DomainListResponse struct {
	Domains []*Domain `json:"domains"`
}

// DomainResponse Domain 详情响应
type DomainResponse struct {
	Domain *Domain `json:"domain"`
}

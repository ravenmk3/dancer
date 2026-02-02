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

// CreateDNSRequest 创建DNS记录请求
type CreateDNSRequest struct {
	Domain string `json:"domain" validate:"required,fqdn"`
	IP     string `json:"ip" validate:"required,ip"`
	TTL    int64  `json:"ttl" validate:"required,min=60"`
}

// UpdateDNSRequest 更新DNS记录请求
type UpdateDNSRequest struct {
	Key    string `json:"key" validate:"required"`
	Domain string `json:"domain" validate:"omitempty,fqdn"`
	IP     string `json:"ip" validate:"omitempty,ip"`
	TTL    int64  `json:"ttl" validate:"omitempty,min=60"`
}

// DeleteDNSRequest 删除DNS记录请求
type DeleteDNSRequest struct {
	Key string `json:"key" validate:"required"`
}

// ListDNSRequest 列出DNS记录请求
type ListDNSRequest struct {
	Domain string `json:"domain"` // 可选，为空时列出所有
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
	User  *User  `json:"user"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []*User `json:"users"`
}

// DNSListResponse DNS记录列表响应
type DNSListResponse struct {
	Records []*DNSRecord `json:"records"`
}

// DNSResponse DNS记录响应
type DNSResponse struct {
	Record *DNSRecord `json:"record"`
}

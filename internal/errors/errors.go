package errors

import "net/http"

// BusinessError 业务错误接口
type BusinessError interface {
	error
	Code() string // 业务错误码
	Status() int  // HTTP 状态码
}

// businessError 业务错误实现
type businessError struct {
	code    string
	status  int
	message string
}

func (e *businessError) Error() string { return e.message }
func (e *businessError) Code() string  { return e.code }
func (e *businessError) Status() int   { return e.status }

// Is 实现 error.Is 接口，按错误码比较
func (e *businessError) Is(target error) bool {
	t, ok := target.(*businessError)
	if !ok {
		return false
	}
	return e.code == t.code
}

// NewBusinessError 创建业务错误
func NewBusinessError(code string, status int, message string) BusinessError {
	return &businessError{
		code:    code,
		status:  status,
		message: message,
	}
}

// 预定义的错误变量
var (
	// 用户相关错误
	ErrUserNotFound       = NewBusinessError("user_not_found", http.StatusNotFound, "user not found")
	ErrUserExists         = NewBusinessError("user_exists", http.StatusConflict, "user already exists")
	ErrInvalidCredentials = NewBusinessError("invalid_credentials", http.StatusUnauthorized, "invalid username or password")
	ErrWrongPassword      = NewBusinessError("wrong_password", http.StatusBadRequest, "wrong password")
	ErrSamePassword       = NewBusinessError("same_password", http.StatusBadRequest, "new password cannot be the same as old password")
	ErrCannotDeleteDefaultAdmin = NewBusinessError("cannot_delete_default_admin", http.StatusForbidden, "cannot delete default admin user")

	// DNS 记录相关错误
	ErrRecordNotFound = NewBusinessError("record_not_found", http.StatusNotFound, "DNS record not found")
	ErrRecordExists   = NewBusinessError("record_exists", http.StatusConflict, "DNS record already exists")

	// Zone 相关错误
	ErrZoneNotFound = NewBusinessError("zone_not_found", http.StatusNotFound, "zone not found")
	ErrZoneExists   = NewBusinessError("zone_exists", http.StatusConflict, "zone already exists")

	// Domain 相关错误
	ErrDomainNotFound = NewBusinessError("domain_not_found", http.StatusNotFound, "domain not found")
	ErrDomainExists   = NewBusinessError("domain_exists", http.StatusConflict, "domain already exists")

	// 认证授权错误
	ErrInvalidToken = NewBusinessError("invalid_token", http.StatusUnauthorized, "invalid token")
	ErrTokenExpired = NewBusinessError("token_expired", http.StatusUnauthorized, "token expired")
	ErrUnauthorized = NewBusinessError("unauthorized", http.StatusUnauthorized, "unauthorized")
	ErrForbidden    = NewBusinessError("forbidden", http.StatusForbidden, "forbidden")

	// 输入错误
	ErrInvalidInput    = NewBusinessError("invalid_input", http.StatusBadRequest, "invalid input")
	ErrPasswordTooLong = NewBusinessError("invalid_input", http.StatusBadRequest, "password exceeds maximum length of 72 bytes")

	// 系统错误
	ErrEtcdUnavailable = NewBusinessError("service_unavailable", http.StatusServiceUnavailable, "etcd service temporarily unavailable, please retry later")
)

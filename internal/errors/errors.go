package errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrWrongPassword      = errors.New("wrong password")
	ErrRecordNotFound     = errors.New("DNS record not found")
	ErrRecordExists       = errors.New("DNS record already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidInput       = errors.New("invalid input")
	ErrEtcdUnavailable    = errors.New("etcd service temporarily unavailable")

	// Zone 相关错误
	ErrZoneNotFound = errors.New("zone not found")
	ErrZoneExists   = errors.New("zone already exists")

	// Domain 相关错误
	ErrDomainNotFound = errors.New("domain not found")
	ErrDomainExists   = errors.New("domain already exists")

	// 其他业务错误
	ErrCannotDeleteDefaultAdmin = errors.New("cannot delete default admin user")

	// 密码相关错误
	ErrPasswordTooLong = errors.New("password exceeds maximum length of 72 bytes")
)

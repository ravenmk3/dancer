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
)

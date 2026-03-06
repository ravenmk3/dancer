package handlers

import (
	"dancer/internal/errors"
	"github.com/go-playground/validator/v10"
)

// handleValidationError 统一处理验证错误
// 将 validator.ValidationErrors 转换为业务错误 ErrInvalidInput
func handleValidationError(err error) error {
	if _, ok := err.(validator.ValidationErrors); ok {
		return errors.ErrInvalidInput
	}
	return err
}

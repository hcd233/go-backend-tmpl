package model

import (
	"fmt"
	"net/http"

	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/enum"
	"github.com/hcd233/aris-api-tmpl/internal/i18n"
)

// Error 错误
//
//	@author centonhuang
//	@update 2025-11-10 19:10:53
type Error struct {
	Code    int    `json:"code" doc:"Code"`
	Message string `json:"message" doc:"Message"`
	// MessageKey i18n 翻译键；非空时 Localize 会按请求语言翻译 Message。
	MessageKey string `json:"-"`
}

// NewError 创建错误
//
//	@param code int
//	@param message string
//	@return *Error
//	@author centonhuang
//	@update 2025-11-10 19:14:00
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithKey 创建带翻译键的错误。
func NewErrorWithKey(code int, message, key string) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		MessageKey: key,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf(constant.ErrorModelTemplate, e.Code, e.Message)
}

// Localize 根据 locale 翻译错误消息；无翻译键时原样返回。
func (e *Error) Localize(locale enum.Locale) *Error {
	if e.MessageKey == "" {
		return e
	}
	return &Error{
		Code:       e.Code,
		Message:    i18n.Translate(locale, e.MessageKey, e.Message),
		MessageKey: e.MessageKey,
	}
}

// StatusCode 将业务错误码映射为 HTTP 状态码。
//
// 业务错误统一走顶层 {"error": {code, message}} 响应后，HTTP 状态码由业务码推导。
// 未显式映射的错误码兜底为 500。
func (e *Error) StatusCode() int {
	switch e.Code {
	case constant.BizErrorCodeUnauthorized: // ErrUnauthorized / ErrJWTDecode
		return http.StatusUnauthorized
	case constant.BizErrorCodeNoPermission: // ErrNoPermission
		return http.StatusForbidden
	case constant.BizErrorCodeDataNotExists: // ErrDataNotExists
		return http.StatusNotFound
	case constant.BizErrorCodeDataExists: // ErrDataExists
		return http.StatusConflict
	case constant.BizErrorCodeTooManyRequests, constant.BizErrorCodeInsufficientQuota: // ErrTooManyRequests / ErrInsufficientQuota
		return http.StatusTooManyRequests
	case constant.BizErrorCodeBadRequest: // ErrBadRequest / ErrValidation
		return http.StatusBadRequest
	case constant.BizErrorCodeResourceLocked: // ErrResourceLocked
		return http.StatusLocked
	default:
		return http.StatusInternalServerError
	}
}

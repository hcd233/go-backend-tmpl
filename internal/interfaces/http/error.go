// Package http 错误转换
//
//	@update 2025-09-30 00:00:00
package http

import (
	"errors"

	appauth "github.com/hcd233/go-backend-tmpl/internal/application/auth"
	appuser "github.com/hcd233/go-backend-tmpl/internal/application/user"
	"github.com/hcd233/go-backend-tmpl/internal/domain/user"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
)

// convertAppErrorToProtocol 将应用层错误转换为协议层错误
func convertAppErrorToProtocol(err error) error {
	if err == nil {
		return nil
	}

	// 用户应用层错误
	if errors.Is(err, appuser.ErrUserNotFound) {
		return protocol.ErrDataNotExists
	}
	if errors.Is(err, appuser.ErrInternalError) {
		return protocol.ErrInternalError
	}

	// 认证应用层错误
	if errors.Is(err, appauth.ErrUnauthorized) {
		return protocol.ErrUnauthorized
	}
	if errors.Is(err, appauth.ErrInvalidProvider) {
		return protocol.ErrBadRequest
	}
	if errors.Is(err, appauth.ErrInternalError) {
		return protocol.ErrInternalError
	}

	// 用户领域层错误
	if errors.Is(err, user.ErrInvalidUserName) {
		return protocol.ErrBadRequest
	}
	if errors.Is(err, user.ErrUserNotFound) {
		return protocol.ErrDataNotExists
	}
	if errors.Is(err, user.ErrUserAlreadyExists) {
		return protocol.ErrDataExists
	}

	// 默认返回内部错误
	return protocol.ErrInternalError
}

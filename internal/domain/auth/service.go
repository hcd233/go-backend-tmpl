// Package auth 认证领域服务
//
//	@update 2025-09-30 00:00:00
package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/oauth2"
)

var (
	// ErrInvalidAuthCode 无效的授权码
	ErrInvalidAuthCode = errors.New("invalid auth code")
	// ErrInvalidState 无效的state
	ErrInvalidState = errors.New("invalid state")
	// ErrFailedToGetUserInfo 获取用户信息失败
	ErrFailedToGetUserInfo = errors.New("failed to get user info")
)

// OAuth2Provider OAuth2提供商接口
//
//	@update 2025-09-30 00:00:00
type OAuth2Provider interface {
	// GetAuthURL 获取授权URL
	GetAuthURL() string

	// ExchangeToken 通过授权码获取Access Token
	ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error)

	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuth2UserInfo, error)

	// GetProvider 获取提供商类型
	GetProvider() Provider
}

// TokenSigner JWT令牌签名器接口
//
//	@update 2025-09-30 00:00:00
type TokenSigner interface {
	// EncodeToken 编码令牌
	EncodeToken(userID uint) (string, error)

	// DecodeToken 解码令牌
	DecodeToken(token string) (uint, error)
}

// Service 认证领域服务
//
//	@update 2025-09-30 00:00:00
type Service struct {
	accessTokenSigner  TokenSigner
	refreshTokenSigner TokenSigner
	stateString        string
}

// NewService 创建认证领域服务
func NewService(accessTokenSigner, refreshTokenSigner TokenSigner, stateString string) *Service {
	return &Service{
		accessTokenSigner:  accessTokenSigner,
		refreshTokenSigner: refreshTokenSigner,
		stateString:        stateString,
	}
}

// ValidateState 验证OAuth2 state参数
func (s *Service) ValidateState(state string) error {
	if state != s.stateString {
		return ErrInvalidState
	}
	return nil
}

// GenerateTokens 为用户生成访问令牌和刷新令牌
func (s *Service) GenerateTokens(userID uint) (*Token, error) {
	accessToken, err := s.accessTokenSigner.EncodeToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshTokenSigner.EncodeToken(userID)
	if err != nil {
		return nil, err
	}

	// 这里简化处理，实际应该从配置中获取过期时间
	// 假设访问令牌1小时后过期
	expiresAt := time.Now().Add(1 * time.Hour)
	token := NewToken(accessToken, refreshToken, expiresAt)

	return token, nil
}

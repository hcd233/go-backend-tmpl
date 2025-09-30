// Package auth 认证领域模型
//
//	@update 2025-09-30 00:00:00
package auth

import (
	"errors"
	"time"
)

// Provider OAuth2提供商类型
type Provider string

const (
	// ProviderGithub GitHub OAuth2提供商
	ProviderGithub Provider = "github"
	// ProviderQQ QQ OAuth2提供商
	ProviderQQ Provider = "qq"
	// ProviderGoogle Google OAuth2提供商
	ProviderGoogle Provider = "google"
)

// IsValid 验证提供商是否有效
func (p Provider) IsValid() bool {
	switch p {
	case ProviderGithub, ProviderQQ, ProviderGoogle:
		return true
	default:
		return false
	}
}

// OAuth2UserInfo OAuth2用户信息值对象
//
//	@update 2025-09-30 00:00:00
type OAuth2UserInfo struct {
	id       string
	name     string
	email    string
	avatar   string
	provider Provider
}

// NewOAuth2UserInfo 创建OAuth2用户信息
func NewOAuth2UserInfo(id, name, email, avatar string, provider Provider) (*OAuth2UserInfo, error) {
	if id == "" {
		return nil, errors.New("OAuth2 user ID cannot be empty")
	}
	if !provider.IsValid() {
		return nil, errors.New("invalid OAuth2 provider")
	}
	return &OAuth2UserInfo{
		id:       id,
		name:     name,
		email:    email,
		avatar:   avatar,
		provider: provider,
	}, nil
}

// ID 获取用户ID
func (o *OAuth2UserInfo) ID() string {
	return o.id
}

// Name 获取用户名
func (o *OAuth2UserInfo) Name() string {
	return o.name
}

// Email 获取邮箱
func (o *OAuth2UserInfo) Email() string {
	return o.email
}

// Avatar 获取头像
func (o *OAuth2UserInfo) Avatar() string {
	return o.avatar
}

// Provider 获取提供商
func (o *OAuth2UserInfo) Provider() Provider {
	return o.provider
}

// Token 认证令牌值对象
//
//	@update 2025-09-30 00:00:00
type Token struct {
	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

// NewToken 创建令牌
func NewToken(accessToken, refreshToken string, expiresAt time.Time) *Token {
	return &Token{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
	}
}

// AccessToken 获取访问令牌
func (t *Token) AccessToken() string {
	return t.accessToken
}

// RefreshToken 获取刷新令牌
func (t *Token) RefreshToken() string {
	return t.refreshToken
}

// ExpiresAt 获取过期时间
func (t *Token) ExpiresAt() time.Time {
	return t.expiresAt
}

// IsExpired 判断令牌是否过期
func (t *Token) IsExpired() bool {
	return time.Now().After(t.expiresAt)
}


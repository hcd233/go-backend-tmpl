// Package oauth2 OAuth2提供商实现
//
//	@update 2025-09-30 00:00:00
package oauth2

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/hcd233/go-backend-tmpl/internal/config"
	"github.com/hcd233/go-backend-tmpl/internal/domain/auth"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

var googleUserScopes = []string{
	"openid",
	"profile",
	"email",
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
}

// GoogleProvider Google OAuth2提供商实现
type GoogleProvider struct {
	oauth2Config *oauth2.Config
}

// GoogleUserInfo Google用户信息结构体
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

// NewGoogleProvider 创建Google OAuth2提供商
func NewGoogleProvider() auth.OAuth2Provider {
	return &GoogleProvider{
		oauth2Config: &oauth2.Config{
			Endpoint:     google.Endpoint,
			Scopes:       googleUserScopes,
			ClientID:     config.Oauth2GoogleClientID,
			ClientSecret: config.Oauth2GoogleClientSecret,
			RedirectURL:  config.Oauth2GoogleRedirectURL,
		},
	}
}

// GetAuthURL 获取授权URL
func (p *GoogleProvider) GetAuthURL() string {
	return p.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
}

// ExchangeToken 通过授权码获取Access Token
func (p *GoogleProvider) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	logger := logger.WithCtx(ctx)

	logger.Info("[GoogleProvider] exchanging code for token",
		zap.String("clientID", p.oauth2Config.ClientID),
		zap.String("redirectURL", p.oauth2Config.RedirectURL))

	token, err := p.oauth2Config.Exchange(ctx, code)
	if err != nil {
		logger.Error("[GoogleProvider] token exchange failed", zap.Error(err))
		return nil, err
	}

	logger.Info("[GoogleProvider] token exchange successful")
	return token, nil
}

// GetUserInfo 获取用户信息
func (p *GoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*auth.OAuth2UserInfo, error) {
	logger := logger.WithCtx(ctx)
	client := p.oauth2Config.Client(ctx, token)

	logger.Info("[GoogleProvider] calling Google UserInfo API")

	resp, err := client.Get(googleUserInfoURL)
	if err != nil {
		logger.Error("[GoogleProvider] failed to call userinfo API", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	logger.Info("[GoogleProvider] userinfo API response", zap.Int("statusCode", resp.StatusCode))

	var userInfo GoogleUserInfo
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		logger.Error("[GoogleProvider] failed to decode userinfo response", zap.Error(err))
		return nil, err
	}

	logger.Info("[GoogleProvider] successfully decoded user info",
		zap.String("userID", userInfo.ID),
		zap.String("userName", userInfo.Name),
		zap.String("userEmail", userInfo.Email))

	// 转换为领域模型
	return auth.NewOAuth2UserInfo(
		userInfo.ID,
		userInfo.Name,
		userInfo.Email,
		userInfo.Picture,
		auth.ProviderGoogle,
	)
}

// GetProvider 获取提供商类型
func (p *GoogleProvider) GetProvider() auth.Provider {
	return auth.ProviderGoogle
}


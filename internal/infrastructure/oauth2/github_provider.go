// Package oauth2 OAuth2提供商实现
//
//	@update 2025-09-30 00:00:00
package oauth2

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/hcd233/go-backend-tmpl/internal/config"
	"github.com/hcd233/go-backend-tmpl/internal/domain/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
)

var githubUserScopes = []string{"user:email", "repo", "read:org"}

// GithubProvider GitHub OAuth2提供商实现
type GithubProvider struct {
	oauth2Config *oauth2.Config
}

// GithubUserInfo Github用户信息结构体
type GithubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GithubEmail Github邮箱信息结构体
type GithubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

// NewGithubProvider 创建GitHub OAuth2提供商
func NewGithubProvider() auth.OAuth2Provider {
	return &GithubProvider{
		oauth2Config: &oauth2.Config{
			Endpoint:     github.Endpoint,
			Scopes:       githubUserScopes,
			ClientID:     config.Oauth2GithubClientID,
			ClientSecret: config.Oauth2GithubClientSecret,
			RedirectURL:  config.Oauth2GithubRedirectURL,
		},
	}
}

// GetAuthURL 获取授权URL
func (p *GithubProvider) GetAuthURL() string {
	return p.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
}

// ExchangeToken 通过授权码获取Access Token
func (p *GithubProvider) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.oauth2Config.Exchange(ctx, code)
}

// GetUserInfo 获取用户信息
func (p *GithubProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*auth.OAuth2UserInfo, error) {
	client := p.oauth2Config.Client(ctx, token)

	// 获取用户基本信息
	resp, err := client.Get(githubUserURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GithubUserInfo
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	// 获取用户邮箱信息
	emailResp, err := client.Get(githubUserEmailURL)
	if err != nil {
		return nil, err
	}
	defer emailResp.Body.Close()

	var emails []GithubEmail
	if err := sonic.ConfigDefault.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	// 选择主邮箱
	for _, email := range emails {
		if email.Primary {
			userInfo.Email = email.Email
			break
		}
	}

	// 转换为领域模型
	return auth.NewOAuth2UserInfo(
		fmt.Sprintf("%d", userInfo.ID),
		userInfo.Login,
		userInfo.Email,
		userInfo.AvatarURL,
		auth.ProviderGithub,
	)
}

// GetProvider 获取提供商类型
func (p *GithubProvider) GetProvider() auth.Provider {
	return auth.ProviderGithub
}

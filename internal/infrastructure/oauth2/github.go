package oauth2

import (
	"context"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/domain/oauth2/vo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
)

var githubUserScopes = []string{"user:email", "repo", "read:org"}

// GithubUserInfo Github 用户信息结构体。
type GithubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GetID 获取 Github 用户 ID。
func (u *GithubUserInfo) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}

// GetName 获取 Github 用户名。
func (u *GithubUserInfo) GetName() string {
	return u.Login
}

// GetEmail 获取 Github 用户邮箱。
func (u *GithubUserInfo) GetEmail() string {
	return u.Email
}

// GetAvatar 获取 Github 用户头像。
func (u *GithubUserInfo) GetAvatar() string {
	return u.AvatarURL
}

// GithubEmail Github 邮箱信息结构体。
type GithubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

type githubPlatform struct {
	oauth2Config *oauth2.Config
}

// NewGithubPlatform 创建 Github OAuth2 平台。
func NewGithubPlatform() Platform {
	return &githubPlatform{
		oauth2Config: &oauth2.Config{
			Endpoint:     github.Endpoint,
			Scopes:       githubUserScopes,
			ClientID:     config.Oauth2GithubClientID,
			ClientSecret: config.Oauth2GithubClientSecret,
			RedirectURL:  config.Oauth2GithubRedirectURL,
		},
	}
}

func (p *githubPlatform) GetAuthURL() string {
	return p.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
}

func (p *githubPlatform) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.oauth2Config.Exchange(ctx, code)
}

func (p *githubPlatform) GetUserInfo(ctx context.Context, token *oauth2.Token) (vo.OAuthUserInfo, error) {
	client := p.oauth2Config.Client(ctx, token)
	resp, err := client.Get(githubUserURL)
	if err != nil {
		return vo.OAuthUserInfo{}, err
	}
	defer func() { _ = resp.Body.Close() }() //nolint:errcheck // best-effort close

	var userInfo GithubUserInfo
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return vo.OAuthUserInfo{}, err
	}
	emailResp, err := client.Get(githubUserEmailURL)
	if err != nil {
		return vo.OAuthUserInfo{}, err
	}
	defer func() { _ = emailResp.Body.Close() }() //nolint:errcheck // best-effort close

	var emails []GithubEmail
	if err := sonic.ConfigDefault.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		return vo.OAuthUserInfo{}, err
	}
	for _, email := range emails {
		if email.Primary {
			userInfo.Email = email.Email
			break
		}
	}
	return vo.NewOAuthUserInfo(userInfo.GetID(), userInfo.GetName(), userInfo.GetEmail(), userInfo.GetAvatar()), nil
}

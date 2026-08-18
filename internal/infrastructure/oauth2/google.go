package oauth2

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/domain/oauth2/vo"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
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

// GoogleUserInfo Google 用户信息结构体。
type GoogleUserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhotoURL string `json:"picture"`
}

type googleUserInfoAPIResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

// GetID 获取 Google 用户 ID。
func (u *GoogleUserInfo) GetID() string {
	return u.ID
}

// GetName 获取 Google 用户名。
func (u *GoogleUserInfo) GetName() string {
	return u.Name
}

// GetEmail 获取 Google 用户邮箱。
func (u *GoogleUserInfo) GetEmail() string {
	return u.Email
}

// GetAvatar 获取 Google 用户头像。
func (u *GoogleUserInfo) GetAvatar() string {
	return u.PhotoURL
}

type googlePlatform struct {
	oauth2Config *oauth2.Config
}

// NewGooglePlatform 创建 Google OAuth2 平台。
func NewGooglePlatform() Platform {
	return &googlePlatform{
		oauth2Config: &oauth2.Config{
			Endpoint:     google.Endpoint,
			Scopes:       googleUserScopes,
			ClientID:     config.Oauth2GoogleClientID,
			ClientSecret: config.Oauth2GoogleClientSecret,
			RedirectURL:  config.Oauth2GoogleRedirectURL,
		},
	}
}

func (p *googlePlatform) GetAuthURL() string {
	return p.oauth2Config.AuthCodeURL(config.Oauth2StateString, oauth2.AccessTypeOffline)
}

func (p *googlePlatform) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	log := logger.WithCtx(ctx)
	log.Info("[GoogleOauth2] exchanging code for token",
		zap.String("clientID", p.oauth2Config.ClientID),
		zap.String("redirectURL", p.oauth2Config.RedirectURL),
		zap.Strings("scopes", p.oauth2Config.Scopes))
	token, err := p.oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Error("[GoogleOauth2] token exchange failed", zap.Error(err))
		return nil, err
	}
	log.Info("[GoogleOauth2] token exchange successful")
	return token, nil
}

func (p *googlePlatform) GetUserInfo(ctx context.Context, token *oauth2.Token) (vo.OAuthUserInfo, error) {
	log := logger.WithCtx(ctx)
	client := p.oauth2Config.Client(ctx, token)
	log.Info("[GoogleOauth2] calling Google UserInfo API")
	resp, err := client.Get(googleUserInfoURL)
	if err != nil {
		log.Error("[GoogleOauth2] failed to call userinfo API", zap.Error(err))
		return vo.OAuthUserInfo{}, err
	}
	defer func() { _ = resp.Body.Close() }() //nolint:errcheck // best-effort close
	log.Info("[GoogleOauth2] userinfo API response", zap.Int("statusCode", resp.StatusCode))

	var userInfoResp googleUserInfoAPIResponse
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&userInfoResp); err != nil {
		log.Error("[GoogleOauth2] failed to decode userinfo response", zap.Error(err))
		return vo.OAuthUserInfo{}, err
	}
	log.Info("[GoogleOauth2] successfully decoded user info",
		zap.String("userID", userInfoResp.ID),
		zap.String("userName", userInfoResp.Name),
		zap.String("userEmail", userInfoResp.Email))
	return vo.NewOAuthUserInfo(userInfoResp.ID, userInfoResp.Name, userInfoResp.Email, userInfoResp.Picture), nil
}

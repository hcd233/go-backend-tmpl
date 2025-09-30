// Package auth 认证应用层DTO
//
//	@update 2025-09-30 00:00:00
package auth

// LoginCommand 登录命令
type LoginCommand struct {
	Provider string
}

// CallbackCommand OAuth2回调命令
type CallbackCommand struct {
	Code     string
	State    string
	Provider string
}

// RefreshTokenCommand 刷新令牌命令
type RefreshTokenCommand struct {
	RefreshToken string
}

// LoginResponse 登录响应
type LoginResponse struct {
	RedirectURL string `json:"redirectURL"`
}

// CallbackResponse OAuth2回调响应
type CallbackResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}


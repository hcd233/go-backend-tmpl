package modules

import (
	identityport "github.com/hcd233/aris-api-tmpl/internal/application/identity/port"
	oauth2port "github.com/hcd233/aris-api-tmpl/internal/application/oauth2/port"
	"github.com/hcd233/aris-api-tmpl/internal/handler"
	"go.uber.org/fx"
)

// HandlerModule 处理器模块：HTTP handler 装配。
var HandlerModule = fx.Module("handler",
	fx.Provide(
		newTokenDependencies,
		newOauth2Dependencies,
		newUserDependencies,
		handler.NewPingHandler,
		handler.NewTokenHandler,
		handler.NewOauth2Handler,
		handler.NewUserHandler,
	),
)

// newTokenDependencies 装配令牌 handler 依赖。
func newTokenDependencies(refresh identityport.RefreshTokensHandler) handler.TokenDependencies {
	return handler.TokenDependencies{Refresh: refresh}
}

// newOauth2Dependencies 装配 OAuth2 handler 依赖。
func newOauth2Dependencies(initiate oauth2port.InitiateLoginHandler, callback oauth2port.HandleCallbackHandler) handler.Oauth2Dependencies {
	return handler.Oauth2Dependencies{Initiate: initiate, Callback: callback}
}

// newUserDependencies 装配用户 handler 依赖。
func newUserDependencies(getCurrentUser identityport.GetCurrentUserHandler, updateProfile identityport.UpdateProfileHandler) handler.UserDependencies {
	return handler.UserDependencies{GetCurrentUser: getCurrentUser, UpdateProfile: updateProfile}
}

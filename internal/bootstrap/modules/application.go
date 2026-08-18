package modules

import (
	identitycommand "github.com/hcd233/aris-api-tmpl/internal/application/identity/command"
	identityport "github.com/hcd233/aris-api-tmpl/internal/application/identity/port"
	identityquery "github.com/hcd233/aris-api-tmpl/internal/application/identity/query"
	oauth2command "github.com/hcd233/aris-api-tmpl/internal/application/oauth2/command"
	oauth2port "github.com/hcd233/aris-api-tmpl/internal/application/oauth2/port"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity"
	identityservice "github.com/hcd233/aris-api-tmpl/internal/domain/identity/service"
	oauth2service "github.com/hcd233/aris-api-tmpl/internal/domain/oauth2/service"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/oauth2"
	"github.com/hcd233/aris-api-tmpl/internal/jwt"
	"go.uber.org/fx"
)

// ApplicationModule 应用层模块：命令/查询处理器。
var ApplicationModule = fx.Module("application",
	fx.Provide(
		fx.Annotate(newAccessTokenSigner, fx.ResultTags(`name:"accessSigner"`)),
		fx.Annotate(newRefreshTokenSigner, fx.ResultTags(`name:"refreshSigner"`)),
		newOauth2Platforms,
		newRefreshTokensHandler,
		identitycommand.NewUpdateProfileHandler,
		identityquery.NewGetCurrentUserHandler,
		oauth2command.NewInitiateLoginHandler,
		newHandleCallbackHandler,
	),
)

// newAccessTokenSigner 构造 access token 签名器。
func newAccessTokenSigner() identityservice.TokenSigner {
	return jwt.GetAccessTokenSigner()
}

// newRefreshTokenSigner 构造 refresh token 签名器。
func newRefreshTokenSigner() identityservice.TokenSigner {
	return jwt.GetRefreshTokenSigner()
}

// newOauth2Platforms 构造 OAuth2 平台注册表。
func newOauth2Platforms() map[string]oauth2service.Platform {
	return map[string]oauth2service.Platform{
		constant.OAuthProviderGithub: oauth2.NewGithubPlatform(),
		constant.OAuthProviderGoogle: oauth2.NewGooglePlatform(),
	}
}

type refreshTokensParams struct {
	fx.In

	UserRepo      identity.UserRepository
	AccessSigner  identityservice.TokenSigner `name:"accessSigner"`
	RefreshSigner identityservice.TokenSigner `name:"refreshSigner"`
}

// newRefreshTokensHandler 构造刷新令牌处理器。
func newRefreshTokensHandler(params refreshTokensParams) identityport.RefreshTokensHandler {
	return identitycommand.NewRefreshTokensHandler(params.UserRepo, params.AccessSigner, params.RefreshSigner)
}

type handleCallbackParams struct {
	fx.In

	Platforms     map[string]oauth2service.Platform
	UserRepo      identity.UserRepository
	AccessSigner  identityservice.TokenSigner `name:"accessSigner"`
	RefreshSigner identityservice.TokenSigner `name:"refreshSigner"`
	DirCreator    oauth2port.ObjectStorageDirCreator
}

// newHandleCallbackHandler 构造回调处理器。
func newHandleCallbackHandler(params handleCallbackParams) oauth2port.HandleCallbackHandler {
	return oauth2command.NewHandleCallbackHandler(
		params.Platforms,
		params.UserRepo,
		params.AccessSigner,
		params.RefreshSigner,
		params.DirCreator,
	)
}

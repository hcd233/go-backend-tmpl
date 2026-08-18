// Package bootstrap wires application startup dependencies.
package bootstrap

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-api-tmpl/internal/api"
	identitycommand "github.com/hcd233/aris-api-tmpl/internal/application/identity/command"
	identityport "github.com/hcd233/aris-api-tmpl/internal/application/identity/port"
	identityquery "github.com/hcd233/aris-api-tmpl/internal/application/identity/query"
	oauth2command "github.com/hcd233/aris-api-tmpl/internal/application/oauth2/command"
	oauth2port "github.com/hcd233/aris-api-tmpl/internal/application/oauth2/port"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity"
	identityservice "github.com/hcd233/aris-api-tmpl/internal/domain/identity/service"
	oauth2service "github.com/hcd233/aris-api-tmpl/internal/domain/oauth2/service"
	"github.com/hcd233/aris-api-tmpl/internal/handler"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/oauth2"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/repository"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/storage"
	"github.com/hcd233/aris-api-tmpl/internal/jwt"
	"go.uber.org/dig"
)

const (
	digNameAccessSigner  = "accessSigner"
	digNameRefreshSigner = "refreshSigner"
)

// Server 是启动阶段解析出的 HTTP 服务对象。
type Server struct {
	container *dig.Container
	App       *fiber.App
	HumaAPI   huma.API
}

// BuildServer 构建启动依赖容器并解析 HTTP 服务对象。
func BuildServer() (*Server, error) {
	container := dig.New()
	if err := provide(container); err != nil {
		return nil, err
	}
	var server *Server
	if err := container.Invoke(func(app *fiber.App, humaAPI huma.API) {
		server = &Server{container: container, App: app, HumaAPI: humaAPI}
	}); err != nil {
		return nil, err
	}
	return server, nil
}

func provide(container *dig.Container) error {
	if err := provideHTTP(container); err != nil {
		return err
	}
	if err := provideInfrastructure(container); err != nil {
		return err
	}
	if err := provideApplication(container); err != nil {
		return err
	}
	if err := provideHandlers(container); err != nil {
		return err
	}
	if err := container.Provide(newAccessTokenSigner, dig.Name(digNameAccessSigner)); err != nil {
		return err
	}
	if err := container.Provide(newRefreshTokenSigner, dig.Name(digNameRefreshSigner)); err != nil {
		return err
	}
	return nil
}

func provideHTTP(container *dig.Container) error {
	if err := container.Provide(api.NewFiberApp); err != nil {
		return err
	}
	if err := container.Provide(api.NewHumaAPI); err != nil {
		return err
	}
	return nil
}

func provideInfrastructure(container *dig.Container) error {
	if err := container.Provide(newUserRepository); err != nil {
		return err
	}
	if err := container.Provide(newOauth2Platforms); err != nil {
		return err
	}
	if err := container.Provide(newAudioDirCreator); err != nil {
		return err
	}
	return nil
}

func provideApplication(container *dig.Container) error {
	if err := container.Provide(newRefreshTokensHandler); err != nil {
		return err
	}
	if err := container.Provide(identitycommand.NewUpdateProfileHandler); err != nil {
		return err
	}
	if err := container.Provide(identityquery.NewGetCurrentUserHandler); err != nil {
		return err
	}
	if err := container.Provide(oauth2command.NewInitiateLoginHandler); err != nil {
		return err
	}
	if err := container.Provide(newHandleCallbackHandler); err != nil {
		return err
	}
	return nil
}

func provideHandlers(container *dig.Container) error {
	if err := container.Provide(newTokenDependencies); err != nil {
		return err
	}
	if err := container.Provide(newOauth2Dependencies); err != nil {
		return err
	}
	if err := container.Provide(newUserDependencies); err != nil {
		return err
	}
	if err := container.Provide(handler.NewPingHandler); err != nil {
		return err
	}
	if err := container.Provide(handler.NewTokenHandler); err != nil {
		return err
	}
	if err := container.Provide(handler.NewOauth2Handler); err != nil {
		return err
	}
	if err := container.Provide(handler.NewUserHandler); err != nil {
		return err
	}
	return nil
}

func newUserRepository() identity.UserRepository {
	return repository.NewUserRepository()
}

func newAudioDirCreator() oauth2port.ObjectStorageDirCreator {
	if config.CosAppID == "" && config.MinioEndpoint == "" {
		return nil
	}
	storage.InitObjectStorage()
	return repository.NewAudioDirCreator()
}

func newAccessTokenSigner() identityservice.TokenSigner {
	return jwt.GetAccessTokenSigner()
}

func newRefreshTokenSigner() identityservice.TokenSigner {
	return jwt.GetRefreshTokenSigner()
}

func newOauth2Platforms() map[string]oauth2service.Platform {
	return map[string]oauth2service.Platform{
		constant.OAuthProviderGithub: oauth2.NewGithubPlatform(),
		constant.OAuthProviderGoogle: oauth2.NewGooglePlatform(),
	}
}

type refreshTokensParams struct {
	dig.In

	UserRepo      identity.UserRepository
	AccessSigner  identityservice.TokenSigner `name:"accessSigner"`
	RefreshSigner identityservice.TokenSigner `name:"refreshSigner"`
}

func newRefreshTokensHandler(params refreshTokensParams) identityport.RefreshTokensHandler {
	return identitycommand.NewRefreshTokensHandler(params.UserRepo, params.AccessSigner, params.RefreshSigner)
}

type handleCallbackParams struct {
	dig.In

	Platforms     map[string]oauth2service.Platform
	UserRepo      identity.UserRepository
	AccessSigner  identityservice.TokenSigner `name:"accessSigner"`
	RefreshSigner identityservice.TokenSigner `name:"refreshSigner"`
	DirCreator    oauth2port.ObjectStorageDirCreator
}

func newHandleCallbackHandler(params handleCallbackParams) oauth2port.HandleCallbackHandler {
	return oauth2command.NewHandleCallbackHandler(
		params.Platforms,
		params.UserRepo,
		params.AccessSigner,
		params.RefreshSigner,
		params.DirCreator,
	)
}

func newTokenDependencies(refresh identityport.RefreshTokensHandler) handler.TokenDependencies {
	return handler.TokenDependencies{Refresh: refresh}
}

func newOauth2Dependencies(initiate oauth2port.InitiateLoginHandler, callback oauth2port.HandleCallbackHandler) handler.Oauth2Dependencies {
	return handler.Oauth2Dependencies{Initiate: initiate, Callback: callback}
}

func newUserDependencies(getCurrentUser identityport.GetCurrentUserHandler, updateProfile identityport.UpdateProfileHandler) handler.UserDependencies {
	return handler.UserDependencies{GetCurrentUser: getCurrentUser, UpdateProfile: updateProfile}
}

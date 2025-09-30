// Package container 使用dig的依赖注入容器
//
//	@update 2025-09-30 00:00:00
package container

import (
	"context"

	appauth "github.com/hcd233/go-backend-tmpl/internal/application/auth"
	appuser "github.com/hcd233/go-backend-tmpl/internal/application/user"
	"github.com/hcd233/go-backend-tmpl/internal/auth"
	"github.com/hcd233/go-backend-tmpl/internal/config"
	domainauth "github.com/hcd233/go-backend-tmpl/internal/domain/auth"
	"github.com/hcd233/go-backend-tmpl/internal/domain/user"
	"github.com/hcd233/go-backend-tmpl/internal/handler"
	"github.com/hcd233/go-backend-tmpl/internal/infrastructure/oauth2"
	"github.com/hcd233/go-backend-tmpl/internal/infrastructure/persistence"
	"github.com/hcd233/go-backend-tmpl/internal/interfaces/http"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	objdao "github.com/hcd233/go-backend-tmpl/internal/resource/storage/obj_dao"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

// Container dig依赖注入容器
type Container struct {
	container *dig.Container
}

// NewContainer 创建dig容器
func NewContainer() (*Container, error) {
	c := dig.New()
	dc := &Container{container: c}

	// 注册所有依赖
	if err := dc.registerDependencies(); err != nil {
		return nil, err
	}

	return dc, nil
}

// registerDependencies 注册所有依赖
func (dc *Container) registerDependencies() error {
	// 1. 注册仓储层
	if err := dc.container.Provide(persistence.NewUserRepository); err != nil {
		return err
	}

	// 2. 注册领域服务
	if err := dc.container.Provide(user.NewService); err != nil {
		return err
	}

	if err := dc.container.Provide(func() *domainauth.Service {
		return domainauth.NewService(
			auth.GetJwtAccessTokenSigner(),
			auth.GetJwtRefreshTokenSigner(),
			config.Oauth2StateString,
		)
	}); err != nil {
		return err
	}

	// 3. 注册OAuth2提供商
	if err := dc.container.Provide(func() map[domainauth.Provider]domainauth.OAuth2Provider {
		return map[domainauth.Provider]domainauth.OAuth2Provider{
			domainauth.ProviderGithub: oauth2.NewGithubProvider(),
			domainauth.ProviderGoogle: oauth2.NewGoogleProvider(),
		}
	}); err != nil {
		return err
	}

	// 4. 注册用户创建回调
	if err := dc.container.Provide(func() func(ctx context.Context, userID uint) error {
		return func(ctx context.Context, userID uint) error {
			logger := logger.WithCtx(ctx)
			imageObjDAO := objdao.GetImageObjDAO()
			thumbnailObjDAO := objdao.GetThumbnailObjDAO()

			if _, err := imageObjDAO.CreateDir(ctx, userID); err != nil {
				logger.Error("[Container] failed to create image dir", zap.Error(err))
				return err
			}

			if _, err := thumbnailObjDAO.CreateDir(ctx, userID); err != nil {
				logger.Error("[Container] failed to create thumbnail dir", zap.Error(err))
				return err
			}

			return nil
		}
	}); err != nil {
		return err
	}

	// 5. 注册应用服务
	if err := dc.container.Provide(appuser.NewService); err != nil {
		return err
	}

	if err := dc.container.Provide(func(
		authService *domainauth.Service,
		userRepo user.Repository,
		userService *user.Service,
		oauth2Providers map[domainauth.Provider]domainauth.OAuth2Provider,
		onUserCreated func(ctx context.Context, userID uint) error,
	) *appauth.Service {
		return appauth.NewService(
			authService,
			userRepo,
			userService,
			oauth2Providers,
			onUserCreated,
		)
	}); err != nil {
		return err
	}

	// 6. 注册Token签名器
	if err := dc.container.Provide(func() auth.JwtTokenSigner {
		return auth.GetJwtAccessTokenSigner()
	}, dig.Name("accessTokenSigner")); err != nil {
		return err
	}

	if err := dc.container.Provide(func() auth.JwtTokenSigner {
		return auth.GetJwtRefreshTokenSigner()
	}, dig.Name("refreshTokenSigner")); err != nil {
		return err
	}

	// 7. 注册HTTP处理器
	if err := dc.container.Provide(http.NewUserHandler); err != nil {
		return err
	}

	if err := dc.container.Provide(http.NewOAuth2Handler); err != nil {
		return err
	}

	if err := dc.container.Provide(newTokenHandler); err != nil {
		return err
	}

	return nil
}

// tokenHandlerParams TokenHandler的依赖参数
type tokenHandlerParams struct {
	dig.In

	AccessTokenSigner  auth.JwtTokenSigner `name:"accessTokenSigner"`
	RefreshTokenSigner auth.JwtTokenSigner `name:"refreshTokenSigner"`
}

// newTokenHandler 创建TokenHandler的工厂函数（用于dig）
func newTokenHandler(p tokenHandlerParams) *handler.TokenHandler {
	return handler.NewTokenHandler(p.AccessTokenSigner, p.RefreshTokenSigner)
}

// GetUserHandler 获取用户处理器
func (dc *Container) GetUserHandler() *http.UserHandler {
	var handler *http.UserHandler
	if err := dc.container.Invoke(func(h *http.UserHandler) {
		handler = h
	}); err != nil {
		panic(err)
	}
	return handler
}

// GetOAuth2Handler 获取OAuth2处理器
func (dc *Container) GetOAuth2Handler() *http.OAuth2Handler {
	var handler *http.OAuth2Handler
	if err := dc.container.Invoke(func(h *http.OAuth2Handler) {
		handler = h
	}); err != nil {
		panic(err)
	}
	return handler
}

// GetTokenHandler 获取Token处理器
func (dc *Container) GetTokenHandler() *handler.TokenHandler {
	var h *handler.TokenHandler
	if err := dc.container.Invoke(func(tokenHandler *handler.TokenHandler) {
		h = tokenHandler
	}); err != nil {
		panic(err)
	}
	return h
}

// Invoke 调用函数并注入依赖
func (dc *Container) Invoke(fn interface{}) error {
	return dc.container.Invoke(fn)
}

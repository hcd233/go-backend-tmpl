// Package auth 认证应用服务
//
//	@update 2025-09-30 00:00:00
package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hcd233/go-backend-tmpl/internal/domain/auth"
	"github.com/hcd233/go-backend-tmpl/internal/domain/user"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// ErrUnauthorized 未授权
	ErrUnauthorized = errors.New("unauthorized")
	// ErrInternalError 内部错误
	ErrInternalError = errors.New("internal error")
	// ErrInvalidProvider 无效的OAuth2提供商
	ErrInvalidProvider = errors.New("invalid OAuth2 provider")
)

// Service 认证应用服务
//
//	@update 2025-09-30 00:00:00
type Service struct {
	authService     *auth.Service
	userRepo        user.Repository
	userService     *user.Service
	oauth2Providers map[auth.Provider]auth.OAuth2Provider
	// 用于创建用户目录的回调
	onUserCreated func(ctx context.Context, userID uint) error
}

// NewService 创建认证应用服务
func NewService(
	authService *auth.Service,
	userRepo user.Repository,
	userService *user.Service,
	oauth2Providers map[auth.Provider]auth.OAuth2Provider,
	onUserCreated func(ctx context.Context, userID uint) error,
) *Service {
	return &Service{
		authService:     authService,
		userRepo:        userRepo,
		userService:     userService,
		oauth2Providers: oauth2Providers,
		onUserCreated:   onUserCreated,
	}
}

// Login 登录
func (s *Service) Login(ctx context.Context, cmd *LoginCommand) (*LoginResponse, error) {
	logger := logger.WithCtx(ctx)

	provider := auth.Provider(cmd.Provider)
	if !provider.IsValid() {
		logger.Error("[AuthAppService] invalid provider", zap.String("provider", cmd.Provider))
		return nil, ErrInvalidProvider
	}

	oauth2Provider, exists := s.oauth2Providers[provider]
	if !exists {
		logger.Error("[AuthAppService] OAuth2 provider not configured", zap.String("provider", cmd.Provider))
		return nil, ErrInvalidProvider
	}

	redirectURL := oauth2Provider.GetAuthURL()
	logger.Info("[AuthAppService] login redirect", zap.String("redirectURL", redirectURL))

	return &LoginResponse{RedirectURL: redirectURL}, nil
}

// Callback OAuth2回调处理
func (s *Service) Callback(ctx context.Context, cmd *CallbackCommand) (*CallbackResponse, error) {
	logger := logger.WithCtx(ctx)

	// 验证state
	if err := s.authService.ValidateState(cmd.State); err != nil {
		logger.Error("[AuthAppService] invalid state", zap.String("state", cmd.State))
		return nil, ErrUnauthorized
	}

	// 获取OAuth2提供商
	provider := auth.Provider(cmd.Provider)
	oauth2Provider, exists := s.oauth2Providers[provider]
	if !exists {
		logger.Error("[AuthAppService] OAuth2 provider not found", zap.String("provider", cmd.Provider))
		return nil, ErrInvalidProvider
	}

	// 交换授权码获取token
	token, err := oauth2Provider.ExchangeToken(ctx, cmd.Code)
	if err != nil {
		logger.Error("[AuthAppService] failed to exchange token", zap.Error(err))
		return nil, ErrUnauthorized
	}

	// 获取OAuth2用户信息
	oauth2UserInfo, err := oauth2Provider.GetUserInfo(ctx, token)
	if err != nil {
		logger.Error("[AuthAppService] failed to get user info", zap.Error(err))
		return nil, ErrInternalError
	}

	logger.Info("[AuthAppService] OAuth2 user info",
		zap.String("id", oauth2UserInfo.ID()),
		zap.String("name", oauth2UserInfo.Name()),
		zap.String("email", oauth2UserInfo.Email()))

	// 查找或创建用户
	email, err := user.NewEmail(oauth2UserInfo.Email())
	if err != nil {
		logger.Error("[AuthAppService] invalid email", zap.Error(err))
		return nil, ErrInternalError
	}

	u, err := s.userRepo.FindByEmail(ctx, *email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[AuthAppService] failed to find user", zap.Error(err))
		return nil, ErrInternalError
	}

	// 用户存在，记录登录
	if u != nil {
		u.RecordLogin()

		// 绑定第三方账号
		switch provider {
		case auth.ProviderGithub:
			_ = u.BindGithub(oauth2UserInfo.ID())
		case auth.ProviderGoogle:
			_ = u.BindGoogle(oauth2UserInfo.ID())
		}

		if err := s.userRepo.Save(ctx, u); err != nil {
			logger.Error("[AuthAppService] failed to save user", zap.Error(err))
			return nil, ErrInternalError
		}
	} else {
		// 创建新用户
		userName := oauth2UserInfo.Name()
		if err := user.ValidateUserName(userName); err != nil {
			// 如果用户名不符合规则，生成一个新的
			userName = "User" + strconv.FormatInt(time.Now().Unix(), 10)
		}

		u, err = s.userService.CreateUser(ctx, userName, oauth2UserInfo.Email(), oauth2UserInfo.Avatar(), user.PermissionReader)
		if err != nil {
			logger.Error("[AuthAppService] failed to create user", zap.Error(err))
			return nil, ErrInternalError
		}

		// 绑定第三方账号
		switch provider {
		case auth.ProviderGithub:
			_ = u.BindGithub(oauth2UserInfo.ID())
		case auth.ProviderGoogle:
			_ = u.BindGoogle(oauth2UserInfo.ID())
		}

		if err := s.userRepo.Save(ctx, u); err != nil {
			logger.Error("[AuthAppService] failed to save user with bind", zap.Error(err))
			return nil, ErrInternalError
		}

		// 执行用户创建后的回调（如创建用户目录）
		if s.onUserCreated != nil {
			if err := s.onUserCreated(ctx, u.ID()); err != nil {
				logger.Error("[AuthAppService] failed to execute onUserCreated callback", zap.Error(err))
				// 不返回错误，因为用户已经创建成功
			}
		}

		logger.Info("[AuthAppService] new user created", zap.Uint("userID", u.ID()))
	}

	// 生成JWT令牌
	tokens, err := s.authService.GenerateTokens(u.ID())
	if err != nil {
		logger.Error("[AuthAppService] failed to generate tokens", zap.Error(err))
		return nil, ErrInternalError
	}

	logger.Info("[AuthAppService] callback success", zap.Uint("userID", u.ID()))

	return &CallbackResponse{
		AccessToken:  tokens.AccessToken(),
		RefreshToken: tokens.RefreshToken(),
	}, nil
}

// RefreshToken 刷新令牌
func (s *Service) RefreshToken(ctx context.Context, cmd *RefreshTokenCommand) (*RefreshTokenResponse, error) {
	logger := logger.WithCtx(ctx)

	// 这里简化处理，实际应该验证refresh token并生成新的access token
	// 由于原代码中没有刷新token的逻辑，这里暂时返回错误
	logger.Error("[AuthAppService] refresh token not implemented")
	return nil, fmt.Errorf("refresh token not implemented")
}


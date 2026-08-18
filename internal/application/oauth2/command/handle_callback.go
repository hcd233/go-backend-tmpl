// Package command 定义 OAuth2 域命令处理器。
package command

import (
	"context"
	"strconv"
	"time"

	oauth2port "github.com/hcd233/aris-api-tmpl/internal/application/oauth2/port"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/common/ierr"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity"
	identityaggregate "github.com/hcd233/aris-api-tmpl/internal/domain/identity/aggregate"
	identityservice "github.com/hcd233/aris-api-tmpl/internal/domain/identity/service"
	identityvo "github.com/hcd233/aris-api-tmpl/internal/domain/identity/vo"
	oauth2service "github.com/hcd233/aris-api-tmpl/internal/domain/oauth2/service"
	oauth2vo "github.com/hcd233/aris-api-tmpl/internal/domain/oauth2/vo"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/hcd233/aris-api-tmpl/internal/util"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type initiateLoginHandler struct {
	platforms map[string]oauth2service.Platform
}

// NewInitiateLoginHandler 构造登录发起处理器。
func NewInitiateLoginHandler(platforms map[string]oauth2service.Platform) oauth2port.InitiateLoginHandler {
	return &initiateLoginHandler{platforms: platforms}
}

// Handle 执行登录发起流程。
func (h *initiateLoginHandler) Handle(ctx context.Context, cmd oauth2port.InitiateLoginCommand) (*oauth2port.InitiateLoginResult, error) {
	platform, ok := h.platforms[cmd.Platform]
	if !ok {
		logger.WithCtx(ctx).Warn("[OAuth2Command] invalid platform on initiate", zap.String("platform", cmd.Platform))
		return nil, ierr.New(ierr.ErrBadRequest, "invalid oauth platform")
	}
	url := platform.GetAuthURL()
	logger.WithCtx(ctx).Info("[OAuth2Command] initiate login", zap.String("platform", cmd.Platform))
	return &oauth2port.InitiateLoginResult{RedirectURL: url}, nil
}

type handleCallbackHandler struct {
	platforms     map[string]oauth2service.Platform
	userRepo      identity.UserRepository
	accessSigner  identityservice.TokenSigner
	refreshSigner identityservice.TokenSigner
	dirCreator    oauth2port.ObjectStorageDirCreator
}

// NewHandleCallbackHandler 构造回调处理器。
func NewHandleCallbackHandler(
	platforms map[string]oauth2service.Platform,
	userRepo identity.UserRepository,
	accessSigner, refreshSigner identityservice.TokenSigner,
	dirCreator oauth2port.ObjectStorageDirCreator,
) oauth2port.HandleCallbackHandler {
	return &handleCallbackHandler{
		platforms:     platforms,
		userRepo:      userRepo,
		accessSigner:  accessSigner,
		refreshSigner: refreshSigner,
		dirCreator:    dirCreator,
	}
}

// Handle 执行 OAuth2 回调流程。
func (h *handleCallbackHandler) Handle(ctx context.Context, cmd oauth2port.HandleCallbackCommand) (*oauth2port.HandleCallbackResult, error) {
	platform, err := h.validateStateAndPlatform(ctx, cmd.State, cmd.Platform)
	if err != nil {
		return nil, err
	}
	userInfo, err := h.exchangeAndFetchUser(ctx, platform, cmd.Code)
	if err != nil {
		return nil, err
	}
	userID, isNewUser, err := h.resolveUser(ctx, cmd.Platform, userInfo)
	if err != nil {
		return nil, err
	}
	accessToken, refreshToken, err := h.signTokenPair(ctx, userID)
	if err != nil {
		return nil, err
	}
	logger.WithCtx(ctx).Info("[OAuth2Command] callback success",
		zap.String("platform", cmd.Platform), zap.Uint("userID", userID), zap.Bool("isNewUser", isNewUser))
	pair := identityvo.NewTokenPair(accessToken, refreshToken)
	return &oauth2port.HandleCallbackResult{TokenPair: &pair, UserID: userID, IsNewUser: isNewUser}, nil
}

func (h *handleCallbackHandler) validateStateAndPlatform(ctx context.Context, state, platform string) (oauth2service.Platform, error) {
	if state != config.Oauth2StateString {
		logger.WithCtx(ctx).Error("[OAuth2Command] invalid state", zap.String("platform", platform))
		return nil, ierr.New(ierr.ErrUnauthorized, "invalid oauth state")
	}
	p, ok := h.platforms[platform]
	if !ok {
		logger.WithCtx(ctx).Error("[OAuth2Command] invalid platform", zap.String("platform", platform))
		return nil, ierr.New(ierr.ErrBadRequest, "invalid oauth platform")
	}
	return p, nil
}

func (h *handleCallbackHandler) exchangeAndFetchUser(ctx context.Context, platform oauth2service.Platform, code string) (oauth2vo.OAuthUserInfo, error) {
	log := logger.WithCtx(ctx)
	log.Info("[OAuth2Command] exchanging token")
	token, err := platform.ExchangeToken(ctx, code)
	if err != nil {
		log.Error("[OAuth2Command] failed to exchange token", zap.Error(err))
		return oauth2vo.OAuthUserInfo{}, ierr.Wrap(ierr.ErrOAuth2Exchange, err, "exchange oauth token")
	}
	h.logTokenReceived(ctx, token)
	userInfo, err := platform.GetUserInfo(ctx, token)
	if err != nil {
		log.Error("[OAuth2Command] failed to get user info", zap.Error(err))
		return oauth2vo.OAuthUserInfo{}, ierr.Wrap(ierr.ErrOAuth2UserInfo, err, "get oauth user info")
	}
	return userInfo, nil
}

func (h *handleCallbackHandler) resolveUser(ctx context.Context, platformName string, userInfo oauth2vo.OAuthUserInfo) (userID uint, isNewUser bool, err error) {
	log := logger.WithCtx(ctx)
	thirdPartyID := userInfo.ID()
	userName := userInfo.Name()
	existing, err := h.findByBindID(ctx, platformName, thirdPartyID)
	if err != nil {
		log.Error("[OAuth2Command] failed to find user by bind id",
			zap.String("platform", platformName), zap.String("thirdPartyID", thirdPartyID), zap.Error(err))
		return 0, false, err
	}
	if existing != nil {
		if err := h.userRepo.TouchLastLogin(ctx, existing.AggregateID()); err != nil {
			log.Error("[OAuth2Command] failed to update user login time", zap.String("platform", platformName), zap.Error(err))
			return 0, false, err
		}
		return existing.AggregateID(), false, nil
	}
	if err := util.ValidateUserName(userName); err != nil {
		userName = constant.DefaultUserNamePrefix + strconv.FormatInt(time.Now().UTC().Unix(), 10)
	}
	user, err := identityaggregate.RegisterUser(
		identityvo.UserName(userName), identityvo.Email(userInfo.Email()), identityvo.Avatar(userInfo.Avatar()),
		platformName, thirdPartyID, time.Now().UTC(),
	)
	if err != nil {
		log.Error("[OAuth2Command] register user aggregate failed", zap.String("platform", platformName), zap.Error(err))
		return 0, false, err
	}
	if err := h.userRepo.Save(ctx, user); err != nil {
		log.Error("[OAuth2Command] failed to save new user", zap.String("platform", platformName), zap.Error(err))
		return 0, false, err
	}
	if h.dirCreator != nil {
		if err := h.dirCreator.CreateDir(ctx, user.AggregateID()); err != nil {
			log.Error("[OAuth2Command] failed to create audio dir", zap.String("platform", platformName), zap.Error(err))
			return 0, false, err
		}
	}
	return user.AggregateID(), true, nil
}

func (h *handleCallbackHandler) signTokenPair(ctx context.Context, userID uint) (accessToken, refreshToken string, err error) {
	accessToken, err = h.accessSigner.EncodeToken(userID)
	if err != nil {
		logger.WithCtx(ctx).Error("[OAuth2Command] failed to encode access token", zap.Error(err))
		return "", "", ierr.Wrap(ierr.ErrJWTEncode, err, "encode access token")
	}
	refreshToken, err = h.refreshSigner.EncodeToken(userID)
	if err != nil {
		logger.WithCtx(ctx).Error("[OAuth2Command] failed to encode refresh token", zap.Error(err))
		return "", "", ierr.Wrap(ierr.ErrJWTEncode, err, "encode refresh token")
	}
	return accessToken, refreshToken, nil
}

func (h *handleCallbackHandler) findByBindID(ctx context.Context, platform, bindID string) (*identityaggregate.User, error) {
	switch platform {
	case constant.OAuthProviderGithub:
		return h.userRepo.FindByGithubBindID(ctx, bindID)
	case constant.OAuthProviderGoogle:
		return h.userRepo.FindByGoogleBindID(ctx, bindID)
	default:
		return nil, ierr.New(ierr.ErrBadRequest, "invalid oauth platform")
	}
}

func (h *handleCallbackHandler) logTokenReceived(ctx context.Context, token *oauth2.Token) {
	logger.WithCtx(ctx).Info("[OAuth2Command] token exchange successful",
		zap.String("tokenType", token.TokenType), zap.Bool("valid", token.Valid()))
}

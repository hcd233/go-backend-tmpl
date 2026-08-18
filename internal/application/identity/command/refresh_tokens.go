package command

import (
	"context"

	identityport "github.com/hcd233/aris-api-tmpl/internal/application/identity/port"
	"github.com/hcd233/aris-api-tmpl/internal/common/ierr"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity/service"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity/vo"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/hcd233/aris-api-tmpl/internal/util"
	"go.uber.org/zap"
)

type refreshTokensHandler struct {
	repo    identity.UserRepository
	access  service.TokenSigner
	refresh service.TokenSigner
}

// NewRefreshTokensHandler 构造刷新令牌处理器。
func NewRefreshTokensHandler(repo identity.UserRepository, access, refresh service.TokenSigner) identityport.RefreshTokensHandler {
	return &refreshTokensHandler{repo: repo, access: access, refresh: refresh}
}

// Handle 执行刷新令牌流程。
func (h *refreshTokensHandler) Handle(ctx context.Context, cmd identityport.RefreshTokensCommand) (*vo.TokenPair, error) {
	log := logger.WithCtx(ctx)
	userID, err := h.refresh.DecodeToken(cmd.RefreshToken)
	if err != nil {
		log.Error("[IdentityCommand] decode refresh token failed",
			zap.String("refreshToken", util.MaskSecret(cmd.RefreshToken)), zap.Error(err))
		return nil, ierr.Wrap(ierr.ErrJWTDecode, err, "decode refresh token")
	}
	user, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		log.Error("[IdentityCommand] find user failed", zap.Error(err), zap.Uint("userID", userID))
		return nil, err
	}
	if user == nil {
		log.Warn("[IdentityCommand] user not found during refresh", zap.Uint("userID", userID))
		return nil, ierr.New(ierr.ErrDataNotExists, "user not found")
	}
	accessToken, err := h.access.EncodeToken(userID)
	if err != nil {
		log.Error("[IdentityCommand] encode access token failed", zap.Error(err))
		return nil, ierr.Wrap(ierr.ErrJWTEncode, err, "encode access token")
	}
	refreshToken, err := h.refresh.EncodeToken(userID)
	if err != nil {
		log.Error("[IdentityCommand] encode refresh token failed", zap.Error(err))
		return nil, ierr.Wrap(ierr.ErrJWTEncode, err, "encode refresh token")
	}
	log.Info("[IdentityCommand] refresh token success", zap.Uint("userID", userID))
	pair := vo.NewTokenPair(accessToken, refreshToken)
	return &pair, nil
}

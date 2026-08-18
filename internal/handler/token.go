package handler

import (
	"context"
	"strings"

	apiutil "github.com/hcd233/aris-api-tmpl/internal/api/util"
	identityport "github.com/hcd233/aris-api-tmpl/internal/application/identity/port"
	"github.com/hcd233/aris-api-tmpl/internal/common/ierr"
	"github.com/hcd233/aris-api-tmpl/internal/dto"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/hcd233/aris-api-tmpl/internal/util"
	"go.uber.org/zap"
)

// TokenHandler 令牌处理器。
type TokenHandler interface {
	HandleRefreshToken(ctx context.Context, req *dto.RefreshTokenReq) (*dto.HTTPResponse[*dto.RefreshTokenRsp], error)
}

// TokenDependencies TokenHandler 依赖项。
type TokenDependencies struct {
	Refresh identityport.RefreshTokensHandler
}

type tokenHandler struct {
	refresh identityport.RefreshTokensHandler
}

// NewTokenHandler 创建令牌处理器。
func NewTokenHandler(deps TokenDependencies) TokenHandler {
	return &tokenHandler{refresh: deps.Refresh}
}

// HandleRefreshToken 刷新令牌。
func (h *tokenHandler) HandleRefreshToken(ctx context.Context, req *dto.RefreshTokenReq) (*dto.HTTPResponse[*dto.RefreshTokenRsp], error) {
	if req == nil || req.Body == nil || strings.TrimSpace(req.Body.RefreshToken) == "" {
		return nil, apiutil.NewHumaBizErrorFromModel(ctx, ierr.ErrValidation.BizError())
	}
	pair, err := h.refresh.Handle(ctx, identityport.RefreshTokensCommand{RefreshToken: req.Body.RefreshToken})
	if err != nil {
		logger.WithCtx(ctx).Warn("[TokenHandler] refresh token failed",
			zap.String("refreshToken", util.MaskSecret(req.Body.RefreshToken)), zap.Error(err))
		return nil, apiutil.NewHumaBizError(ctx, err, ierr.ErrInternal.BizError())
	}
	rsp := &dto.RefreshTokenRsp{
		AccessToken:  pair.AccessToken(),
		RefreshToken: pair.RefreshToken(),
	}
	return util.WrapHTTPResponse(rsp, nil)
}

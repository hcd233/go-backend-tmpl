// Package handler HTTP处理器（兼容旧路由）
package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/go-backend-tmpl/internal/auth"
	"github.com/hcd233/go-backend-tmpl/internal/constant"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
	"github.com/hcd233/go-backend-tmpl/internal/util"
	"go.uber.org/zap"
)

// TokenHandler Token处理器
type TokenHandler struct {
	accessTokenSigner  auth.JwtTokenSigner
	refreshTokenSigner auth.JwtTokenSigner
}

// NewTokenHandler 创建Token处理器（使用依赖注入）
func NewTokenHandler(
	accessTokenSigner auth.JwtTokenSigner,
	refreshTokenSigner auth.JwtTokenSigner,
) *TokenHandler {
	return &TokenHandler{
		accessTokenSigner:  accessTokenSigner,
		refreshTokenSigner: refreshTokenSigner,
	}
}

// HandleRefreshToken 刷新Token
//
//	@Summary		刷新访问令牌
//	@Description	使用refresh token获取新的access token
//	@Tags			token
//	@Accept			json
//	@Produce		json
//	@Param			body	body		protocol.RefreshTokenBody	true	"刷新令牌请求"
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.RefreshTokenResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/token/refresh [post]
func (h *TokenHandler) HandleRefreshToken(c *fiber.Ctx) error {
	logger := logger.WithCtx(c.Context())
	body := c.Locals(constant.CtxKeyBody).(*protocol.RefreshTokenBody)

	// 验证refresh token
	userID, err := h.refreshTokenSigner.DecodeToken(body.RefreshToken)
	if err != nil {
		logger.Error("[TokenHandler] invalid refresh token", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrUnauthorized)
		return nil
	}

	// 生成新的access token
	accessToken, err := h.accessTokenSigner.EncodeToken(userID)
	if err != nil {
		logger.Error("[TokenHandler] failed to generate access token", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return nil
	}

	rsp := &protocol.RefreshTokenResponse{
		AccessToken: accessToken,
	}

	logger.Info("[TokenHandler] token refreshed", zap.Uint("userID", userID))
	util.SendHTTPResponse(c, rsp, nil)
	return nil
}

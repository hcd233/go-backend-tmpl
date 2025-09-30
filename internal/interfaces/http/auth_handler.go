// Package http HTTP接口层
//
//	@update 2025-09-30 00:00:00
package http

import (
	"github.com/gofiber/fiber/v2"
	appauth "github.com/hcd233/go-backend-tmpl/internal/application/auth"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
	"github.com/hcd233/go-backend-tmpl/internal/util"
)

// OAuth2Handler OAuth2 HTTP处理器
type OAuth2Handler struct {
	authAppService *appauth.Service
}

// NewOAuth2Handler 创建OAuth2 HTTP处理器
func NewOAuth2Handler(authAppService *appauth.Service) *OAuth2Handler {
	return &OAuth2Handler{
		authAppService: authAppService,
	}
}

// HandleLogin 处理登录请求
//
//	@Summary		OAuth2登录
//	@Description	获取OAuth2授权URL
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string	true	"OAuth2提供商 (github/google)"
//	@Success		200			{object}	protocol.HTTPResponse{data=appauth.LoginResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/oauth2/{provider}/login [get]
func (h *OAuth2Handler) HandleLogin(c *fiber.Ctx) error {
	provider := c.Params("provider")

	cmd := &appauth.LoginCommand{
		Provider: provider,
	}

	rsp, err := h.authAppService.Login(c.Context(), cmd)
	if err != nil {
		protocolErr := convertAppErrorToProtocol(err)
		util.SendHTTPResponse(c, nil, protocolErr)
		return nil
	}

	util.SendHTTPResponse(c, rsp, nil)
	return nil
}

// HandleCallback 处理OAuth2回调
//
//	@Summary		OAuth2回调
//	@Description	处理OAuth2授权回调
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string	true	"OAuth2提供商 (github/google)"
//	@Param			code		query		string	true	"授权码"
//	@Param			state		query		string	true	"状态码"
//	@Success		200			{object}	protocol.HTTPResponse{data=appauth.CallbackResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/oauth2/{provider}/callback [get]
func (h *OAuth2Handler) HandleCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	param := c.Locals("param").(*protocol.OAuth2CallbackParam)

	cmd := &appauth.CallbackCommand{
		Code:     param.Code,
		State:    param.State,
		Provider: provider,
	}

	rsp, err := h.authAppService.Callback(c.Context(), cmd)
	if err != nil {
		protocolErr := convertAppErrorToProtocol(err)
		util.SendHTTPResponse(c, nil, protocolErr)
		return nil
	}

	util.SendHTTPResponse(c, rsp, nil)
	return nil
}

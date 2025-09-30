package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
	"github.com/hcd233/go-backend-tmpl/internal/service"
	"github.com/hcd233/go-backend-tmpl/internal/util"
)

// Oauth2Handler OAuth2处理器接口
type Oauth2Handler interface {
	HandleLogin(c *fiber.Ctx) error
	HandleCallback(c *fiber.Ctx) error
}

type oauth2Handler struct {
	svc service.Oauth2Service
}

// NewGithubOauth2Handler 创建Github OAuth2处理器
//
//	return Oauth2Handler
//	author centonhuang
//	update 2025-01-05 13:43:43
func NewGithubOauth2Handler() Oauth2Handler {
	return &oauth2Handler{
		svc: service.NewGithubOauth2Service(),
	}
}

// NewGoogleOauth2Handler 创建Google OAuth2处理器
func NewGoogleOauth2Handler() Oauth2Handler {
	return &oauth2Handler{
		svc: service.NewGoogleOauth2Service(),
	}
}

// HandleLogin OAuth2登录
//
//	@Summary		OAuth2登录
//	@Description	OAuth2登录请求,返回重定向URL
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	protocol.HTTPResponse{data=protocol.LoginResponse,error=nil}
//	@Failure		400	{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401	{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403	{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500	{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/oauth2/{provider}/login [get]
//	receiver h *oauth2Handler
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-05 13:43:42
func (h *oauth2Handler) HandleLogin(c *fiber.Ctx) error {
	req := &protocol.LoginRequest{}

	rsp, err := h.svc.Login(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

// HandleCallback OAuth2回调
//
//	@Summary		OAuth2回调
//	@Description	OAuth2回调请求,验证code和state
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"授权码"
//	@Param			state	query		string	true	"状态码"
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.CallbackResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/oauth2/{provider}/callback [get]
//	receiver h *oauth2Handler
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-05 13:43:36
func (h *oauth2Handler) HandleCallback(c *fiber.Ctx) error {
	params := protocol.OAuth2CallbackParam{}
	if err := c.QueryParser(&params); err != nil {
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return nil
	}

	req := &protocol.CallbackRequest{
		Code:  params.Code,
		State: params.State,
	}

	rsp, err := h.svc.Callback(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

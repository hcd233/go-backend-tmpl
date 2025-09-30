// Package http HTTP接口层
//
//	@update 2025-09-30 00:00:00
package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/go-backend-tmpl/internal/application/user"
	"github.com/hcd233/go-backend-tmpl/internal/constant"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
	"github.com/hcd233/go-backend-tmpl/internal/util"
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	userAppService *user.Service
}

// NewUserHandler 创建用户HTTP处理器
func NewUserHandler(userAppService *user.Service) *UserHandler {
	return &UserHandler{
		userAppService: userAppService,
	}
}

// HandleGetCurUserInfo 获取当前用户信息
//
//	@Summary		获取当前用户信息
//	@Description	获取当前用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=user.CurUserInfoResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/user/current [get]
func (h *UserHandler) HandleGetCurUserInfo(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)

	query := &user.GetCurUserInfoQuery{
		UserID: userID,
	}

	rsp, err := h.userAppService.GetCurUserInfo(c.Context(), query)
	if err != nil {
		// 转换应用层错误为protocol错误
		protocolErr := convertAppErrorToProtocol(err)
		util.SendHTTPResponse(c, nil, protocolErr)
		return nil
	}

	util.SendHTTPResponse(c, rsp, nil)
	return nil
}

// HandleGetUserInfo 获取用户信息
//
//	@Summary		获取用户信息
//	@Description	获取用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userID	path		uint	true	"用户ID"
//	@Success		200		{object}	protocol.HTTPResponse{data=user.UserInfoResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/user/{userID} [get]
func (h *UserHandler) HandleGetUserInfo(c *fiber.Ctx) error {
	uri := c.Locals(constant.CtxKeyURI).(*protocol.UserURI)

	query := &user.GetUserInfoQuery{
		UserID: uri.UserID,
	}

	rsp, err := h.userAppService.GetUserInfo(c.Context(), query)
	if err != nil {
		protocolErr := convertAppErrorToProtocol(err)
		util.SendHTTPResponse(c, nil, protocolErr)
		return nil
	}

	util.SendHTTPResponse(c, rsp, nil)
	return nil
}

// HandleUpdateInfo 更新用户信息
//
//	@Summary		更新用户信息
//	@Description	更新用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body	body		protocol.UpdateUserBody	true	"更新用户信息请求"
//	@Success		200		{object}	protocol.HTTPResponse{data=user.UpdateUserInfoResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/user [patch]
func (h *UserHandler) HandleUpdateInfo(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)
	body := c.Locals(constant.CtxKeyBody).(*protocol.UpdateUserBody)

	cmd := &user.UpdateUserInfoCommand{
		UserID:      userID,
		UpdatedName: body.UserName,
	}

	rsp, err := h.userAppService.UpdateUserInfo(c.Context(), cmd)
	if err != nil {
		protocolErr := convertAppErrorToProtocol(err)
		util.SendHTTPResponse(c, nil, protocolErr)
		return nil
	}

	util.SendHTTPResponse(c, rsp, nil)
	return nil
}


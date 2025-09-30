package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
	"github.com/hcd233/go-backend-tmpl/internal/util"
)

// PingHandler 健康检查处理器
//
//	author centonhuang
//	update 2025-01-04 15:52:48
type PingHandler interface {
	HandlePing(c *fiber.Ctx) error
}

type pingHandler struct{}

// NewPingHandler 创建健康检查处理器
//
//	return PingHandler
//	author centonhuang
//	update 2025-01-04 15:52:48
func NewPingHandler() PingHandler {
	return &pingHandler{}
}

// HandlePing 处理健康检查请求
//
//	@Summary		健康检查
//	@Description	健康检查
//	@Tags			ping
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.PingResponse,error=nil}
//	@Router			/ [get]
//	receiver h *pingHandler
//	param c *fiber.Ctx
//	author centonhuang
//	update 2025-01-04 20:47:48
func (h *pingHandler) HandlePing(c *fiber.Ctx) error {
	rsp := protocol.PingResponse{
		Status: "ok",
	}

	util.SendHTTPResponse(c, rsp, nil)
	return nil
}

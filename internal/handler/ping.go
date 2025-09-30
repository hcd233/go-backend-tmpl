// Package handler HTTP处理器（兼容旧路由）
package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
)

// PingHandler Ping处理器
type PingHandler struct{}

// NewPingHandler 创建Ping处理器
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// HandlePing 处理Ping请求
//
//	@Summary		Ping
//	@Description	Health check endpoint
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	protocol.HTTPResponse{data=string,error=nil}
//	@Router			/ [get]
func (h *PingHandler) HandlePing(c *fiber.Ctx) error {
	return c.JSON(protocol.HTTPResponse{
		Data: "ok",
	})
}

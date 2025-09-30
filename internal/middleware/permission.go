package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/go-backend-tmpl/internal/constant"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
	"github.com/hcd233/go-backend-tmpl/internal/resource/database/model"
	"github.com/hcd233/go-backend-tmpl/internal/util"
	"go.uber.org/zap"
)

// LimitUserPermissionMiddleware 限制用户权限中间件
//
//	param serviceName string
//	param requiredPermission model.Permission
//	return fiber.Handler
//	author centonhuang
//	update 2025-01-05 15:07:08
func LimitUserPermissionMiddleware(serviceName string, requiredPermission model.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permission := c.Locals(constant.CtxKeyPermission).(model.Permission)
		if model.PermissionLevelMapping[permission] < model.PermissionLevelMapping[requiredPermission] {
			logger.WithFCtx(c).Info("[LimitUserPermissionMiddleware] permission denied",
				zap.String("serviceName", serviceName),
				zap.String("requiredPermission", string(requiredPermission)),
				zap.String("permission", string(permission)))
			util.SendHTTPResponse(c, nil, protocol.ErrNoPermission)
			return c.Status(fiber.StatusForbidden).JSON(protocol.HTTPResponse{
				Error: protocol.ErrNoPermission.Error(),
			})
		}

		return c.Next()
	}
}

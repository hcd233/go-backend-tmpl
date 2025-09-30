package middleware

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"go.uber.org/zap"
)

// RecoverMiddleware 恢复中间件
//
//	@return fiber.Handler
//	@author centonhuang
//	@update 2025-08-18 20:21:14
func RecoverMiddleware() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			logger.WithFCtx(c).Panic("[Panic Recovery] recovered panic",
				zap.Any("error", e),
				zap.ByteString("stack", debug.Stack()))
		},
	})
}

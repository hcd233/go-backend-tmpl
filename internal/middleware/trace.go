package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hcd233/go-backend-tmpl/internal/constant"
)

// TraceMiddleware 追踪中间件
//
//	return fiber.Handler
//	author centonhuang
//	update 2025-01-05 15:30:00
func TraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("X-Trace-Id")

		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Locals(constant.CtxKeyTraceID, traceID)

		c.Set("X-Trace-Id", traceID)

		return c.Next()
	}
}

package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/i18n"
)

// LocaleMiddleware 从 Accept-Language 头检测语言环境并写入请求上下文，
// 供错误本地化链路（ierr → model.Error.Localize）消费。
func LocaleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		locale := i18n.DetectLocale(c.Get(constant.HTTPHeaderAcceptLanguage))
		c.Locals(constant.CtxKeyLocale, locale)
		return c.Next()
	}
}

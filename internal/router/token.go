package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/go-backend-tmpl/internal/config"
	"github.com/hcd233/go-backend-tmpl/internal/infrastructure/container"
	"github.com/hcd233/go-backend-tmpl/internal/middleware"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
)

// initTokenRouter 初始化Token路由（使用依赖注入）
func initTokenRouter(r fiber.Router, c *container.Container) {
	tokenHandler := c.GetTokenHandler()

	tokenRouter := r.Group("/token")
	{
		tokenRouter.Post(
			"/refresh",
			middleware.RateLimiterMiddleware("refreshToken", "", config.JwtAccessTokenExpired/4, 2),
			middleware.ValidateBodyMiddleware(&protocol.RefreshTokenBody{}),
			tokenHandler.HandleRefreshToken,
		)
	}
}

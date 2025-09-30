// Package router 路由器
//
//	@update 2025-09-30 00:00:00
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/hcd233/go-backend-tmpl/internal/handler"
	"github.com/hcd233/go-backend-tmpl/internal/infrastructure/container"
	"github.com/hcd233/go-backend-tmpl/internal/middleware"
	"github.com/hcd233/go-backend-tmpl/internal/protocol"
	"github.com/samber/lo"
)

// RegisterRouter 注册所有路由
//
//	param app *fiber.App
//	author centonhuang
//	update 2025-01-04 15:32:40
func RegisterRouter(app *fiber.App) {
	// swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Ping健康检查
	pingHandler := handler.NewPingHandler()
	app.Get("/", pingHandler.HandlePing)

	// 创建DI容器
	digContainer := lo.Must1(container.NewContainer())

	// v1 API路由组
	v1Router := app.Group("/v1")
	{
		// 用户路由
		initUserRouter(v1Router, digContainer)
		// OAuth2路由
		initOAuth2Router(v1Router, digContainer)
		// Token路由（使用依赖注入）
		initTokenRouter(v1Router, digContainer)
	}
}

// initUserRouter 初始化用户路由
func initUserRouter(r fiber.Router, c *container.Container) {
	userHandler := c.GetUserHandler()

	userRouter := r.Group("/user", middleware.JwtMiddleware())
	{
		userRouter.Get("/current", userHandler.HandleGetCurUserInfo)
		userRouter.Patch("/", middleware.ValidateBodyMiddleware(&protocol.UpdateUserBody{}), userHandler.HandleUpdateInfo)

		userNameRouter := userRouter.Group("/:userID", middleware.ValidateURIMiddleware(&protocol.UserURI{}))
		{
			userNameRouter.Get("/", userHandler.HandleGetUserInfo)
		}
	}
}

// initOAuth2Router 初始化OAuth2路由
func initOAuth2Router(r fiber.Router, c *container.Container) {
	oauth2Handler := c.GetOAuth2Handler()

	oauth2Router := r.Group("/oauth2")
	{
		providerRouter := oauth2Router.Group("/:provider")
		{
			providerRouter.Get("/login", oauth2Handler.HandleLogin)
			providerRouter.Get("/callback", middleware.ValidateParamMiddleware(&protocol.OAuth2CallbackParam{}), oauth2Handler.HandleCallback)
		}
	}
}

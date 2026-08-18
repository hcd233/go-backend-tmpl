package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	apiutil "github.com/hcd233/aris-api-tmpl/internal/api/util"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/enum"
	"github.com/samber/lo"
)

func init() {
	// 统一 200 错误契约：huma 框架错误（422/404 等）也转换为
	// 顶层 {"error": {code, message}} 结构，与 handler 错误、中间件错误一致。
	huma.NewError = apiutil.FrameworkError
}

// NewHumaAPI 创建 Huma API 实例。
func NewHumaAPI(app *fiber.App) huma.API {
	return humafiber.New(app, huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:       "Aris API Tmpl",
				Description: "Aris API Tmpl is a RESTful API Template.",
				Version:     "1.0",
				Contact: &huma.Contact{
					Name:  "hcd233",
					Email: "lvlvko233@qq.com",
					URL:   "https://github.com/hcd233",
				},
				License: &huma.License{
					Name: "Apache 2.0",
					URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
				},
			},
			Components: &huma.Components{
				Schemas: huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer),
				SecuritySchemes: map[string]*huma.SecurityScheme{
					"jwtAuth": {
						Type:        "apiKey",
						Name:        "Authorization",
						In:          "header",
						Description: "JWT Authentication，Please pass the JWT token in the Authorization header.",
					},
				},
			},
		},
		OpenAPIPath:   lo.If(config.Env != enum.EnvProduction, "/openapi").Else(""),
		DocsPath:      "",
		SchemasPath:   lo.If(config.Env != enum.EnvProduction, "/schemas").Else(""),
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
	})
}

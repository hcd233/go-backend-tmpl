package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/go-backend-tmpl/internal/config"
	"github.com/hcd233/go-backend-tmpl/internal/cron"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"go.uber.org/zap"

	"github.com/hcd233/go-backend-tmpl/internal/middleware"
	"github.com/hcd233/go-backend-tmpl/internal/resource/cache"
	"github.com/hcd233/go-backend-tmpl/internal/resource/database"
	"github.com/hcd233/go-backend-tmpl/internal/resource/llm"
	"github.com/hcd233/go-backend-tmpl/internal/resource/storage"
	"github.com/hcd233/go-backend-tmpl/internal/router"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "服务器命令组",
	Long:  `包含服务器相关操作的命令组`,
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "启动API服务器",
	Long:  `启动并运行API服务器，监听指定的主机和端口`,
	Run: func(cmd *cobra.Command, _ []string) {
		defer func() {
			if r := recover(); r != nil {
				logger.Logger().Error("[Server] Start server panic", zap.Any("error", r), zap.ByteString("stack", debug.Stack()))
				os.Exit(1)
			}
			os.Exit(0)
		}()
		host, port := lo.Must1(cmd.Flags().GetString("host")), lo.Must1(cmd.Flags().GetString("port"))

		database.InitDatabase()
		cache.InitCache()
		storage.InitObjectStorage()
		llm.InitOpenAIClient()
		cron.InitCronJobs()

		app := fiber.New(fiber.Config{
			Prefork:      false,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  120 * time.Second,
			JSONEncoder:  sonic.Marshal,
			JSONDecoder:  sonic.Unmarshal,
		})

		// 中间件
		app.Use(
			middleware.FgprofMiddleware(),
			middleware.TraceMiddleware(),
			middleware.LogMiddleware(),
			middleware.CORSMiddleware(),
			middleware.CompressMiddleware(),
			middleware.RecoverMiddleware(),
		)

		router.RegisterRouter(app)

		lo.Must0(app.Listen(fmt.Sprintf("%s:%s", host, port)))
	},
}

func init() {
	serverCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(serverCmd)

	startServerCmd.Flags().StringP("host", "", "localhost", "监听的主机")
	startServerCmd.Flags().StringP("port", "p", "8080", "监听的端口")
}

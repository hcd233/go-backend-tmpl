package cmd

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/hcd233/aris-api-tmpl/internal/bootstrap"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server Command Group",
	Long:  `Server command group for starting and managing the API server`,
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the API server",
	Long:  `Start and run the API server, listening on the specified host and port`,
	Run: func(cmd *cobra.Command, _ []string) {
		defer func() {
			if r := recover(); r != nil {
				logger.Logger().Error("[Server] Start server panic", zap.Any("error", r), zap.ByteString("stack", debug.Stack()))
				os.Exit(1)
			}
		}()
		host, port := lo.Must1(cmd.Flags().GetString("host")), lo.Must1(cmd.Flags().GetString("port"))

		logger.Logger().Info("[Server] Environment",
			zap.String("env", config.Env),
			zap.Duration("readTimeout", config.ReadTimeout),
			zap.Duration("writeTimeout", config.WriteTimeout),
			zap.Int("maxHeaderBytes", config.MaxHeaderBytes),
			zap.Int("poolWorkers", config.PoolWorkers),
			zap.Int("poolQueueSize", config.PoolQueueSize),
			zap.Strings("trustedProxies", config.TrustedProxies),
		)

		// fx 应用：依赖注入 + 生命周期钩子（中间件/路由/优雅退出均在 Invoke 中注册）。
		app := bootstrap.BuildFxApp(host, port)

		startCtx, startCancel := context.WithTimeout(context.Background(), constant.ShutdownTimeout)
		defer startCancel()
		if err := app.Start(startCtx); err != nil {
			logger.Logger().Error("[Server] Start server failed", zap.Error(err))
			os.Exit(1)
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		logger.Logger().Info("[Server] Received shutdown signal, starting graceful shutdown...", zap.String("signal", sig.String()))

		// Stop 执行 fx OnStop 钩子（逆序）：Cron → Inflight → Pool → HTTP → Logger → DB → Redis
		stopCtx, stopCancel := context.WithTimeout(context.Background(), constant.ShutdownTimeout)
		defer stopCancel()
		if err := app.Stop(stopCtx); err != nil {
			logger.Logger().Error("[Server] Graceful shutdown failed", zap.Error(err))
			os.Exit(1)
		}
		logger.Logger().Info("[Server] Graceful shutdown completed")
	},
}

func init() {
	serverCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(serverCmd)

	startServerCmd.Flags().StringP("host", "", "localhost", "监听的主机")
	startServerCmd.Flags().StringP("port", "p", "8080", "监听的端口")
}

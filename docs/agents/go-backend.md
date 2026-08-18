# Go 后端编码契约

> **使用场景**：编写或修改 Go 后端代码时加载。涵盖测试、代码风格、Context、DTO/API、路由命名全部硬约束。

## 测试契约

- 单元测试放在源码同目录同包（`internal/<package>/<file>_test.go`）；集成/端到端测试放在 `test/<topic>/`。
- 测试数据放对应目录的 `fixtures/` 或 Go 标准 `testdata/`。
- 测试 helper 必须调用 `t.Helper()`；优先表驱动测试；断言失败信息必须包含上下文。
- 禁止 `time.Sleep()` 做同步；使用 channel、WaitGroup 或 deadline。

## 代码契约

- 业务错误创建/包装统一走 `internal/common/ierr`；禁止 `fmt.Errorf` 或 `errors.New`（lint 的 forbidigo 强制）。
- 内部链路（command/query/domain service）使用 `ierr.Wrap(sentinel, cause, msg)` 或 `ierr.New(sentinel, msg)` 传递错误。
- Handler 错误出口：`return nil, apiutil.NewHumaBizError(ctx, err, ierr.ErrInternal.BizError())`；校验错误用 `NewHumaBizErrorFromModel`。禁止手动写 `rsp.Error = ...` 再返回 200。
- 序列化统一 `github.com/bytedance/sonic`；禁止 `encoding/json`（lint 的 depguard 强制）。
- 日志使用 `logger.WithCtx(ctx)` 或 `logger.WithFCtx(c)`；消息前缀为 `[PascalCaseModule]`；key/token/secret/password 必须掩码（`util.MaskSecret` / `util.MaskHTTPHeadersForLog`）。
- Redis key、存储路径、字符串模板放 `internal/common/constant/`；禁止业务包内定义 `const` 块。
- HTTP 状态码使用 `fiber.StatusXxx`，禁止裸数字。
- DTO 时间字段用 `time.Time`；禁止 Service 层提前格式化为字符串。
- DTO 包禁止导入 `internal/infrastructure/database/model`（lint 的 depguard 强制）。

## Context 契约

- handler/service/middleware/DAO 调用链应从上层传递 `context.Context`。
- 接口逻辑层禁止随意创建 `context.Background()` 或 `context.TODO()`。
- 允许根 context 的场景：启动、基础设施初始化、cron 入口、命令行一次性任务。
- 异步协程池任务必须使用 `util.CopyContextValues(ctx)`，禁止直接持有原始请求 context。
- context key 统一使用 `constant.CtxKey*` 强类型（禁止裸字符串，lint 的 staticcheck SA1029 强制），新 key 必须注册到 `internal/common/constant/ctx.go`。

## 统一 200 错误契约

- 管理 API 一律返回 HTTP 200，错误语义由顶层 `{"error": {code, message}}` 体承载，前端只判断 `body.error`。
- 业务码定义在 `internal/common/constant/error.go`；`model.Error.StatusCode()` 负责业务码 → HTTP 状态码推导，`apiutil.FrameworkError` 负责框架错误（422/404）反向推导业务码——新增业务码时必须同步维护两个方向的映射。
- 错误消息支持 i18n：哨兵错误带 `MessageKey`（`model.NewErrorWithKey`），`LocaleMiddleware` 注入 locale 后由 `Localize` 翻译；翻译表在 `internal/i18n/locales/`。
- 例外：inflight 503（优雅退出流量摘除信号）不受统一 200 契约约束。

## 路由命名规范

- 资源名单数小写：`/user`、`/session`；禁止裸尾斜杠。
- 操作通过 Path 段表达（`/user/current`），不依赖 HTTP Method 表达语义差异。
- 路由在 `internal/router/<module>.go` 的 `initXxxRouter` 注册；中间件在 router 层装配（Group 级 `UseMiddleware` / Operation 级 `Middlewares`）。
- 新路由路径常量优先放 `internal/common/constant/route.go`。

## 依赖注入（fx）

- 依赖统一由 fx 容器注入（`internal/bootstrap/modules/` 按 infra/repository/application/handler/cron 拆分）。
- handler 层只依赖 `application/<module>/port` 接口，不依赖 command/query 具体实现。
- 仓储接口定义在 `domain/<module>/repository.go`，实现在 `infrastructure/repository/`；应用级单例配置可在 `application/<module>/port` 定义仓储接口，不硬造 domain 聚合。
- 生命周期钩子注册在 `internal/bootstrap/lifecycle.go`；停止顺序 Cron → Inflight → Pool → HTTP → Logger → DB → Redis。

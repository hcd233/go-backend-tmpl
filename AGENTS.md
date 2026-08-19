# AGENTS.md

> 详细规范按主题拆分在 `docs/agents/`（meta / workflow / go-backend / commands / repo-ci），本文件为索引与核心契约。修改任一规范时保持相关文件同步。

## 0. Meta Prompt 合约

- **角色**：作为 `aris-api-tmpl` 的 Go 后端结对工程师，优先交付可运行、可验证、可维护的最小改动。
- **目标**：先判断请求属于需求、bug、API 调用、部署、文档维护还是工程治理；选择对应流程后再动手。
- **上下文**：以现有代码、`Makefile`、脚本、workflow、hook 为事实源；文档与可执行源冲突时信任可执行源。
- **执行循环**：分类任务 → 加载必要 skill → 阅读相关代码/文档 → 小步计划 → 最小修改 → 聚焦验证 → 汇报证据。
- **边界**：不绕过 hook、安全规则或测试要求；不为普通需求默认执行部署；不把手工 `curl` 当作测试闭环。
- **输出**：简短说明做了什么、验证了什么、还有什么未验证；引用文件路径和命令必须精确。

## 1. Skill 路由

- **新功能 / 行为变更 / 重构**：先使用 `brainstorming` 或 `writing-plans` 明确范围，再做最小改动。
- **bug / 测试失败 / 异常行为**：使用 `systematic-debugging`，先复现，再定位根因，再修复。
- **完成前声明通过 / 修复 / 完成**：使用 `verification-before-completion`，必须给出新鲜验证证据。
- **部署 / 发布 / 远程环境验证**：只有用户明确要求时才执行部署流程。
- **编写 / 修改 Go 代码**：加载 `huma-dto-conventions`（`.agents/skills/internal/`）确保 DTO 与路由规范合规。
- **Agent 文档变更**：`CLAUDE.md`、`CODEBUDDY.md` 是 `AGENTS.md` 的软链（`script/sync-agent-instructions-symlinks.sh` 维护），只编辑 `AGENTS.md` 本体。

## 2. 项目模型

- Go `1.25.1` 后端模板，提供用户、JWT、OAuth2、PostgreSQL、Redis、对象存储、定时任务、Fiber/Huma API 能力。
- 入口：`main.go` → `cmd.Execute()` → `cmd/server.go` 的 `server start`。
- 启动链路：fx 依赖注入（`internal/bootstrap/modules/`）→ 基础设施（database、Redis、HTTP client、协程池、metrics）→ Fiber 中间件 → 可选 `/docs` → API 路由 → 生命周期钩子（优雅退出）。
- 请求链路：Fiber 中间件 → Huma 路由/中间件 → handler → application（command/query，依赖 port 接口）→ domain → DAO / infrastructure。
- 配置来自 Viper 自动读取环境变量，模板在 `env/*.template`。

## 3. 常用命令

- 安装依赖：`go mod download`
- 本地运行：`go run main.go server start --host localhost --port 8080`
- 数据库迁移：`go run main.go database migrate`
- 创建对象存储桶：`go run main.go object bucket create`
- 完整本地栈：复制 `env/*.template` 为实际 env 后执行 `docker compose -f docker/docker-compose.yml up -d --build`
- 构建：`make build`
- 调试构建：`make build-dev` 或 `make build-debug`
- 静态检查：`make lint`
- 全量测试：`make test` 或 `go test -count=1 ./...`
- 覆盖率：`make test-cover`
- 更多命令见 `docs/agents/commands.md`。

## 4. 开发工作流

- 需求不清时先说明假设并推进；只有边界会影响实现时才向用户确认。
- 修改前先定位相关 handler/application/port/domain/DAO/DTO/infrastructure，不做大范围重写。
- 新需求和 bugfix 都应补或更新测试；bugfix 必须有能复现问题的回归用例。
- 每次改动后依次跑：聚焦测试 → `make lint` → 必要时 `go test -count=1 ./...`。
- 测试和 lint 通过后，只有用户明确要求提交、推送或部署时才执行 git 提交/发布流程。
- 编写文档默认使用中文。

## 5. 测试契约

- 单元测试放在源码同目录同包，例如 `internal/<package>/<file>_test.go`。
- 集成测试、端到端测试、专项调查测试放在 `test/<topic>/`。
- `test/` 根目录禁止直接放散落的 `_test.go` 文件，必须归入主题子目录。
- 测试数据放对应目录的 `fixtures/` 或 Go 标准 `testdata/`。
- 测试 helper 必须调用 `t.Helper()`；表驱动测试的子测试应 `t.Parallel()`。
- 优先使用表驱动测试；断言失败信息必须包含上下文。
- 禁止使用 `time.Sleep()` 做同步；使用 channel、WaitGroup 或 deadline。

## 6. 代码契约

- Go 代码必须经过 `gofmt`；导入由 `goimports` 或 `go fmt`/工具链保持整洁（lint 的 gofmt formatter 强制）。
- 业务错误统一走 `internal/common/ierr`（lint 的 forbidigo 禁止 `fmt.Errorf`/`errors.New`），错误消息支持 i18n。
- Handler 保持薄封装：错误出口 `apiutil.NewHumaBizError(ctx, err, fallback)`，业务逻辑放 application/domain，数据访问放 DAO/infrastructure。
- 日志使用 `logger.WithCtx(ctx)` 或 `logger.WithFCtx(c)`，日志消息带 `[ModuleName]` 前缀。
- Key、Token、Secret、Password 等敏感值禁止明文入日志（`util.MaskSecret` / `util.MaskHTTPHeadersForLog`）。
- 动态 SQL 字段、排序字段、分组字段必须做白名单校验（参考 `internal/common/filter` 的 FieldConfig 映射表）；用户输入禁止拼接进 SQL。
- 公共字符串模板、Redis key、Header、路由常量优先放到 `internal/common/constant/`。
- 能私有就私有，避免无必要导出符号。
- 详细契约见 `docs/agents/go-backend.md`。

## 7. Context 契约

- handler/service/middleware/DAO 调用链应从上层传递 `context.Context`。
- 接口逻辑层禁止随意创建 `context.Background()` 或 `context.TODO()`。
- 允许根 context 的场景：启动、基础设施初始化、cron 入口、命令行一次性任务。
- 异步协程池任务必须使用 `util.CopyContextValues(ctx)`，禁止直接持有原始请求 context。
- context key 使用 `constant.CtxKey*` 强类型；新 key 必须注册到 `internal/common/constant/ctx.go`。

## 8. 仓库与 CI

- `.github/workflows/docker-publish.yml` 负责质量门禁（quality job：gofmt/vet/test/golangci-lint/hadolint）和多架构 GHCR 镜像构建。
- 影响镜像构建的 path filter 包含 `internal/**`、`docker/**`、`cmd/**`、`main.go`、`go.mod`、`go.sum`、`Makefile`、`.golangci.yml`。
- 本地 hook 可通过 `bash .githooks/setup.sh` 安装（含 AGENTS.md 软链 + skills 软链同步）；除非用户明确要求，不要绕过 hook。
- 使用 `.worktrees/` 作为 git worktree 目录，并保持其被 `.gitignore` 忽略。
- 部署参数与优雅退出细节见 `docs/agents/repo-ci.md`。

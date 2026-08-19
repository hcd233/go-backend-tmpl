# 开发工作流

> **使用场景**：开始实现需求或修复 bug 前加载。定义分支、提交、验证与排障流程。

## 分支与提交

- 开发在 `.worktrees/` 下的 git worktree 中进行（`git worktree add .worktrees/<branch> -b <branch>`）。
- 分支命名：`{feature|bugfix|refactor|chore|docs|test|hotfix}/{描述}-{yyyy-mm-dd}`。
- 提交信息：`{type}({scope}): {描述}`，如 `feat(bootstrap): implement DDD dependency injection with fx`。
- 只有用户明确要求提交、推送或部署时才执行 git 提交/发布流程。

## 本地 hook

- `bash .githooks/setup.sh` 安装 pre-commit（含 AGENTS.md 软链同步 + skills 软链同步）。
- 除非用户明确要求，不要绕过 hook。

## 验证循环

- 每次改动后依次跑：聚焦测试（`go test -count=1 ./<pkg>/`）→ `make lint` → 必要时全量 `go test -count=1 ./...`。
- 完成前声明通过/修复/完成必须给出新鲜验证证据（测试输出、lint 输出）。
- 新需求和 bugfix 都应补或更新测试；bugfix 必须有能复现问题的回归用例。

## 排障流程

- 按 `systematic-debugging` 流程：先复现，再定位根因，再修复。
- 用 `/docs` 的 OpenAPI 文档核对协议；用 `go run main.go server start` 本地起服务验证。
- 不把手工 `curl` 当作测试闭环——复现问题后必须补回归用例。

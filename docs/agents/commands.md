# 常用命令

> **使用场景**：执行构建、测试、运行、部署相关命令前加载。

## 本地开发

- 安装依赖：`go mod download`
- 本地运行：`go run main.go server start --host localhost --port 8080`
- 数据库迁移：`go run main.go database migrate`
- 创建对象存储桶：`go run main.go object bucket create`
- 完整本地栈：复制 `env/*.template` 为实际 env 后执行 `docker compose -f docker/docker-compose.yml up -d --build`

## 构建

- 生产构建：`make build`（strip 符号）
- 极致压缩：`make build-upx`（需安装 upx）
- 调试构建：`make build-dev` / `make build-debug`（dlv 调试用）

## 质量检查

- 静态检查：`make lint`（go vet + golangci-lint）
- 全量测试：`make test`（`go test -count=1 ./...`）
- 覆盖率：`make test-cover`
- 直接跑 golangci-lint：`golangci-lint run ./...`

## 部署

- 部署脚本：`script/deploy.sh` / `script/deploy-dev.sh`
- K8s 部署清单：`k8s/`（见 `docs/agents/repo-ci.md`）
- 只有用户明确要求时才执行部署流程。

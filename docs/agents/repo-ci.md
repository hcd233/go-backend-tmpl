# 仓库与 CI

> **使用场景**：涉及 CI/CD、镜像构建、K8s 部署参数时加载。

## CI 质量门禁

- `.github/workflows/docker-publish.yml`：`quality` job（gofmt 校验、go vet、go test、golangci-lint、hadolint）→ `build` job（linux/amd64 + arm64 矩阵、cosign 签名）→ `merge` job（manifest list）。
- 影响镜像构建的 path filter：`internal/**`、`docker/**`、`cmd/**`、`main.go`、`go.mod`、`go.sum`、`Makefile`、`.golangci.yml`。

## K8s 部署参数（k8s/deployment.yaml）

- 滚动更新：`maxUnavailable: 0` + `maxSurge: 2`（蓝绿式滚动，零中断）。
- 优雅退出：`terminationGracePeriodSeconds: 480`，预算构成：preStop 10s + inflight drain soft 5min + hard 30s + 收尾。
- `lifecycle.preStop: sleep 10`：配合 `/ready` 探针失效做流量摘除窗口。
- `GOMEMLIMIT=440MiB`（limits 512Mi 的 ~86%）：提前 GC 防贴边 OOMKill。
- 双探针分离：`livenessProbe /health`（15s/20s/3s×3）与 `readinessProbe /ready`（5s/10s/3s×6，draining 时返回 503）。
- 资源限制：`requests 50m/128Mi, limits 300m/512Mi`；`emptyDir` 日志卷带 `sizeLimit`。

## 优雅退出链路

- 信号 → `app.Stop` → fx OnStop 钩子逆序执行：Cron → Inflight → Pool → HTTP → Logger → DB → Redis。
- Inflight 排空：soft 5min 等自然完成 → 广播取消（`CancelOnDrain`）→ hard 30s 收尾；draining 期间新请求返回 503。
- 运行时指标：每 pod `Flusher` 每 5s 把 Prometheus 快照写入 Redis ZSET（`metrics:runtime:*`），留存 24h。

## 依赖与工具

- 反哺自 aris-proxy-api 的通用能力：fx DI、统一 200 错误契约、i18n、inflight、metrics、guard、filter。
- 外部技能经 `skills-lock.json` 锁定；项目自研技能放 `.agents/skills/internal/`。

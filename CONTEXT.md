# aris-api-tmpl

本文件是项目的领域词汇表（glossary），只定义术语，不含实现细节。开发前必读；新出现的领域概念、术语或语义边界必须同步补充到本文件，保持词汇表与代码一致。

## Identity & Access（身份与访问）

**User（用户）**:
一个使用平台的自然人。通过 OAuth2（GitHub / Google）注册和登录，档案包含 Name、Email、Avatar，绑定平台第三方 ID。权限体系为 pending（待审）→ user（普通）→ admin（管理员），通过 `enum.Permission` 的 `Level()` 比较等级。
_Avoid_: account, member, operator

**TokenPair（令牌对）**:
OAuth2 完成后下发的 JWT 访问令牌对，含 AccessToken（短时效，用于 API 鉴权）和 RefreshToken（长时效，用于静默续期）。签发与校验带 Issuer/Audience（`constant.ProjectName`），防 token 混淆。
_Avoid_: jwt token, session token

**OAuthProvider（OAuth 平台）**:
支持的第三方 OAuth2 登录平台，当前为 GitHub 和 Google。用户通过任一平台提供的 OAuthUserInfo（id、name、email、avatar）完成注册或登录绑定。
_Avoid_: social login, sso provider

## 通用基础设施

**Inflight（在途请求）**:
优雅退出期间对进行中请求的跟踪与两阶段排空（soft 5min 自然等待 → 广播取消 → hard 30s 收尾）。draining 期间新请求返回 503（K8s 流量摘除信号）。
_Avoid_: draining, request drain

**RuntimeMetrics（运行时指标）**:
每个 pod 周期性（5s）采集 Prometheus 快照（goroutines/heap/CPU/HTTP 时延/请求结果）写入 Redis ZSET（`metrics:runtime:*`），留存 24h，供多实例聚合查询。
_Avoid_: telemetry, monitoring data

**Guard（路由扫描防护）**:
对路由未命中（404）按 IP 记 strike，观察窗口内达到阈值即封禁（返回 403）。白名单见 `config.GuardAllowIPs`。
_Avoid_: firewall, waf

**Filter（列表筛选 DSL）**:
`internal/common/filter` 提供的字符串表达式筛选（`field:value field2:!value2`），经 `FieldConfig` 白名单映射为 SQL 条件。字段必须先在 `FieldConfig` 注册，未注册字段直接拒绝。
_Avoid_: query dsl, search syntax

---
name: huma-dto-conventions
description: Use when creating, editing, or reviewing HTTP DTOs in internal/dto/, adding new huma routes, modifying request/response payloads, or debugging "field is always zero / nil body / 422 validation" issues. Mention this skill whenever the work touches `*Req` / `*Rsp` structs, huma path/query/body bindings, or OpenAPI schema fields — even if the user only says "add an API" or "tweak the request payload".
---

# Huma DTO 规范 (huma-dto-conventions)

本模板用 [huma v2](https://huma.rocks) 把 Go 结构体直接映射成 OpenAPI schema 和 HTTP 处理器。huma 的字段绑定语义和 `encoding/json` 的"反射就反序列化"很不一样：**字段必须用 huma 认识的标签或包装结构告诉它来源于哪里（path/query/header/body/cookie）**，否则 huma 会**静默忽略**这部分输入，让字段保持零值。

本 skill 的目标是让任何人在写 / 改 / review `internal/dto/**` 时，一次就遵守 huma 规范，避免 zero-value 数据流进数据库或 redis 触发"返回错位记录""沉默丢字段"这类难排查的 bug。

## 真实事故（必看）

请求体是 `{"id": 1}`，但服务端读到的字段是 0，查询返回了主键最小的那条记录。根因在三层：

1. **DTO 层**：请求结构体把字段直接平铺在顶层（仅有 `json:"id"`），既没有 `path`/`query` 标签，也没有按 huma 约定包在 `Body *XxxReqBody` 里。huma 看不出这个字段来自哪里，于是直接跳过 body 反序列化，字段永远是零值。
2. **业务层**：没有拒绝零值字段（如 `id == 0`），把 0 当作合法值继续处理。
3. **数据访问层**：ORM 的 `First(&Model{ID: 0})` 经典坑——零值字段被当作"无 where 条件"，于是返回了主键最小的那条记录。

修复主战场是第 1 层。第 2、3 层是防御性补丁，但只要第 1 层的规范没破，后两层永远不会被触发。**这条规范是最便宜的防御，请把它写对一次就别再退化。**

## 核心规则

### 规则 1：根据字段来源选标签或包装

| 字段来源 | 怎么写 | 例 |
|---|---|---|
| URL Path 参数 | `path:"name"` 标签 | `ID uint \`path:"id" required:"true" minimum:"1"\`` |
| URL Query 参数 | `query:"name"` 标签 | `Page int \`query:"page" minimum:"1"\`` |
| HTTP Header | `header:"name"` 标签 | `Auth string \`header:"Authorization"\`` |
| Cookie | `cookie:"name"` 标签 | `SID string \`cookie:"sid"\`` |
| **JSON Body（POST / PATCH / PUT）** | **必须**在外层 Req 上声明 `Body *XxxReqBody` 字段，body 字段写在 `XxxReqBody` 里 | 见下文模板 |

**最容易踩的坑**：只写 `json:"xxx"` 而不包 Body，huma 会当成无来源的字段忽略掉 body。**`json:` 标签本身不是 huma 的"来源"标签**，它只描述序列化时的键名，不告诉 huma 该字段从哪来。

### 规则 2：JSON Body 的标准模板

所有写入 JSON body 的接口（`POST` / `PATCH` / `PUT`）一律按下面两段式写，外层 Req 只承载来源标签的字段（path / Body 包装）：

```go
// XxxReq 外层：只放来源标签字段
type XxxReq struct {
    ID   uint           `path:"id" required:"true" minimum:"1"`
    Body *XxxReqBody    `json:"body"`
}

// XxxReqBody 内层：JSON body 字段
type XxxReqBody struct {
    Name string `json:"name" required:"true" minLength:"1" maxLength:"50"`
    Age  int    `json:"age" minimum:"0"`
}
```

### 规则 3：响应统一走 CommonRsp + 顶层结构

- 成功响应：`rsp := &dto.XxxRsp{...}`，经 `util.WrapHTTPResponse(rsp, nil)` 返回，huma 序列化为 `{"data": {...}}`。
- 错误响应：**不要**手动往 `rsp.Error` 塞错误再返回 200——统一返回 `nil, apiutil.NewHumaBizError(ctx, err, ierr.ErrInternal.BizError())`（业务错误）或 `nil, apiutil.NewHumaBizErrorFromModel(ctx, ierr.ErrBadRequest.BizError())`（校验错误），huma 输出顶层 `{"error": {code, message}}`，HTTP 状态恒为 200。
- 响应体禁止出现裸 `any`/`interface{}`/`json.RawMessage` 类型（lint 强制）；需要透传时用 `dto/schema` 的 RawJSON。

### 规则 4：时间字段

DTO 时间字段用 `time.Time`；禁止在 Service 层提前格式化为字符串。字符串化（如 `time.DateTime`）只发生在 handler 组装响应时。

### 规则 5：分层依赖

`internal/dto/` 禁止导入 `internal/infrastructure/database/model`（lint 的 depguard 强制）。需要数据库字段时，将具体字段作为参数传入。

### 规则 6：路由命名

- 资源名单数小写：`/user`、`/session`。
- 禁止裸尾斜杠：`POST /endpoint` ✅，`POST /endpoint/` ❌。
- 操作通过 Path 段表达（`/user/current`），不依赖 HTTP Method 表达语义差异。
- 新路由在 `internal/router/<module>.go` 的 `initXxxRouter` 中注册，中间件在 router 层装配（`UseMiddleware` / Operation 级 `Middlewares`）。

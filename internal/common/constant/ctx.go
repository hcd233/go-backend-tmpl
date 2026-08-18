// Package constant 常量
package constant

// CtxKey context 值键的强类型，避免与内置类型 string 冲突（SA1029）。
type CtxKey string

const (
	// CtxKeyUserID undefined
	//	@update 2025-09-30 15:57:05
	CtxKeyUserID CtxKey = "userID"

	// CtxKeyUserName undefined
	//	@update 2025-09-30 15:57:07
	CtxKeyUserName CtxKey = "userName"

	// CtxKeyPermission undefined
	//	@update 2025-09-30 15:57:08
	CtxKeyPermission CtxKey = "permission"

	// CtxKeyTraceID undefined
	//	@update 2025-09-30 15:57:13
	CtxKeyTraceID CtxKey = "traceID"

	// CtxKeyLimiter undefined
	//	@update 2025-09-30 15:57:14
	CtxKeyLimiter CtxKey = "limiter"

	// CtxKeyLocale 请求语言环境（由 LocaleMiddleware 注入）。
	CtxKeyLocale CtxKey = "locale"
)

package util

import (
	"context"

	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/common/enum"
)

// CtxValueUint 从 context 中读取 uint 值。
func CtxValueUint(ctx context.Context, key constant.CtxKey) uint {
	value, ok := ctx.Value(key).(uint)
	if !ok {
		return 0
	}
	return value
}

// CtxValueString 从 context 中读取 string 值。
func CtxValueString(ctx context.Context, key constant.CtxKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}

// CtxValuePermission 从 context 中读取用户权限。
func CtxValuePermission(ctx context.Context, key constant.CtxKey) enum.Permission {
	permission, ok := ctx.Value(key).(enum.Permission)
	if !ok {
		return ""
	}
	return permission
}

// CopyContextValues 复制请求上下文的跟踪/身份字段到独立 context，
// 供异步协程池任务使用——禁止直接持有原始请求 context。
func CopyContextValues(src context.Context) (dst context.Context) {
	dst = context.Background()
	dst = context.WithValue(dst, constant.CtxKeyTraceID, src.Value(constant.CtxKeyTraceID))
	dst = context.WithValue(dst, constant.CtxKeyUserID, src.Value(constant.CtxKeyUserID))
	dst = context.WithValue(dst, constant.CtxKeyUserName, src.Value(constant.CtxKeyUserName))
	dst = context.WithValue(dst, constant.CtxKeyPermission, src.Value(constant.CtxKeyPermission))
	dst = context.WithValue(dst, constant.CtxKeyLimiter, src.Value(constant.CtxKeyLimiter))
	dst = context.WithValue(dst, constant.CtxKeyLocale, src.Value(constant.CtxKeyLocale))
	return dst
}

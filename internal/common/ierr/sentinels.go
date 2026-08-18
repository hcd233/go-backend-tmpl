package ierr

import (
	"github.com/hcd233/aris-api-tmpl/internal/common/model"
)

// ==================== 哨兵错误定义 ====================
//
// 每个哨兵错误对应一个错误类别，并映射到一个业务错误 (*model.Error)。
// 使用方通过 ierr.Wrap / ierr.New 等函数基于哨兵创建具体的 error 实例。
// Service 层可通过 ierr.ToBizError 自动从 error 中提取对应的业务错误。
//
// 命名规范：Err + 领域 + 具体错误，如 ErrDBQuery、ErrJWTDecode、ErrDTOConvert
//
// bizError code 分段：
//   - 10000-10099: 通用业务错误（与 constant/error.go 保持一致）
//   - 无需为每个内部错误分配新 code，内部错误复用已有的业务错误 code

var (
	// ==================== 通用错误 ====================

	// ErrInternal 通用内部错误（兜底）
	ErrInternal = newFromSentinel(newSentinel("internal_error", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ErrBadRequest 请求参数错误
	ErrBadRequest = newFromSentinel(newSentinel("bad_request", model.NewErrorWithKey(10006, "BadRequest", "error.bad_request")))

	// ErrUnauthorized 未授权
	ErrUnauthorized = newFromSentinel(newSentinel("unauthorized", model.NewErrorWithKey(10001, "Unauthorized", "error.unauthorized")))

	// ErrNoPermission 没有权限
	ErrNoPermission = newFromSentinel(newSentinel("no_permission", model.NewErrorWithKey(10002, "NoPermission", "error.no_permission")))

	// ErrDataNotExists 数据不存在
	ErrDataNotExists = newFromSentinel(newSentinel("data_not_exists", model.NewErrorWithKey(10003, "DataNotExists", "error.data_not_exists")))

	// ErrDataExists 数据已存在
	ErrDataExists = newFromSentinel(newSentinel("data_exists", model.NewErrorWithKey(10004, "DataExists", "error.data_exists")))

	// ErrTooManyRequests 请求过于频繁
	ErrTooManyRequests = newFromSentinel(newSentinel("too_many_requests", model.NewErrorWithKey(10005, "TooManyRequests", "error.too_many_requests")))

	// ErrInsufficientQuota 配额不足
	ErrInsufficientQuota = newFromSentinel(newSentinel("insufficient_quota", model.NewErrorWithKey(10007, "InsufficientQuota", "error.insufficient_quota")))

	// ErrResourceLocked 资源锁定
	ErrResourceLocked = newFromSentinel(newSentinel("resource_locked", model.NewErrorWithKey(10009, "ResourceLocked", "error.resource_locked")))

	// ==================== 数据库错误 ====================

	// ErrDBQuery 数据库查询错误
	ErrDBQuery = newFromSentinel(newSentinel("db_query", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ErrDBCreate 数据库创建错误
	ErrDBCreate = newFromSentinel(newSentinel("db_create", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ErrDBUpdate 数据库更新错误
	ErrDBUpdate = newFromSentinel(newSentinel("db_update", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ErrDBClose 数据库关闭错误
	ErrDBClose = newFromSentinel(newSentinel("db_close", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ==================== 认证与鉴权错误 ====================

	// ErrJWTDecode JWT 解码错误
	ErrJWTDecode = newFromSentinel(newSentinel("jwt_decode", model.NewErrorWithKey(10001, "Unauthorized", "error.unauthorized")))

	// ErrJWTEncode JWT 编码错误
	ErrJWTEncode = newFromSentinel(newSentinel("jwt_encode", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ErrOAuth2Exchange OAuth2 token 交换错误
	ErrOAuth2Exchange = newFromSentinel(newSentinel("oauth2_exchange", model.NewErrorWithKey(10001, "Unauthorized", "error.unauthorized")))

	// ErrOAuth2UserInfo OAuth2 用户信息获取错误
	ErrOAuth2UserInfo = newFromSentinel(newSentinel("oauth2_user_info", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ==================== DTO 转换错误 ====================

	// ErrDTOConvert DTO 格式转换错误
	ErrDTOConvert = newFromSentinel(newSentinel("dto_convert", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ErrDTOMarshal DTO 序列化错误
	ErrDTOMarshal = newFromSentinel(newSentinel("dto_marshal", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ErrDTOUnmarshal DTO 反序列化错误
	ErrDTOUnmarshal = newFromSentinel(newSentinel("dto_unmarshal", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ==================== SSE 流式处理错误 ====================

	// ErrSSEParse SSE 事件解析错误
	ErrSSEParse = newFromSentinel(newSentinel("sse_parse", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ErrSSEUnknownEvent SSE 未知事件类型
	ErrSSEUnknownEvent = newFromSentinel(newSentinel("sse_unknown_event", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ==================== 对象存储错误 ====================

	// ErrObjStorage 对象存储操作错误
	ErrObjStorage = newFromSentinel(newSentinel("obj_storage", model.NewErrorWithKey(10000, "InternalError", "error.internal")))

	// ==================== 校验错误 ====================

	// ErrValidation 输入校验错误
	ErrValidation = newFromSentinel(newSentinel("validation", model.NewErrorWithKey(10006, "BadRequest", "error.validation")))
)

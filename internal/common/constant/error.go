package constant

// ==================== 业务错误码 ====================
//
// 统一 200 错误契约下，业务错误码是错误语义的唯一载体：
// model.Error.StatusCode() 负责业务码 → HTTP 状态码的推导，
// apiutil.FrameworkError 负责框架错误（422/404 等）反向推导业务码。
// 新增业务错误时必须同步维护两个方向的映射。
//
// 注：业务错误实例统一由 internal/common/ierr 的哨兵错误创建，
// 本包只定义错误码常量，避免 constant ↔ model 循环依赖。

const (
	// BizErrorCodeInternal 内部错误（兜底）
	BizErrorCodeInternal = 10000

	// BizErrorCodeUnauthorized 未授权
	BizErrorCodeUnauthorized = 10001

	// BizErrorCodeNoPermission 没有权限
	BizErrorCodeNoPermission = 10002

	// BizErrorCodeDataNotExists 数据不存在
	BizErrorCodeDataNotExists = 10003

	// BizErrorCodeDataExists 数据已存在
	BizErrorCodeDataExists = 10004

	// BizErrorCodeTooManyRequests 请求过于频繁
	BizErrorCodeTooManyRequests = 10005

	// BizErrorCodeBadRequest 请求参数错误
	BizErrorCodeBadRequest = 10006

	// BizErrorCodeInsufficientQuota 配额不足
	BizErrorCodeInsufficientQuota = 10007

	// BizErrorCodeNoImplement 未实现
	BizErrorCodeNoImplement = 10008

	// BizErrorCodeResourceLocked 资源锁定
	BizErrorCodeResourceLocked = 10009

	// BizErrorDetailSep 框架错误 message 与字段错误细节之间的分隔符
	BizErrorDetailSep = ": "

	// BizErrorDetailJoinSep 多个字段错误细节之间的分隔符
	BizErrorDetailJoinSep = "; "
)

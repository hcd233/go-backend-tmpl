package constant

// Filter 解析错误消息模板。
const (
	FilterErrEmptyFieldName           = "empty field name in filter expression: %s"
	FilterErrInvalidExpr              = "invalid filter expression: %s"
	FilterErrUnknownField             = "unknown filter field: %s"
	FilterErrNullValueOp              = "operator %s not supported for NULL value"
	FilterErrUnsupportedOp            = "unsupported operator: %s"
	FilterErrInvalidRange             = "invalid range value: %s"
	FilterErrMultiValueWithComparison = "multi-value not supported with comparison operator: %s"
)

// Filter SQL 片段。
const (
	FilterSQLAND       = " AND "
	FilterSQLISNULL    = " IS NULL"
	FilterSQLISNOTNULL = " IS NOT NULL"
	FilterSQLLIKE      = " LIKE ?"
	FilterSQLNOTLIKE   = " NOT LIKE ?"
	FilterSQLEQ        = " = ?"
	FilterSQLNEQ       = " != ?"
	FilterSQLGT        = " > ?"
	FilterSQLLT        = " < ?"
	FilterSQLGTE       = " >= ?"
	FilterSQLLTE       = " <= ?"
	FilterSQLIN        = " IN (?)"
	FilterSQLNOTIN     = " NOT IN (?)"
	FilterSQLOR        = " OR "
)

package constant

// RedisDB Redis 数据库编号（避免散落魔法数字）。
const RedisDB = 0

const (
	// RuntimeMetricsInstancesKey 运行时指标-实例注册表（ZSET：member=instanceID, score=最后flush的unix秒）。
	RuntimeMetricsInstancesKey = "metrics:runtime:instances"

	// RuntimeMetricsDataKeyTemplate 运行时指标-单实例快照时序（ZSET：member=快照payload, score=快照unix秒），%s = instanceID。
	RuntimeMetricsDataKeyTemplate = "metrics:runtime:data:%s"

	// RedisZRangePositiveInfinity ZSET 范围查询的上界。
	RedisZRangePositiveInfinity = "+inf"
)

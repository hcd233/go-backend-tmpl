package lock

import (
	"context"
	"time"

	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/cache"
	"github.com/redis/go-redis/v9"
)

// Locker 锁接口
//
//	@param ctx context.Context
//	@param key string
//	@param value string
//	@return err error
//	@author centonhuang
//	@update 2025-11-11 16:54:41
type Locker interface {
	Lock(ctx context.Context, key string, value string, expire time.Duration) (success bool, err error)
	// Refresh 仅持有者可续期锁，返回是否续期成功（value 不匹配时返回 false）。
	Refresh(ctx context.Context, key string, value string, expire time.Duration) (success bool, err error)
	Unlock(ctx context.Context, key string, value string) (err error)
}

// NewLocker 创建锁
//
//	@return Locker
//	@author centonhuang
//	@update 2025-11-11 17:49:18
func NewLocker() Locker {
	return &redisLocker{
		redis: cache.GetRedisClient(),
	}
}

type redisLocker struct {
	redis *redis.Client
}

func (l *redisLocker) Lock(ctx context.Context, key, value string, expire time.Duration) (success bool, err error) {
	return l.redis.SetNX(ctx, key, value, expire).Result()
}

func (l *redisLocker) Refresh(ctx context.Context, key, value string, expire time.Duration) (success bool, err error) {
	res, err := l.redis.Eval(ctx, constant.LuaRefreshLock, []string{key}, value, expire.Milliseconds()).Int64()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func (l *redisLocker) Unlock(ctx context.Context, key, value string) (err error) {
	return l.redis.Eval(ctx, constant.LuaUnlockLock, []string{key}, value).Err()
}

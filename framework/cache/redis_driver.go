package cache

import (
	"context"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/redis"
	"time"
)

type RedisDriver struct {
	redis  redis.Client
	prefix string
}

var _ contract.Cache = (*RedisDriver)(nil)

func NewRedisDriver(redis redis.Client, prefix string) *RedisDriver {
	return &RedisDriver{
		redis:  redis,
		prefix: prefix,
	}
}

func (r *RedisDriver) Get(ctx context.Context, key string) ([]byte, error) {
	return r.redis.Get(ctx, r.prefix+key).Bytes()
}

func (r *RedisDriver) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.redis.Set(ctx, r.prefix+key, value, ttl).Err()
}

func (r *RedisDriver) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	for i, v := range keys {
		keys[i] = r.prefix + v
	}
	return r.redis.Del(ctx, keys...).Err()
}

func (r *RedisDriver) Has(ctx context.Context, key string) bool {
	return r.redis.Exists(ctx, r.prefix+key).Val() == 1
}

func (r *RedisDriver) ClearPrefix(ctx context.Context, prefix string) error {
	key := r.prefix + prefix
	var cursor uint64
	var n int
	var err error
	for {
		var keys []string
		keys, cursor, err = r.redis.Scan(ctx, cursor, key, 10000).Result()
		if err != nil {
			return err
		}
		n += len(keys)
		r.redis.Del(ctx, keys...)

		if cursor == 0 {
			break
		}
	}
	return err
}

func (r *RedisDriver) Clear(ctx context.Context) error {
	return r.redis.FlushDB(ctx).Err()
}

package cache

import (
	"context"
	"github.com/allegro/bigcache"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"strings"
	"time"
)

type BigcacheDriver struct {
	bigcache *bigcache.BigCache
	prefix   string
}

var _ contract.Cache = (*BigcacheDriver)(nil)

func NewBigcacheDriver(config bigcache.Config, prefix string) *BigcacheDriver {
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		panic(err)
	}

	b := &BigcacheDriver{
		cache,
		prefix,
	}
	return b
}

func (b *BigcacheDriver) Get(_ context.Context, key string) ([]byte, error) {
	return b.bigcache.Get(b.prefix + key)
}

func (b *BigcacheDriver) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	return b.bigcache.Set(b.prefix+key, value)
}

func (b *BigcacheDriver) Delete(_ context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		_ = b.bigcache.Delete(b.prefix + key)
	}
	return nil
}

func (b *BigcacheDriver) Has(_ context.Context, key string) bool {
	if _, err := b.bigcache.Get(b.prefix + key); err != nil {
		return false
	}
	return true
}

func (b *BigcacheDriver) ClearPrefix(_ context.Context, prefix string) error {
	iterator := b.bigcache.Iterator()
	prefix = b.prefix + prefix
	for iterator.SetNext() {
		current, err := iterator.Value()
		if err == nil {
			key := current.Key()
			if strings.HasPrefix(key, prefix) {
				_ = b.bigcache.Delete(key)
			}
		}
	}
	return nil
}

func (b *BigcacheDriver) Clear(_ context.Context) error {
	return b.bigcache.Reset()
}

package contract

import (
	"context"
	"time"
)

type Cache interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, []byte, time.Duration) error
	Delete(context.Context, ...string) error
	Has(context.Context, string) bool
	ClearPrefix(context.Context, string) error
	Clear(context.Context) error
}

type CacheManager interface {
	Cache
	Driver(...string) Cache
	Register(string, Cache)
}

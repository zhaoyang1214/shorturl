package contract

import (
	"github.com/go-redis/redis/v8"
)

type Redis interface {
	redis.UniversalClient
	Connection(name ...string) redis.UniversalClient
	Resolve(name string) redis.UniversalClient
}

package redis

import (
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/redis"
)

type Redis struct {
}

var _ contract.Provider = (*Redis)(nil)

func (r *Redis) Build(container contract.Container, params ...interface{}) (interface{}, error) {
	appServer, err := container.Get("app")
	if err != nil {
		return nil, err
	}

	return redis.NewRedis(appServer.(contract.Application)), nil
}

package cache

import (
	"github.com/zhaoyang1214/ginco/framework/cache"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

type Cache struct {
}

var _ contract.Provider = (*Cache)(nil)

func (c *Cache) Build(container contract.Container, params ...interface{}) (interface{}, error) {
	appServer, err := container.Get("app")
	if err != nil {
		return nil, err
	}

	return cache.NewCache(appServer.(contract.Application)), nil
}

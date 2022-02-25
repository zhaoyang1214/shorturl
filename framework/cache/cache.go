package cache

import (
	"github.com/allegro/bigcache"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"time"
)

type Cache struct {
	contract.Cache
	app     contract.Application
	drivers map[string]contract.Cache
}

var _ contract.CacheManager = (*Cache)(nil)

func NewCache(app contract.Application) *Cache {
	c := &Cache{
		app:     app,
		drivers: make(map[string]contract.Cache),
	}
	c.Cache = c.Driver()
	return c
}

func (c *Cache) Driver(driver ...string) contract.Cache {
	var name string
	if len(driver) == 0 {
		name = c.app.GetI("config").(contract.Config).GetString("cache.default")
	} else {
		name = driver[0]
	}
	if cache, ok := c.drivers[name]; ok {
		return cache
	}
	c.drivers[name] = c.resolve(name)
	return c.drivers[name]
}

func (c *Cache) Register(driver string, cache contract.Cache) {
	c.drivers[driver] = cache
}

func (c *Cache) resolve(name string) contract.Cache {
	cacheConf := c.app.GetI("config").(contract.Config).Sub("cache")
	storeConf := cacheConf.Sub("stores." + name)
	if storeConf == nil {
		panic("Cache config [" + name + "] is not defined")
	}

	var prefix string
	if prefix = storeConf.GetString("prefix"); prefix == "" {
		prefix = cacheConf.GetString("prefix")
	}
	driver := storeConf.GetString("driver")
	switch driver {
	case "redis":
		return c.createRedisDriver(storeConf.GetString("connection"), prefix)
	case "database":
		return c.createDatabaseDriver(storeConf.GetString("connection"), storeConf.GetString("table"), prefix)
	case "memory":
		return c.createBigcacheDriver(prefix, storeConf)
	}
	panic("Cache driver [" + driver + "] is not supported")
}

func (c *Cache) createRedisDriver(connection, prefix string) contract.Cache {
	return NewRedisDriver(c.app.GetI("redis").(contract.Redis).Connection(connection), prefix)
}

func (c *Cache) createDatabaseDriver(connection, table, prefix string) contract.Cache {
	var conn []string
	if connection != "" {
		conn = append(conn, connection)
	}
	return NewDatabaseDriver(c.app.GetI("db").(contract.Database).Connection(conn...), table, prefix)
}

func (c *Cache) createBigcacheDriver(prefix string, conf contract.Config) contract.Cache {
	lifeWindow := conf.GetDuration("life_window")
	if lifeWindow <= 0 {
		lifeWindow = 10
	}
	bigcacheConf := bigcache.DefaultConfig(lifeWindow * time.Minute)

	if shards := conf.GetInt("shards"); shards > 0 {
		bigcacheConf.Shards = shards
	}

	if conf.Has("clean_window") {
		cleanWindow := conf.GetDuration("clean_window")
		if cleanWindow > 0 {
			cleanWindow = cleanWindow * time.Minute
		}
		bigcacheConf.CleanWindow = cleanWindow
	}

	if maxEntriesInWindow := conf.GetInt("max_entries_in_window"); maxEntriesInWindow > 0 {
		bigcacheConf.MaxEntriesInWindow = maxEntriesInWindow
	}

	if maxEntrySize := conf.GetInt("max_entry_size"); maxEntrySize > 0 {
		bigcacheConf.MaxEntrySize = maxEntrySize
	}

	if conf.Has("verbose") {
		bigcacheConf.Verbose = conf.GetBool("verbose")
	}

	if conf.Has("hard_max_cache_size") {
		bigcacheConf.HardMaxCacheSize = conf.GetInt("hard_max_cache_size")
	}

	return NewBigcacheDriver(bigcacheConf, prefix)
}

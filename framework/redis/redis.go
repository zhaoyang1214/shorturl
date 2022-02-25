package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"time"
)

type Client = redis.UniversalClient

type Redis struct {
	Client
	app     contract.Application
	clients map[string]Client
}

var _ contract.Redis = (*Redis)(nil)

func NewRedis(app contract.Application) *Redis {
	r := &Redis{
		app:     app,
		clients: make(map[string]Client),
	}
	r.Client = r.Connection()
	return r
}

func (r *Redis) Connection(names ...string) Client {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}
	if c, ok := r.clients[name]; ok {
		return c
	}
	r.clients[name] = r.Resolve(name)
	return r.clients[name]
}

func (r *Redis) Resolve(name string) Client {
	conf := r.app.GetI("config").(contract.Config).Sub("redis." + name)
	if conf == nil {
		panic("Redis config [" + name + "] is not defined")
	}

	addrs := conf.GetStringSlice("addrs")
	db := conf.GetInt("db")
	password := conf.GetString("password")

	options := &redis.UniversalOptions{
		Addrs:    addrs,
		DB:       db,
		Password: password,
	}

	if sentinelPassword := conf.GetString("sentinel_password"); sentinelPassword != "" {
		options.SentinelPassword = sentinelPassword
	}

	if username := conf.GetString("username"); username != "" {
		options.Username = username
	}

	if conf.Has("max_retries") {
		options.MaxRetries = conf.GetInt("max_retries")
	}

	if conf.Has("min_retry_backoff") {
		options.MinRetryBackoff = conf.GetDuration("min_retry_backoff") * time.Millisecond
	}

	if conf.Has("max_retry_backoff") {
		options.MaxRetryBackoff = conf.GetDuration("max_retry_backoff") * time.Millisecond
	}

	if conf.Has("dial_timeout") {
		options.DialTimeout = conf.GetDuration("dial_timeout") * time.Second
	}

	if conf.Has("read_timeout") {
		options.ReadTimeout = conf.GetDuration("read_timeout") * time.Second
	}

	if conf.Has("write_timeout") {
		options.WriteTimeout = conf.GetDuration("write_timeout") * time.Second
	}

	if conf.Has("pool_fifo") {
		options.PoolFIFO = conf.GetBool("pool_fifo")
	}

	if conf.Has("pool_size") {
		options.PoolSize = conf.GetInt("pool_size")
	}

	if conf.Has("min_idle_conns") {
		options.MinIdleConns = conf.GetInt("min_idle_conns")
	}

	if conf.Has("max_conn_age") {
		options.MaxConnAge = conf.GetDuration("max_conn_age") * time.Millisecond
	}

	if conf.Has("pool_timeout") {
		options.PoolTimeout = conf.GetDuration("pool_timeout") * time.Second
	}

	if conf.Has("idle_timeout") {
		idleTimeout := conf.GetDuration("idle_timeout")
		if idleTimeout != -1 {
			idleTimeout *= time.Minute
		}
		options.IdleTimeout = idleTimeout
	}

	if conf.Has("idle_check_frequency") {
		idleCheckFrequency := conf.GetDuration("idle_check_frequency")
		if idleCheckFrequency != -1 {
			idleCheckFrequency *= time.Minute
		}
		options.IdleCheckFrequency = idleCheckFrequency
	}

	if conf.Has("max_redirects") {
		options.MaxRedirects = conf.GetInt("max_redirects")
	}

	if conf.Has("read_only") {
		options.ReadOnly = conf.GetBool("read_only")
	}

	if conf.Has("route_by_latency") {
		options.RouteByLatency = conf.GetBool("route_by_latency")
	}

	if conf.Has("route_randomly") {
		options.RouteRandomly = conf.GetBool("route_randomly")
	}

	if conf.Has("master_name") {
		options.MasterName = conf.GetString("master_name")
	}

	return redis.NewUniversalClient(options)
}

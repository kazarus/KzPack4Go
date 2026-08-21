package KzPack4Go

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type TRedis4Go struct {
	RedisTyp int64
	Client1  *redis.Client
	Client2  *redis.ClusterClient
}

func NewClient(RedisTyp int64, RedisSrv string, RedisUsr string, RedisPwd string) *TRedis4Go {

	var Redis4Go = TRedis4Go{}
	Redis4Go.RedisTyp = RedisTyp

	switch RedisTyp {
	case 0, 1:
		{
			options := redis.Options{
				Network:         "tcp",
				Addr:            fmt.Sprintf("%s", RedisSrv),
				Dialer:          nil,
				OnConnect:       nil,
				Password:        fmt.Sprintf("%s", RedisPwd),
				DB:              11,
				MaxRetries:      0,
				MinRetryBackoff: 0,
				MaxRetryBackoff: 0,
				DialTimeout:     0,
				ReadTimeout:     1000 * time.Second,
				WriteTimeout:    1000 * time.Second,
				PoolSize:        0,
				MinIdleConns:    0,
				MaxIdleConns:    0,
				PoolTimeout:     0,
				TLSConfig:       nil,
				Limiter:         nil,
			}

			Redis4Go.Client1 = redis.NewClient(&options)
		}
	case 2:
		{
			var ListAddr = strings.Split(RedisSrv, ",")
			options := &redis.ClusterOptions{
				Addrs:    ListAddr,
				Password: RedisPwd,
			}

			Redis4Go.Client2 = redis.NewClusterClient(options)
		}
	}

	return &Redis4Go
}

func (this TRedis4Go) Ping(ctx context.Context) *redis.StatusCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.Ping(ctx)
		}
	case 2:
		{
			return this.Client2.Ping(ctx)
		}
	}

	return nil
}

func (this TRedis4Go) Close() error {

	switch this.RedisTyp {
	case 0, 1:
		{
			if this.Client1 != nil {
				return this.Client1.Close()
			}
		}
	case 2:
		{
			if this.Client2 != nil {
				return this.Client2.Close()
			}
		}
	}

	return nil
}

func (this TRedis4Go) Get(ctx context.Context, key string) *redis.StringCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.Get(ctx, key)
		}
	case 2:
		{
			return this.Client2.Get(ctx, key)
		}
	}

	return nil
}

func (this TRedis4Go) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.Set(ctx, key, value, expiration)
		}
	case 2:
		{
			return this.Client2.Set(ctx, key, value, expiration)
		}
	}

	return nil
}

func (this TRedis4Go) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.Keys(ctx, pattern)
		}
	case 2:
		{
			return this.Client2.Keys(ctx, pattern)
		}
	}

	return nil
}

func (this TRedis4Go) Exists(ctx context.Context, keys ...string) *redis.IntCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.Exists(ctx, keys...)
		}
	case 2:
		{
			return this.Client2.Exists(ctx, keys...)
		}
	}

	return nil
}

func (this TRedis4Go) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.HSet(ctx, key, values...)
		}
	case 2:
		{
			return this.Client2.HSet(ctx, key, values...)
		}
	}

	return nil
}

func (this TRedis4Go) HGet(ctx context.Context, key, field string) *redis.StringCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.HGet(ctx, key, field)
		}
	case 2:
		{
			return this.Client2.HGet(ctx, key, field)
		}
	}

	return nil
}

func (this TRedis4Go) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.HScan(ctx, key, cursor, match, count)
		}
	case 2:
		{
			return this.Client2.HScan(ctx, key, cursor, match, count)
		}
	}

	return nil
}

func (this TRedis4Go) Del(ctx context.Context, keys ...string) *redis.IntCmd {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.Del(ctx, keys...)
		}
	case 2:
		{
			return this.Client2.Del(ctx, keys...)
		}
	}

	return nil
}

func (this TRedis4Go) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {

	switch this.RedisTyp {
	case 0, 1:
		{
			return this.Client1.PSubscribe(ctx, channels...)
		}
	case 2:
		{
			return this.Client2.PSubscribe(ctx, channels...)
		}
	}

	return nil
}

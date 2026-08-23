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

// TRedisOptions Redis 客户端可选参数;零值字段使用新默认值
// (单机 DB=0,读写超时 3s)。如需旧行为(DB=11,读写 1000s)请显式传参
// 或继续使用 NewClient。
type TRedisOptions struct {
	DB           int           // #单机库号;零值即 DB 0
	ReadTimeout  time.Duration // #读超时;<=0 时使用默认 3s
	WriteTimeout time.Duration // #写超时;<=0 时使用默认 3s
}

// NewClient 保留旧默认值(单机 DB=11,读写超时 1000s)以兼容历史调用点。
//
// Deprecated: 旧默认值反直觉,新代码请使用 NewClientWithOptions
// (零值字段使用新默认:DB 0、读写超时 3s)。
func NewClient(RedisTyp int64, RedisSrv string, RedisUsr string, RedisPwd string) *TRedis4Go {

	return NewClientWithOptions(RedisTyp, RedisSrv, RedisUsr, RedisPwd,
		TRedisOptions{DB: 11, ReadTimeout: 1000 * time.Second, WriteTimeout: 1000 * time.Second})
}

// NewClientWithOptions 构造 Redis 客户端,允许自定义单机库号与读写超时;
// 零值字段使用新默认值(DB 0,读写超时 3s)。
func NewClientWithOptions(RedisTyp int64, RedisSrv string, RedisUsr string, RedisPwd string, Options TRedisOptions) *TRedis4Go {

	var Redis4Go = TRedis4Go{}
	Redis4Go.RedisTyp = RedisTyp

	// #零值回退到新默认值
	db := Options.DB
	if db < 0 {
		db = 0
	}
	readTimeout := Options.ReadTimeout
	if readTimeout <= 0 {
		readTimeout = 3 * time.Second
	}
	writeTimeout := Options.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = 3 * time.Second
	}

	switch RedisTyp {
	case 0, 1:
		{
			options := redis.Options{
				Network:      "tcp",
				Addr:         RedisSrv,
				Username:     RedisUsr,
				Password:     RedisPwd,
				DB:           db,
				ReadTimeout:  readTimeout,
				WriteTimeout: writeTimeout,
			}

			Redis4Go.Client1 = redis.NewClient(&options)
		}
	case 2:
		{
			var ListAddr = strings.Split(RedisSrv, ",")
			options := &redis.ClusterOptions{
				Addrs:    ListAddr,
				Username: RedisUsr,
				Password: RedisPwd,
			}

			Redis4Go.Client2 = redis.NewClusterClient(options)
		}
	default:
		{
			//#非法类型立即失败,避免后续命令方法返回nil引发远程panic
			panic(fmt.Sprintf("Redis4Go: 非法的RedisTyp=%d; 有效值{0,1:单机;2:集群}", RedisTyp))
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

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

// TRedisOptions Redis 客户端可选参数;零值沿用 NewClient 的旧默认值
// (单机 DB=11,读写超时 1000s)。如确需使用库 0,可在构造后通过
// Client1.Options().DB 直接调整。
type TRedisOptions struct {
	DB           int           // #单机库号;<=0 时使用默认 11
	ReadTimeout  time.Duration // #读超时;<=0 时使用默认 1000s
	WriteTimeout time.Duration // #写超时;<=0 时使用默认 1000s
}

func NewClient(RedisTyp int64, RedisSrv string, RedisUsr string, RedisPwd string) *TRedis4Go {

	return NewClientWithOptions(RedisTyp, RedisSrv, RedisUsr, RedisPwd, TRedisOptions{})
}

// NewClientWithOptions 与 NewClient 相同,但允许自定义单机库号与读写超时;
// 零值字段回退到旧默认值,保证与 NewClient 行为一致。
func NewClientWithOptions(RedisTyp int64, RedisSrv string, RedisUsr string, RedisPwd string, Options TRedisOptions) *TRedis4Go {

	var Redis4Go = TRedis4Go{}
	Redis4Go.RedisTyp = RedisTyp

	// #零值回退到旧默认值(与 NewClient 一致)
	db := Options.DB
	if db <= 0 {
		db = 11
	}
	readTimeout := Options.ReadTimeout
	if readTimeout <= 0 {
		readTimeout = 1000 * time.Second
	}
	writeTimeout := Options.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = 1000 * time.Second
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

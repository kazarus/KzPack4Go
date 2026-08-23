package KzPack4Go

import (
	"context"
	"net"
	"os"
	"testing"
	"time"
)

// ---------- 单元测试:无需真实 Redis 服务 ----------

// RedisTyp 0/1 应路由到单机客户端,且不初始化集群客户端
func TestNewClientSingle(t *testing.T) {
	for _, typ := range []int64{0, 1} {
		r := NewClient(typ, "127.0.0.1:6379", "", "")
		if r.RedisTyp != typ {
			t.Errorf("RedisTyp = %d; want %d", r.RedisTyp, typ)
		}
		if r.Client1 == nil {
			t.Errorf("RedisTyp=%d: 单机模式 Client1 未初始化", typ)
		}
		if r.Client2 != nil {
			t.Errorf("RedisTyp=%d: 单机模式不应初始化 Client2", typ)
		}
	}
}

// RedisTyp 2 应路由到集群客户端,并正确解析逗号分隔的地址列表
func TestNewClientCluster(t *testing.T) {
	addrs := "10.0.0.1:6379,10.0.0.2:6379,10.0.0.3:6379"
	r := NewClient(2, addrs, "", "")
	if r.Client2 == nil {
		t.Fatal("集群模式 Client2 未初始化")
	}
	if r.Client1 != nil {
		t.Error("集群模式不应初始化 Client1")
	}

	got := r.Client2.Options().Addrs
	if len(got) != 3 {
		t.Fatalf("集群地址数 = %d; want 3", len(got))
	}
	if got[0] != "10.0.0.1:6379" || got[1] != "10.0.0.2:6379" || got[2] != "10.0.0.3:6379" {
		t.Errorf("集群地址解析错误: %v", got)
	}
}

// 未知 RedisTyp:修复后应在构造期立即 panic(快速失败),而非返回nil客户端
func TestNewClientUnknownType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewClient 对未知 RedisTyp 应立即 panic(快速失败)")
		}
	}()
	NewClient(9, "127.0.0.1:6379", "", "")
}

// NewClientWithOptions 应透传自定义库号与读写超时;零值字段回退旧默认值(与 NewClient 一致)
func TestNewClientWithOptions(t *testing.T) {
	r := NewClientWithOptions(0, "127.0.0.1:6379", "", "", TRedisOptions{DB: 5, ReadTimeout: 2 * time.Second, WriteTimeout: 3 * time.Second})
	if r.Client1 == nil {
		t.Fatal("Client1 未初始化")
	}
	opts := r.Client1.Options()
	if opts.DB != 5 {
		t.Errorf("DB = %d; want 5", opts.DB)
	}
	if opts.ReadTimeout != 2*time.Second {
		t.Errorf("ReadTimeout = %v; want 2s", opts.ReadTimeout)
	}
	if opts.WriteTimeout != 3*time.Second {
		t.Errorf("WriteTimeout = %v; want 3s", opts.WriteTimeout)
	}

	// 零值回退旧默认(单机 DB=11,读写超时 1000s)
	r2 := NewClientWithOptions(0, "127.0.0.1:6379", "", "", TRedisOptions{})
	opts2 := r2.Client1.Options()
	if opts2.DB != 11 {
		t.Errorf("零值 DB = %d; want 11", opts2.DB)
	}
	if opts2.ReadTimeout != 1000*time.Second {
		t.Errorf("零值 ReadTimeout = %v; want 1000s", opts2.ReadTimeout)
	}
	if opts2.WriteTimeout != 1000*time.Second {
		t.Errorf("零值 WriteTimeout = %v; want 1000s", opts2.WriteTimeout)
	}
}

// Close 对未建立过连接的客户端应安全返回 nil
func TestCloseWithoutConnect(t *testing.T) {
	r := NewClient(0, "127.0.0.1:16379", "", "")
	if err := r.Close(); err != nil {
		t.Errorf("Close() = %v; want nil", err)
	}
}

// 未知 RedisTyp 下 Close 应返回 nil 而非 panic(直接构造结构体,绕过 NewClient 校验)
func TestCloseUnknownType(t *testing.T) {
	r := &TRedis4Go{RedisTyp: 9}
	if err := r.Close(); err != nil {
		t.Errorf("Close() = %v; want nil", err)
	}
}

// 未知 RedisTyp 下命令方法返回 nil 命令对象——记录行为边界(nil 上调用 Result 会 panic;
// 经 NewClient 构造的实例已在构造期拦截,仅直接构造结构体可触达此分支)
func TestCommandsReturnNilOnUnknownType(t *testing.T) {
	r := &TRedis4Go{RedisTyp: 9}
	ctx := context.Background()

	if cmd := r.Get(ctx, "k"); cmd != nil {
		t.Errorf("Get() = %v; want nil", cmd)
	}
	if cmd := r.Set(ctx, "k", "v", 0); cmd != nil {
		t.Errorf("Set() = %v; want nil", cmd)
	}
	if cmd := r.Del(ctx, "k"); cmd != nil {
		t.Errorf("Del() = %v; want nil", cmd)
	}
}

// HScan:单机模式返回非 nil 命令对象(无需真实连接);未知类型返回 nil
func TestHScan(t *testing.T) {
	r := NewClient(0, "127.0.0.1:6379", "", "")
	ctx := context.Background()

	if cmd := r.HScan(ctx, "hkey", 0, "*", 10); cmd == nil {
		t.Error("单机模式 HScan 返回 nil")
	}
	if cmd := (&TRedis4Go{RedisTyp: 9}).HScan(ctx, "hkey", 0, "*", 10); cmd != nil {
		t.Error("未知类型 HScan 应返回 nil")
	}
}

// 集群模式命令方法返回非 nil 命令对象。
// 注意:go-redis v9 命令方法会立即执行(拨号/路由),集群地址不可达时必须用
// 短超时 ctx 让命令快速失败,此处仅验证命令对象构造与路由逻辑非 nil。
func TestClusterCommandsNonNil(t *testing.T) {
	r := NewClient(2, "10.0.0.1:6379,10.0.0.2:6379,10.0.0.3:6379", "", "")
	defer r.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	if r.Ping(ctx) == nil || r.Get(ctx, "k") == nil || r.Set(ctx, "k", "v", 0) == nil ||
		r.Keys(ctx, "*") == nil || r.Exists(ctx, "k") == nil || r.HSet(ctx, "h", "f", "v") == nil ||
		r.HGet(ctx, "h", "f") == nil || r.HScan(ctx, "h", 0, "*", 10) == nil || r.Del(ctx, "k") == nil {
		t.Error("集群模式命令方法返回 nil")
	}
}

// ---------- 集成测试:需要本机 Redis,不可用时自动跳过 ----------

// redisAddr 支持通过环境变量 REDIS_ADDR 覆盖目标地址
func redisAddr() string {
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:6379"
}

func skipIfNoRedis(t *testing.T) string {
	t.Helper()
	addr := redisAddr()
	conn, eror := net.DialTimeout("tcp", addr, time.Second)
	if eror != nil {
		t.Skipf("Redis(%s) 不可用,跳过集成测试: %v", addr, eror)
	}
	conn.Close()
	return addr
}

// 集成测试:Ping / Set / Get / Exists / HSet / HGet / Keys / PSubscribe / Del 回环
// 注意:单机模式固定写入 DB 11(见 Redis4Go.go NewClient)
func TestRedisIntegration(t *testing.T) {
	addr := skipIfNoRedis(t)

	r := NewClient(0, addr, "", "")
	defer r.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//#连通性
	if eror := r.Ping(ctx).Err(); eror != nil {
		t.Skipf("Redis(%s) 无响应: %v", addr, eror)
	}

	//#Set/Get 回环
	key := "kzpack4go:test:string"
	if eror := r.Set(ctx, key, "hello", time.Minute).Err(); eror != nil {
		t.Fatalf("Set: %v", eror)
	}
	if got := r.Get(ctx, key).Val(); got != "hello" {
		t.Errorf("Get = %q; want hello", got)
	}

	//#不存在的键应返回 redis.Nil
	if eror := r.Get(ctx, "kzpack4go:test:missing").Err(); eror == nil {
		t.Error("Get 缺失键应返回错误(redis.Nil)")
	}

	//#Exists
	if n := r.Exists(ctx, key).Val(); n != 1 {
		t.Errorf("Exists = %d; want 1", n)
	}

	//#HSet/HGet 回环
	hkey := "kzpack4go:test:hash"
	if eror := r.HSet(ctx, hkey, "field", "hvalue").Err(); eror != nil {
		t.Fatalf("HSet: %v", eror)
	}
	if got := r.HGet(ctx, hkey, "field").Val(); got != "hvalue" {
		t.Errorf("HGet = %q; want hvalue", got)
	}

	//#Keys 模式匹配
	keys := r.Keys(ctx, "kzpack4go:test:*").Val()
	if len(keys) < 2 {
		t.Errorf("Keys 匹配数 = %d; want >= 2", len(keys))
	}

	//#PSubscribe 订阅与关闭
	ps := r.PSubscribe(ctx, "kzpack4go:test:*")
	if ps == nil {
		t.Fatal("PSubscribe 返回 nil")
	}
	if eror := ps.Close(); eror != nil {
		t.Errorf("PubSub.Close: %v", eror)
	}

	//#清理
	if n := r.Del(ctx, key, hkey).Val(); n != 2 {
		t.Errorf("Del = %d; want 2", n)
	}
}

package KzPack4Go

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// ---------- TTL 与容量上限 ----------

func TestNewBlackList(t *testing.T) {
	bl := NewBlackList(100, 5*time.Minute)
	if bl.MaxSize != 100 || bl.TTL != 5*time.Minute {
		t.Errorf("NewBlackList 字段 = %+v", bl)
	}
}

// 过期后 Select 视为不存在;惰性清除(条目保留到下次 Append 整表清扫)
func TestBlackListTTL(t *testing.T) {
	fake := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	bl := &TBlackList{TTL: time.Minute, now: func() time.Time { return fake }}

	bl.Append("u1", TErrorUsr{UserCode: "u1", ErrorCnt: 1})
	if usr, ok := bl.Select("u1"); !ok || usr.ErrorCnt != 1 {
		t.Fatalf("未过期应命中: %+v, %v", usr, ok)
	}

	fake = fake.Add(61 * time.Second) // 过期
	if _, ok := bl.Select("u1"); ok {
		t.Error("过期后 Select 应返回不存在")
	}
	if len(bl.BlackList) != 1 {
		t.Error("Select 不应直接删除条目(惰性清除)")
	}
}

// 只设 TTL 时,追加会按 TTL/2 间隔整表清扫过期条目,防止驻留
func TestBlackListTTLPurgeOnAppend(t *testing.T) {
	fake := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	bl := &TBlackList{TTL: time.Minute, now: func() time.Time { return fake }}

	bl.Append("u1", TErrorUsr{UserCode: "u1"})
	bl.Append("u2", TErrorUsr{UserCode: "u2"})

	fake = fake.Add(61 * time.Second)          // 全部过期
	bl.Append("u3", TErrorUsr{UserCode: "u3"}) // 触发整表清扫

	if len(bl.BlackList) != 1 {
		t.Errorf("清扫后残留 %d 条; want 1", len(bl.BlackList))
	}
	if _, ok := bl.Select("u3"); !ok {
		t.Error("新条目 u3 应存在")
	}
}

// 容量超限(TTL=0)按 FIFO 淘汰最早插入的条目
func TestBlackListCapacityFIFO(t *testing.T) {
	fake := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	bl := &TBlackList{MaxSize: 3, now: func() time.Time { return fake }}

	bl.Append("a", TErrorUsr{UserCode: "a"})
	bl.Append("b", TErrorUsr{UserCode: "b"})
	bl.Append("c", TErrorUsr{UserCode: "c"})
	bl.Append("d", TErrorUsr{UserCode: "d"}) // 触发淘汰

	if len(bl.BlackList) != 3 {
		t.Fatalf("容量超限后大小 = %d; want 3", len(bl.BlackList))
	}
	if _, ok := bl.Select("a"); ok {
		t.Error("最早插入的 a 应被淘汰")
	}
	for _, k := range []string{"b", "c", "d"} {
		if _, ok := bl.Select(k); !ok {
			t.Errorf("%s 应保留", k)
		}
	}
}

// 容量超限 + TTL:在 TTL/2 清扫窗口内,淘汰"最早到期"的条目
func TestBlackListCapacityWithTTL(t *testing.T) {
	fake := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	bl := &TBlackList{MaxSize: 3, TTL: time.Hour, now: func() time.Time { return fake }}

	bl.Append("a", TErrorUsr{UserCode: "a"}) // 13:00 到期
	bl.Append("b", TErrorUsr{UserCode: "b"}) // 13:00 到期
	bl.Append("c", TErrorUsr{UserCode: "c"}) // 13:00 到期

	fake = fake.Add(10 * time.Minute) // 12:10,未到清扫阈值(30min),无过期条目
	bl.Append("d", TErrorUsr{UserCode: "d"})

	if len(bl.BlackList) != 3 {
		t.Fatalf("容量超限后大小 = %d; want 3", len(bl.BlackList))
	}
	if _, ok := bl.Select("a"); ok {
		t.Error("最早到期的 a 应被淘汰")
	}
	for _, k := range []string{"b", "c", "d"} {
		if _, ok := bl.Select(k); !ok {
			t.Errorf("%s 应保留", k)
		}
	}
}

// 更新已有条目不触发淘汰,并刷新到期时间
func TestBlackListUpdateRefreshes(t *testing.T) {
	fake := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	bl := &TBlackList{MaxSize: 2, TTL: time.Minute, now: func() time.Time { return fake }}

	bl.Append("a", TErrorUsr{UserCode: "a", ErrorCnt: 1})
	bl.Append("b", TErrorUsr{UserCode: "b", ErrorCnt: 1})

	// 更新 a 不淘汰 b
	bl.Append("a", TErrorUsr{UserCode: "a", ErrorCnt: 2})
	if len(bl.BlackList) != 2 {
		t.Fatalf("更新不应触发淘汰: size = %d", len(bl.BlackList))
	}
	if usr, ok := bl.Select("a"); !ok || usr.ErrorCnt != 2 {
		t.Errorf("更新后 a = %+v, %v", usr, ok)
	}

	// 30s 后刷新 a:a 的到期时间顺延为 12:01:30,b 仍为 12:01
	fake = fake.Add(30 * time.Second)
	bl.Append("a", TErrorUsr{UserCode: "a", ErrorCnt: 3})
	if _, ok := bl.Select("a"); !ok {
		t.Error("a 未到期应存在")
	}

	// 12:01:10:a(12:01:30 到期)仍有效;b(12:01 到期)已过期
	fake = fake.Add(40 * time.Second)
	if _, ok := bl.Select("a"); !ok {
		t.Error("a 在 12:01:10 应仍有效(到期 12:01:30)")
	}
	if _, ok := bl.Select("b"); ok {
		t.Error("b 未刷新,已过期应不存在")
	}
}

// ---------- 基础行为(零值) ----------

func TestBlackListBasic(t *testing.T) {
	bl := TBlackList{}

	bl.Append("u001", TErrorUsr{UserCode: "u001", ErrorCnt: 1})
	bl.Append("u002", TErrorUsr{UserCode: "u002", ErrorCnt: 2})

	got, ok := bl.Select("u001")
	if !ok || got.ErrorCnt != 1 {
		t.Errorf("Select(u001) = %+v, %v", got, ok)
	}

	bl.Remove("u001")
	if _, ok := bl.Select("u001"); ok {
		t.Error("u001 should be removed")
	}
	// 不存在的 key 返回零值
	if _, ok := bl.Select("nobody"); ok {
		t.Error("Select(nobody) should be false")
	}
}

// 并发读写,配合 go test -race 验证无数据竞争
func TestBlackListConcurrent(t *testing.T) {
	bl := TBlackList{}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("u%03d", n)
			bl.Append(key, TErrorUsr{UserCode: key, ErrorCnt: int64(n)})
		}(i)
	}
	wg.Wait()

	var wg2 sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg2.Add(1)
		go func(n int) {
			defer wg2.Done()
			key := fmt.Sprintf("u%03d", n)
			_, _ = bl.Select(key)
		}(i)
	}
	wg2.Wait()
}

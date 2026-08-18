package KzPack4Go

import (
	"fmt"
	"sync"
	"testing"
)

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

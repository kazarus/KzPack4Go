package KzPack4Go

import (
	"strings"
	"testing"
)

func TestStartEd(t *testing.T) {
	dbg := TKzDebug{}
	dbg.StartEd(false)

	if dbg.RunDebug {
		t.Error("StartEd(false) 后 RunDebug 应为 false")
	}
	if dbg.ThisTick != 0 || dbg.LastTick != 0 || dbg.TotlTick != 0 {
		t.Errorf("StartEd 应清零计时字段: %+v", dbg)
	}
}

func TestTickFmt(t *testing.T) {
	dbg := TKzDebug{}
	dbg.StartEd(false) // 不打印,纯计时

	first := dbg.TickFmt("start")
	if !strings.Contains(first, "start") {
		t.Errorf("TickFmt 首条消息缺少标签: %s", first)
	}
	if dbg.ThisTick == 0 || dbg.LastTick == 0 || dbg.TotlTick == 0 {
		t.Fatalf("TickFmt 后计时字段应为非零: %+v", dbg)
	}

	second := dbg.TickFmt("step2")
	if !strings.Contains(second, "step2") {
		t.Errorf("TickFmt 第二条缺少标签: %s", second)
	}
	if dbg.LastTick != dbg.ThisTick {
		t.Error("TickFmt 后 LastTick 应更新为 ThisTick")
	}
}

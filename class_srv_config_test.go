package KzPack4Go

import (
	"strings"
	"testing"
)

func TestSrvConfigToDefault(t *testing.T) {
	cfg := TSrvConfig{}
	cfg.ToDefault()

	if cfg.UploadAt != "upload" || cfg.Log4goAt != "log4go" {
		t.Errorf("默认目录异常: upload=%q log4go=%q", cfg.UploadAt, cfg.Log4goAt)
	}
	if cfg.WorkerID != 1 || cfg.SynzFrom != "NTOW" || cfg.LoginMin != 9 {
		t.Errorf("默认业务字段异常: %+v", cfg)
	}
	if cfg.MaxOpenConns != 500 || cfg.MaxIdleConns != 200 || cfg.MaxLifetimes != 60 {
		t.Errorf("默认连接池异常: %+v", cfg)
	}
	if cfg.Path4ExtExcludes != ".html,.js,.ts,.exe" || cfg.AccessControlAllowOrigin != "*" {
		t.Errorf("默认安全字段异常: %+v", cfg)
	}
	if cfg.MakeYearWhenTodo != 1 || cfg.VoidYearWhenTodo != 0 {
		t.Errorf("默认年度字段异常: %+v", cfg)
	}
	// 旧字段迁移
	if cfg.LockYear != cfg.LockKJND || cfg.RealYear != cfg.RealKJND {
		t.Error("LockYear/RealYear 应迁移自 LockKJND/RealKJND")
	}
}

func TestSrvConfigToUpFormat(t *testing.T) {
	cfg := TSrvConfig{UploadAt: "up"}

	// 默认(0):up/{who}/
	if got := cfg.ToUpFormat(7); got != "up/7/" {
		t.Errorf("UpFormat=0 = %s", got)
	}
	// 1:up/{who}/{年度}01/
	cfg.UpFormat = 1
	got := cfg.ToUpFormat(7)
	if !strings.HasPrefix(got, "up/7/") || !strings.HasSuffix(got, "01/") {
		t.Errorf("UpFormat=1 = %s", got)
	}
	// 2:up/{who}/{年月}/
	cfg.UpFormat = 2
	got = cfg.ToUpFormat(7)
	if !strings.HasPrefix(got, "up/7/") || !strings.HasSuffix(got, "/") {
		t.Errorf("UpFormat=2 = %s", got)
	}
	if !strings.HasSuffix(got, "/"+IntToStr(ToKJQJ())+"/") {
		t.Errorf("UpFormat=2 应包含当前年月: %s", got)
	}
	// 3:up/{who}/{年周}/
	cfg.UpFormat = 3
	got = cfg.ToUpFormat(7)
	if !strings.HasSuffix(got, "/"+IntToStr(ToWeek())+"/") {
		t.Errorf("UpFormat=3 应包含当前周: %s", got)
	}
}

func TestSrvConfigToRealPath(t *testing.T) {
	cfg := TSrvConfig{}
	if got := cfg.ToRealPath(7); got != "7/" {
		t.Errorf("ToRealPath 默认 = %s", got)
	}
	cfg.UpFormat = 2
	if got := cfg.ToRealPath(7); !strings.HasPrefix(got, "7/") {
		t.Errorf("ToRealPath UpFormat=2 = %s", got)
	}
}

func TestSrvConfigClearAttr(t *testing.T) {
	cfg := TSrvConfig{
		EndPoint: "ep", CloudURL: "cloud", PrintURL: "print", FrontEnd: "3000",
		HttpPort: "8080", ListenAt: "8080", Ding4APP: "a", Ding4Dom: "d",
		Ding4GOV: "g", Ding4URL: "u", Path4ExtExcludes: ".exe", RealYear: 2026,
		RedisKey: "rk", RedisPwd: "pwd", RedisSrv: "srv", SocketAt: "sock",
		SilentAt: "silent", UploadAt: "up", Log4goAt: "log",
	}
	cfg.ClearAttr()

	if cfg.EndPoint != "" || cfg.CloudURL != "" || cfg.PrintURL != "" || cfg.FrontEnd != "" ||
		cfg.HttpPort != "" || cfg.ListenAt != "" || cfg.Ding4APP != "" || cfg.Ding4Dom != "" ||
		cfg.Ding4GOV != "" || cfg.Ding4URL != "" || cfg.Path4ExtExcludes != "" || cfg.RealYear != 0 ||
		cfg.RedisKey != "" || cfg.RedisPwd != "" || cfg.RedisSrv != "" || cfg.SocketAt != "" ||
		cfg.SilentAt != "" || cfg.UploadAt != "" || cfg.Log4goAt != "" {
		t.Errorf("ClearAttr 后敏感/内部字段未清空: %+v", cfg)
	}
}

func TestSrvConfigToRabbitMQ(t *testing.T) {
	cfg := TSrvConfig{HttpPort: "8080"}
	if got := cfg.ToRabbitMQ("order"); got != "tech.cxcloud.8080.order" {
		t.Errorf("ToRabbitMQ = %s", got)
	}
}

func TestSrvConfigGetMakeYearWhenTodo(t *testing.T) {
	cfg := TSrvConfig{MakeYearWhenTodo: 1}
	if got := cfg.GetMakeYearWhenTodo(2026); got != 2026 {
		t.Errorf("MakeYearWhenTodo=1 应返回源年度: %d", got)
	}
	cfg.MakeYearWhenTodo = 0
	if got := cfg.GetMakeYearWhenTodo(2026); got != -1 {
		t.Errorf("MakeYearWhenTodo=0 应返回 -1: %d", got)
	}
}

func TestSrvConfigGetVoidYearWhenTodo(t *testing.T) {
	cfg := TSrvConfig{VoidYearWhenTodo: 1}
	if got := cfg.GetVoidYearWhenTodo(2026); got != 2026 {
		t.Errorf("VoidYearWhenTodo=1 应返回源年度: %d", got)
	}
	cfg.VoidYearWhenTodo = 0
	if got := cfg.GetVoidYearWhenTodo(2026); got != -1 {
		t.Errorf("VoidYearWhenTodo=0 应返回 -1: %d", got)
	}
}

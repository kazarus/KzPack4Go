package KzPack4Go

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSysStatusInitialize(t *testing.T) {
	st := TSYSSTATUS{}
	st.Initialize()

	if st.DataStat != 1 || st.KJNDKJQJ != -1 {
		t.Errorf("Initialize 默认状态异常: %+v", st)
	}
	if st.DataGUID == "" || st.EntityID == "" {
		t.Error("Initialize 应生成 GUID/EntityID")
	}
	if st.DataDate != ToDate() || st.DataTime != ToTime() {
		t.Error("Initialize 时间戳异常")
	}
}

func TestSysStatusInTokenUsr(t *testing.T) {
	st := TSYSSTATUS{}
	usr := TTokenUsr{
		UserName: "张三",
		EnrollID: "E001",
		EntityID: "U001",
		DataFrom: CONST_DATA_FROM_TO_钉钉,
	}
	st.InTokenUsr(usr)

	if st.MastName != "张三" || st.MastGUID != "U001" || st.EnrollID != "E001" {
		t.Errorf("InTokenUsr 主标识异常: %+v", st)
	}
	if st.CreateBy != "U001" || st.UpdateBy != "U001" || st.CreateAt == 0 || st.UpdateAt == 0 {
		t.Errorf("InTokenUsr 审计字段异常: %+v", st)
	}
	if st.DataFrom != CONST_DATA_FROM_TO_钉钉 {
		t.Error("InTokenUsr 终端来源异常")
	}
}

func TestSysStatusToMarkDown(t *testing.T) {
	st := TSYSSTATUS{Name4App: "app", HttpPort: "8080", FrontEnd: "3000", DataBase: "db"}

	md0 := st.ToMarkDown0("标题", "一切正常")
	if !strings.Contains(md0, "标题") || !strings.Contains(md0, "一切正常") {
		t.Errorf("ToMarkDown0 内容异常: %s", md0)
	}
	if !strings.Contains(md0, "app") || !strings.Contains(md0, "8080") || !strings.Contains(md0, "db") {
		t.Errorf("ToMarkDown0 缺少服务信息: %s", md0)
	}
	// 校验 JSON 结构
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(md0), &m); err != nil {
		t.Fatalf("ToMarkDown0 非法 JSON: %v", err)
	}
	if m["msgtype"] != "markdown" {
		t.Errorf("ToMarkDown0 msgtype = %v", m["msgtype"])
	}

	if !strings.Contains(st.ToMarkDown1(), "服务重启") {
		t.Error("ToMarkDown1 应为服务重启提醒")
	}
	if !strings.Contains(st.ToMarkDown2(), "服务暂停") {
		t.Error("ToMarkDown2 应为服务暂停提醒")
	}
}

func TestSysStatusToFeedCard(t *testing.T) {
	st := TSYSSTATUS{
		EntityID: "{12345678-1234-1234-1234-123456789012}",
		Name4App: "app", Path4URL: "http://x/log", Url4Path: "/api/x",
		ClientIP: "10.0.0.1", EndPoint: "ep", HttpPort: "8080", FrontEnd: "3000",
		MastName: "张三", EnrollID: "E001", DataBase: "db", ErrorMsg: "err",
		FromText: "桌面终端",
	}

	card := st.ToFeedCard1()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(card), &m); err != nil {
		t.Fatalf("ToFeedCard1 非法 JSON: %v", err)
	}
	if m["msgtype"] != "feedCard" {
		t.Errorf("ToFeedCard1 msgtype = %v", m["msgtype"])
	}
	fc, ok := m["feedCard"].(map[string]interface{})
	if !ok {
		t.Fatalf("ToFeedCard1 缺 feedCard: %v", m)
	}
	links, ok := fc["links"].([]interface{})
	if !ok || len(links) < 3 {
		t.Errorf("ToFeedCard1 links 数量异常: %v", fc["links"])
	}

	// ToFeedCard2 应输出完整实体 JSON
	card2 := st.ToFeedCard2()
	if !strings.Contains(card2, "\"whoBuild\"") || !strings.Contains(card2, "\"app\"") {
		t.Errorf("ToFeedCard2 内容异常: %s", card2)
	}
}

func TestSysStatusToErrorLog(t *testing.T) {
	// 日志写入可执行文件所在目录,测试后清理
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	exeDir, _ := filepath.Split(exe)
	defer os.RemoveAll(filepath.Join(exeDir, "kzlog"))

	st := TSYSSTATUS{}
	st.Initialize()
	st.DataFrom = CONST_DATA_FROM_TO_桌面
	st.RqSource = "GET /api HTTP/1.1"

	usr := TTokenUsr{UserName: "张三", EnrollID: "E001", EntityID: "U001"}
	cfg := TSrvConfig{Log4goAt: "kzlog", CloudURL: "http://cloud"}

	if err := st.ToErrorLog(usr, cfg, TBotConfig{}, nil); err != nil {
		t.Fatal(err)
	}

	if st.DiskPath == "" || st.Path4Log == "" || st.Path4Req == "" || st.Path4URL == "" {
		t.Errorf("日志路径未生成: %+v", st)
	}
	if !IsFile(st.DiskPath) {
		t.Error("日志文件未落盘")
	}
	if !IsFile(st.Path4Req) {
		t.Error("请求文件未落盘")
	}
	if st.FromText != "桌面终端" {
		t.Errorf("FromText = %s", st.FromText)
	}
	if !strings.HasPrefix(st.Path4URL, "http://cloud/kzlog/") {
		t.Errorf("Path4URL = %s", st.Path4URL)
	}
}

package KzPack4Go

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Ding4Bot=1:真实钉钉机器人路径,校验签名参数与请求体
func TestToErrorMsg1DingTalk(t *testing.T) {
	var gotRawQuery string
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRawQuery = r.URL.RawQuery
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	bot := TBotConfig{Ding4Bot: 1, EndPoint: srv.URL, AccessID: "tok", SecretID: "sec"}
	bot.ToErrorMsg1(`{"msgtype":"text","text":{"content":"hi"}}`)

	if !strings.Contains(gotRawQuery, "access_token=tok") {
		t.Errorf("缺少 access_token: %s", gotRawQuery)
	}
	if !strings.Contains(gotRawQuery, "timestamp=") || !strings.Contains(gotRawQuery, "sign=") {
		t.Errorf("缺少 timestamp/sign: %s", gotRawQuery)
	}
	if !strings.Contains(gotBody, "msgtype") {
		t.Errorf("请求体异常: %s", gotBody)
	}
}

// Ding4Bot=0:内部网关路径,POST form 到 /core/tsysstatus/syncfrom
func TestToErrorMsg1Gateway(t *testing.T) {
	var gotPath string
	var gotForm string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.PostFormValue("params")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	bot := TBotConfig{Ding4Bot: 0, EndPoint: srv.URL}
	bot.ToErrorMsg1("some error")

	if gotPath != "/core/tsysstatus/syncfrom" {
		t.Errorf("请求路径 = %s", gotPath)
	}
	if gotForm != "some error" {
		t.Errorf("form params = %s", gotForm)
	}
}

// 空消息/空网关直接返回,不发请求
func TestToErrorMsg1Empty(t *testing.T) {
	bot := TBotConfig{Ding4Bot: 0, EndPoint: "http://127.0.0.1:1"}
	bot.ToErrorMsg1("") // 空消息
	bot2 := TBotConfig{Ding4Bot: 0, EndPoint: ""}
	bot2.ToErrorMsg1("err") // 空网关
}

// ToErrorMsg2:JSON 消息 → TSYSSTATUS → 推送 feedCard JSON(跳过附件上传)
func TestToErrorMsg2Gateway(t *testing.T) {
	var gotForm string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotForm = r.PostFormValue("params")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	st := TSYSSTATUS{EntityID: "e1", Name4App: "app", Url4Path: "/x"}
	data, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}

	bot := TBotConfig{Ding4Bot: 0, EndPoint: srv.URL, UploadD2: -1}
	bot.ToErrorMsg2(string(data))

	if !strings.Contains(gotForm, "\"entityID\":\"e1\"") {
		t.Errorf("推送内容异常: %s", gotForm)
	}
}

// ToErrorMsg2:附件文件不存在时跳过上传,仍完成推送
func TestToErrorMsg2GatewayUploadMissing(t *testing.T) {
	var gotForm string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotForm = r.PostFormValue("params")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	st := TSYSSTATUS{EntityID: "e1", Path4Req: "/nonexistent/req.txt"}
	data, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}

	bot := TBotConfig{Ding4Bot: 0, EndPoint: srv.URL, UploadD2: 0}
	bot.ToErrorMsg2(string(data))

	if !strings.Contains(gotForm, "\"entityID\":\"e1\"") {
		t.Errorf("附件缺失时也应完成推送: %s", gotForm)
	}
}

// ToErrorMsg2 钉钉路径:签名参数正确,消息体直传
func TestToErrorMsg2DingTalk(t *testing.T) {
	var gotRawQuery string
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRawQuery = r.URL.RawQuery
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	bot := TBotConfig{Ding4Bot: 1, EndPoint: srv.URL, AccessID: "tok", SecretID: "sec"}
	bot.ToErrorMsg2(`{"entityID":"e1"}`)

	if !strings.Contains(gotRawQuery, "access_token=tok") {
		t.Errorf("缺少 access_token: %s", gotRawQuery)
	}
	if gotBody == "" {
		t.Error("请求体为空")
	}
}

// 非法 JSON 消息直接返回,不发请求
func TestToErrorMsg2BadJSON(t *testing.T) {
	bot := TBotConfig{Ding4Bot: 1, EndPoint: "http://127.0.0.1:1"}
	bot.ToErrorMsg2("not-json") // 不应 panic
}

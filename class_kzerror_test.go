package KzPack4Go

import (
	"encoding/json"
	"testing"
)

func TestToErrorCnt(t *testing.T) {
	// ErrorCnt 参数必须被赋值
	result := ToErrorCnt("密码错误", "PWD", 3)
	if result.ErrorCnt != 3 {
		t.Errorf("ToErrorCnt.ErrorCnt = %d; want 3", result.ErrorCnt)
	}
	if result.ErrorKey != "PWD" {
		t.Errorf("ToErrorCnt.ErrorKey = %s; want PWD", result.ErrorKey)
	}
	if result.OnStatus {
		t.Error("ToErrorCnt.OnStatus = true; want false")
	}
}

func TestToLite(t *testing.T) {
	if got := ToLite(""); got != "" {
		t.Errorf("ToLite(\"\") = %s", got)
	}
	// 短字符串不应 panic,原样返回
	if got := ToLite("abc"); got != "abc" {
		t.Errorf("ToLite(abc) = %s", got)
	}
	if got := ToLite("{12345678-1234-1234-1234-123456789012}"); got != "1234...9012" {
		t.Errorf("ToLite(uuid) = %s", got)
	}
}

func TestToNext19(t *testing.T) {
	a := ToNext19()
	b := ToNext19()
	if a == 0 || b == 0 {
		t.Fatalf("ToNext19 returned zero: %d, %d", a, b)
	}
	if a == b {
		t.Errorf("ToNext19 returned duplicate: %d", a)
	}
}

func TestToUUID(t *testing.T) {
	a := ToUUID()
	b := ToUUID()
	if a == b {
		t.Errorf("ToUUID returned duplicate: %s", a)
	}
	if len(a) != 38 { // {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
		t.Errorf("ToUUID length = %d; want 38", len(a))
	}
}

func TestTERORMarshal(t *testing.T) {
	result := ToMEMO("ok")
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["ONSTATUS"] != true {
		t.Errorf("ONSTATUS = %v; want true", m["ONSTATUS"])
	}
	if m["ERORCODE"] != "TRUE" {
		t.Errorf("ERORCODE = %v; want TRUE", m["ERORCODE"])
	}
	// 空 ListData 应输出为数组而非 null
	if m["LISTDATA"] == nil {
		t.Error("LISTDATA is nil; want []")
	}
}

func TestAtRoot(t *testing.T) {
	if got := AtRoot(); got == "" {
		t.Error("AtRoot() returned empty")
	}
}

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

// ---------- TEROR 构造器 ----------

func TestToEROR(t *testing.T) {
	r := ToEROR("出错了")
	if r.OnStatus || r.ErorCode != "EROR" || r.ErorMemo != "出错了" {
		t.Errorf("ToEROR = %+v", r)
	}
	if r.DataDate != ToDate() || r.DataTime != ToTime() {
		t.Error("ToEROR 时间戳异常")
	}
}

func TestToTRUE(t *testing.T) {
	r := ToTRUE([]int{1, 2})
	if !r.OnStatus || r.ErorCode != "TRUE" {
		t.Errorf("ToTRUE = %+v", r)
	}
	if list, ok := r.ListData.([]int); !ok || len(list) != 2 {
		t.Errorf("ToTRUE ListData = %v", r.ListData)
	}
}

func TestToSIZE(t *testing.T) {
	r := ToSIZE([]string{"a"}, 1)
	if !r.OnStatus || r.DataSize != 1 {
		t.Errorf("ToSIZE = %+v", r)
	}
}

func TestToPage(t *testing.T) {
	r := ToPage([]int{1}, 1, 2, 10, 25)
	if !r.OnStatus || r.RealPage != 2 || r.PageSize != 10 || r.TotalCnt != 25 || r.DataSize != 1 {
		t.Errorf("ToPage = %+v", r)
	}
}

func TestToDATA(t *testing.T) {
	r := ToDATA(true, "OK", "memo", 3, "data")
	if !r.OnStatus || r.ErorCode != "OK" || r.ErorMemo != "memo" || r.DataSize != 3 || r.ListData != "data" {
		t.Errorf("ToDATA = %+v", r)
	}
}

func TestToErrorKey(t *testing.T) {
	r := ToErrorKey("密码错误", "PWD")
	if r.OnStatus || r.ErrorKey != "PWD" || r.ErrorCnt != 0 {
		t.Errorf("ToErrorKey = %+v", r)
	}
}

func TestToErrLevel(t *testing.T) {
	r := ToErrLevel("严重", 2, 3, "KEY")
	if r.ErrLevel != 2 || r.ErrorCnt != 3 || r.ErrorKey != "KEY" {
		t.Errorf("ToErrLevel = %+v", r)
	}
}

func TestToCryptTxt(t *testing.T) {
	r := ToCryptTxt(true, "OK", "memo", 1, "cipher", nil)
	if !r.OnStatus || r.CryptTxt != "cipher" || r.DataSize != 1 {
		t.Errorf("ToCryptTxt = %+v", r)
	}
}

func TestTErrorMsgMarshal(t *testing.T) {
	em := TErrorMsg{RowIndex: 1, ErrorKey: "K", ErrorMsg: "m", ErrorCnt: 2}
	data, err := json.Marshal(em)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["rowIndex"] != float64(1) || m["errorKey"] != "K" || m["errorMsg"] != "m" || m["errorCnt"] != float64(2) {
		t.Errorf("TErrorMsg JSON = %v", m)
	}
}

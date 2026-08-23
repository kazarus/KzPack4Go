package KzPack4Go

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToHash(t *testing.T) {
	// md5("abc") = 900150983cd24fb0d6963f7d28e17f72
	got, err := ToHash("abc", "md5", false)
	if err != nil || got != "900150983cd24fb0d6963f7d28e17f72" {
		t.Errorf("ToHash(md5) = %s, %v", got, err)
	}
	// sha1("abc") = a9993e364706816aba3e25717850c26c9cd0d89d
	got, err = ToHash("abc", "sha1", false)
	if err != nil || got != "a9993e364706816aba3e25717850c26c9cd0d89d" {
		t.Errorf("ToHash(sha1) = %s, %v", got, err)
	}
	// sha256("abc") = ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad
	got, err = ToHash("abc", "sha256", false)
	if err != nil || got != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Errorf("ToHash(sha256) = %s, %v", got, err)
	}
	// 不支持的类型返回错误而非字符串(修复前返回错误提示字符串)
	if got, err := ToHash("abc", "md4", false); err == nil || got != "" {
		t.Errorf("ToHash(md4) = %q, %v; want 空串+错误", got, err)
	}
	if got, err := ToHash("abc", "ripemd160", false); err == nil || got != "" {
		t.Errorf("ToHash(ripemd160) = %q, %v; want 空串+错误", got, err)
	}
	// isHex=true 时输入非法 hex 应返回错误
	if got, err := ToHash("zzzz", "md5", true); err == nil || got != "" {
		t.Errorf("ToHash(非法hex) = %q, %v; want 空串+错误", got, err)
	}
}

func TestFileExist(t *testing.T) {
	dir := t.TempDir()
	exist := filepath.Join(dir, "a.txt")
	missing := filepath.Join(dir, "b.txt")

	if err := os.WriteFile(exist, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	ok, err := FileExist(exist)
	if !ok || err != nil {
		t.Errorf("FileExist(exist) = %v, %v; want true, nil", ok, err)
	}

	// 文件不存在时应返回 (false, nil),而非 (false, 非nil错误)
	ok, err = FileExist(missing)
	if ok || err != nil {
		t.Errorf("FileExist(missing) = %v, %v; want false, nil", ok, err)
	}
}

func TestSubStr(t *testing.T) {
	if got := SubStr("KzPack4Go", 2, 6); got != "Pack" {
		t.Errorf("SubStr = %s", got)
	}
	// 越界不应 panic
	if got := SubStr("KzPack4Go", -1, 3); got != "" {
		t.Errorf("SubStr(-1) = %s", got)
	}
	if got := SubStr("KzPack4Go", 0, 99); got != "" {
		t.Errorf("SubStr(99) = %s", got)
	}
	// ended < start 不应 panic(修复前 slice bounds out of range)
	if got := SubStr("KzPack4Go", 6, 2); got != "" {
		t.Errorf("SubStr(6, 2) = %s", got)
	}
}

func TestNextKJQJ(t *testing.T) {
	if got := NextKJQJ(202512); got != 202601 {
		t.Errorf("NextKJQJ(202512) = %d", got)
	}
	if got := NextKJQJ(202509); got != 202510 {
		t.Errorf("NextKJQJ(202509) = %d", got)
	}
}

func TestFloat0Text(t *testing.T) {
	if got := Float0Text(1.5000); got != "1.5" {
		t.Errorf("Float0Text(1.5) = %s", got)
	}
	if got := Float0Text(0); got != "0" {
		t.Errorf("Float0Text(0) = %s", got)
	}
	if got := Float0Text(2.0); got != "2" {
		t.Errorf("Float0Text(2.0) = %s", got)
	}
}

func TestTimeToStr(t *testing.T) {
	if got := TimeToStr(20200618105004); got != "2020-06-18 10:50:04" {
		t.Errorf("TimeToStr = %s", got)
	}
	if got := TimeToStr(20200618); got != "-1" {
		t.Errorf("TimeToStr(short) = %s", got)
	}
}

func TestDateToStr(t *testing.T) {
	if got := DateToStr(20200618); got != "2020-06-18" {
		t.Errorf("DateToStr = %s", got)
	}
	if got := DateToStr(2020061); got != "-1" {
		t.Errorf("DateToStr(short) = %s", got)
	}
}

func TestContains(t *testing.T) {
	if !Contains("a,b,c", "b") {
		t.Error("Contains(a,b,c, b) = false")
	}
	if Contains("a,b,c", "d") {
		t.Error("Contains(a,b,c, d) = true")
	}
}

func TestFormatFloat(t *testing.T) {
	if got := FormatFloat(1.234567, 2); got != 1.23 {
		t.Errorf("FormatFloat(1.234567, 2) = %v", got)
	}
	if got := FormatFloat(1.234567, -1); got != 1 {
		t.Errorf("FormatFloat(1.234567, -1) = %v", got)
	}
}

// NewIdGenerator 的 WorkerID 参数必须生效(修复前被无条件覆盖为3)
// 默认位布局:ID = 时间戳<<12 | WorkerId<<6 | 序列号,故 worker 位段 = (id>>6)&0x3F
func TestNewIdGenerator(t *testing.T) {
	NewIdGenerator(5)

	id := ToNext()
	if id == 0 {
		t.Fatal("ToNext() = 0")
	}
	if got := (id >> 6) & 0x3F; got != 5 {
		t.Errorf("生成ID中的WorkerId = %d; want 5(参数应生效)", got)
	}
}

// ---------- 其余工具函数 ----------

func TestMinMax(t *testing.T) {
	if got := Min(1, 2); got != 1 {
		t.Errorf("Min(1,2) = %d", got)
	}
	if got := Max(1, 2); got != 2 {
		t.Errorf("Max(1,2) = %d", got)
	}
}

func TestTimeStampFuncs(t *testing.T) {
	// ToMsec 紧凑数字格式:20260822150405.000 级别
	if got := ToMsec(); got < 20000000000000 || got > 21000101000000 {
		t.Errorf("ToMsec 越界: %v", got)
	}
	if got := ToHour(); got < 0 || got > 23 {
		t.Errorf("ToHour 越界: %d", got)
	}
	if got := ToWeek(); got < 200001 || got%100 < 1 || got%100 > 53 {
		t.Errorf("ToWeek 越界: %d", got)
	}
	if got := ToKJND(); got < 2000 || got > 2100 {
		t.Errorf("ToKJND 越界: %d", got)
	}
	if got := ToKJQJ(); got/100 != ToKJND() || got%100 < 1 || got%100 > 12 {
		t.Errorf("ToKJQJ 异常: %d", got)
	}
	// TimePlus/DatePlus 与当期一致,+1 天严格增大
	if TimePlus(0) != ToTime() {
		t.Error("TimePlus(0) 应等于 ToTime()")
	}
	if TimePlus(1) <= TimePlus(0) {
		t.Error("TimePlus(1) 应大于 TimePlus(0)")
	}
	if DatePlus(0) != ToDate() {
		t.Error("DatePlus(0) 应等于 ToDate()")
	}
	if DatePlus(1) <= DatePlus(0) {
		t.Error("DatePlus(1) 应大于 DatePlus(0)")
	}
}

func TestToCron(t *testing.T) {
	got := ToCron(60)
	parts := strings.Split(got, " ")
	if len(parts) != 6 || parts[5] != "?" {
		t.Errorf("ToCron 格式异常: %q", got)
	}
}

func TestFormat1Time2Time(t *testing.T) {
	if got := Format1Time("20200618"); got != "2020-06-18 00:00:00" {
		t.Errorf("Format1Time = %s", got)
	}
	if got := Format2Time("20200618"); got != "2020-06-18 23:59:59" {
		t.Errorf("Format2Time = %s", got)
	}
	if got := Format1Time("2020"); got != "-1" {
		t.Errorf("Format1Time(非法) = %s", got)
	}
	if got := Format2Time("2020"); got != "-1" {
		t.Errorf("Format2Time(非法) = %s", got)
	}
}

func TestStrConvert(t *testing.T) {
	if got := StrToInt("123"); got != 123 {
		t.Errorf("StrToInt = %d", got)
	}
	if got := IntToStr(123); got != "123" {
		t.Errorf("IntToStr = %s", got)
	}
	if got := StrToFloat("1.5"); got != 1.5 {
		t.Errorf("StrToFloat = %v", got)
	}
}

func TestFloatFormat(t *testing.T) {
	if got := Float2Text(1.2345); got != "1.2345" {
		t.Errorf("Float2Text = %s", got)
	}
	if got := Format2Float(3.14159); got != 3.14 {
		t.Errorf("Format2Float = %v", got)
	}
	if got := Format4Float(3.14159265); got != 3.1416 {
		t.Errorf("Format4Float = %v", got)
	}
	if got := Format6Float(1.23456789); got != 1.234568 {
		t.Errorf("Format6Float = %v", got)
	}
}

func TestFormatFileSize(t *testing.T) {
	cases := []struct {
		size int64
		want string
	}{
		{512, "512.00B"},
		{2048, "2.00KB"},
		{5 * 1024 * 1024, "5.00MB"},
		{3 * 1024 * 1024 * 1024, "3.00GB"},
		{2 * 1024 * 1024 * 1024 * 1024, "2.00TB"},
	}
	for _, c := range cases {
		if got := FormatFileSize(c.size); got != c.want {
			t.Errorf("FormatFileSize(%d) = %s; want %s", c.size, got, c.want)
		}
	}
}

func TestIfThen(t *testing.T) {
	if IfThenI(true, 1, 0) != 1 || IfThenI(false, 1, 0) != 0 {
		t.Error("IfThenI 异常")
	}
	if IfThenS(true, "y", "n") != "y" || IfThenS(false, "y", "n") != "n" {
		t.Error("IfThenS 异常")
	}
	if IfThenF(true, 1.5, 0.5) != 1.5 || IfThenF(false, 1.5, 0.5) != 0.5 {
		t.Error("IfThenF 异常")
	}
	if IfThenB(true, true, false) != true || IfThenB(false, true, false) != false {
		t.Error("IfThenB 异常")
	}
}

func TestContain(t *testing.T) {
	if !Contain("b", []string{"a", "b", "c"}) {
		t.Error("Contain(b) = false")
	}
	if Contain("z", []string{"a", "b"}) {
		t.Error("Contain(z) = true")
	}
}

func TestIsFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if !IsFile(f) {
		t.Error("IsFile(文件) = false")
	}
	if IsFile(dir) {
		t.Error("IsFile(目录) = true")
	}
	if IsFile(filepath.Join(dir, "missing")) {
		t.Error("IsFile(不存在) = true")
	}
}

func TestToFileInFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")

	cfg := map[string]interface{}{"a": 1}
	if err := ToFile(path, cfg); err != nil {
		t.Fatal(err)
	}
	if !IsFile(path) {
		t.Fatal("ToFile 未生成文件")
	}

	// InFile 读取已有文件
	var got map[string]interface{}
	if err := InFile(path, &got); err != nil {
		t.Fatal(err)
	}
	if got["a"] != float64(1) {
		t.Errorf("InFile 读取 = %v", got)
	}

	// InFile 文件不存在时创建
	missing := filepath.Join(dir, "new.json")
	var v map[string]interface{}
	if err := InFile(missing, &v); err != nil {
		t.Fatal(err)
	}
	if !IsFile(missing) {
		t.Error("InFile 未创建缺失文件")
	}
}

func TestAnyToVal(t *testing.T) {
	if got := AnyToVal(nil); got != "" {
		t.Errorf("AnyToVal(nil) = %q", got)
	}
	if got := AnyToVal(42); got != "42" {
		t.Errorf("AnyToVal(int) = %q", got)
	}
	if got := AnyToVal(int64(42)); got != "42" {
		t.Errorf("AnyToVal(int64) = %q", got)
	}
	if got := AnyToVal(1.5); got != "1.5" {
		t.Errorf("AnyToVal(float) = %q", got)
	}
	if got := AnyToVal("str"); got != "str" {
		t.Errorf("AnyToVal(string) = %q", got)
	}
	if got := AnyToVal([]byte("hi")); got != "hi" {
		t.Errorf("AnyToVal([]byte) = %q", got)
	}
	if got := AnyToVal(true); got != "true" {
		t.Errorf("AnyToVal(bool) = %q", got)
	}
	if got := AnyToVal(struct{ X int }{X: 1}); got != `{"X":1}` {
		t.Errorf("AnyToVal(struct) = %q", got)
	}
}

func TestToHashHexInput(t *testing.T) {
	// isHex=true:输入为 hex 编码的明文,"616263" = "abc"
	got, err := ToHash("616263", "md5", true)
	if err != nil || got != "900150983cd24fb0d6963f7d28e17f72" {
		t.Errorf("ToHash(hex) = %s, %v", got, err)
	}
}

func TestFormatFileSizeEB(t *testing.T) {
	if got := FormatFileSize(2 * 1024 * 1024 * 1024 * 1024 * 1024); got != "2.00EB" {
		t.Errorf("FormatFileSize(2EB) = %s", got)
	}
}

func TestAnyToValMoreTypes(t *testing.T) {
	for _, v := range []interface{}{uint(42), int8(42), uint8(42), int16(42), uint16(42), int32(42), uint32(42), uint64(42)} {
		if got := AnyToVal(v); got != "42" {
			t.Errorf("AnyToVal(%T) = %q", v, got)
		}
	}
	if got := AnyToVal(float32(1.5)); got != "1.5" {
		t.Errorf("AnyToVal(float32) = %q", got)
	}
}

func TestKzUnicodeMarshal(t *testing.T) {
	data, err := json.Marshal(KzUnicode("中文"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"\u4e2d\u6587"` {
		t.Errorf("KzUnicode JSON = %s", data)
	}
}

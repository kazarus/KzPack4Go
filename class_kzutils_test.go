package KzPack4Go

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToHash(t *testing.T) {
	// md5("abc") = 900150983cd24fb0d6963f7d28e17f72
	if got := ToHash("abc", "md5", false); got != "900150983cd24fb0d6963f7d28e17f72" {
		t.Errorf("ToHash(md5) = %s", got)
	}
	// sha1("abc") = a9993e364706816aba3e25717850c26c9cd0d89d
	if got := ToHash("abc", "sha1", false); got != "a9993e364706816aba3e25717850c26c9cd0d89d" {
		t.Errorf("ToHash(sha1) = %s", got)
	}
	// sha256("abc") = ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad
	if got := ToHash("abc", "sha256", false); got != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Errorf("ToHash(sha256) = %s", got)
	}
	// 不支持的类型不再 panic,返回错误提示
	if got := ToHash("abc", "md4", false); got != "unsupported hash type: md4" {
		t.Errorf("ToHash(md4) = %s", got)
	}
	if got := ToHash("abc", "ripemd160", false); got != "unsupported hash type: ripemd160" {
		t.Errorf("ToHash(ripemd160) = %s", got)
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

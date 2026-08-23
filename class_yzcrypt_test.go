package KzPack4Go

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestAESRoundTrip(t *testing.T) {
	key := []byte("0123456789abcdef") // 16字节密钥

	cases := [][]byte{
		[]byte("hello"),
		[]byte(""),
		[]byte("中文内容 äöü 😀"),
		bytes.Repeat([]byte("A"), 31), // 非整分组
		bytes.Repeat([]byte("B"), 32), // 整分组
	}

	for _, origin := range cases {
		crypted, err := YzCrypt.AES.Encrypt(key, origin)
		if err != nil {
			t.Fatalf("Encrypt(%q) error: %v", origin, err)
		}
		plain, err := YzCrypt.AES.Decrypt(key, crypted)
		if err != nil {
			t.Fatalf("Decrypt error: %v", err)
		}
		if !bytes.Equal(plain, origin) {
			t.Errorf("roundtrip mismatch: got %q want %q", plain, origin)
		}
	}
}

// 同一明文+同一密钥,两次加密必须产生不同密文(随机 IV 生效)
func TestAESRandomIV(t *testing.T) {
	key := []byte("0123456789abcdef")
	origin := []byte("same plaintext")

	c1, err := YzCrypt.AES.Encrypt(key, origin)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := YzCrypt.AES.Encrypt(key, origin)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(c1, c2) {
		t.Error("两次加密密文相同:随机IV未生效")
	}

	// 密文长度 = 16(前置IV) + 16*ceil(len/16)
	wantLen := 16 + ((len(origin)+15)/16)*16
	if len(c1) != wantLen {
		t.Errorf("密文长度 = %d; want %d(含前置IV)", len(c1), wantLen)
	}
}

// 旧接口保持旧格式(固定 IV=密钥前16字节),历史密文仍可解密
func TestAESLegacyCompat(t *testing.T) {
	key := []byte("0123456789abcdef")
	origin := []byte("legacy data")

	legacy, err := EncryptAES(key, origin)
	if err != nil {
		t.Fatal(err)
	}
	// 旧格式:无前置 IV,密文长度 == 填充后明文长度
	wantLen := ((len(origin) + 15) / 16) * 16
	if len(legacy) != wantLen {
		t.Errorf("旧格式密文长度 = %d; want %d(无前置IV)", len(legacy), wantLen)
	}

	plain, err := DecryptAES(key, legacy)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plain, origin) {
		t.Errorf("legacy roundtrip mismatch: got %q want %q", plain, origin)
	}

	// 旧格式固定 IV 特性:同明文同密钥两次加密结果一致(验证格式未变)
	c1, _ := EncryptAES(key, origin)
	if !bytes.Equal(c1, legacy) {
		t.Error("旧格式两次加密结果不同(固定IV特性丢失)")
	}
}

func TestSimplyEncryptDecryptAES(t *testing.T) {
	key := "0123456789abcdef"
	origin := "simply hex"

	hexText := SimplyEncryptAES(key, origin)
	if hexText == "" {
		t.Fatal("SimplyEncryptAES = empty")
	}
	if got := SimplyDecryptAES(key, hexText); got != origin {
		t.Errorf("SimplyDecryptAES = %q; want %q", got, origin)
	}
}

// 密文长度不足或长度不合法时应报错而非 panic
func TestAESDecryptMalformed(t *testing.T) {
	key := []byte("0123456789abcdef")

	if _, err := YzCrypt.AES.Decrypt(key, []byte("short")); err == nil {
		t.Error("长度不足的密文应返回错误")
	}
	if _, err := YzCrypt.AES.Decrypt(key, make([]byte, 16)); err == nil {
		t.Error("仅IV无分组体应返回错误")
	}
	// IV(16B)+21B:21 不是 16 的整数倍,修复前 CryptBlocks 直接 panic
	if _, err := YzCrypt.AES.Decrypt(key, make([]byte, 37)); err == nil {
		t.Error("非分组整数倍的密文应返回错误而非 panic")
	}

	// 旧接口同样不应 panic
	if _, err := DecryptAES(key, make([]byte, 21)); err == nil {
		t.Error("旧接口非分组整数倍的密文应返回错误而非 panic")
	}
}

// ---------- SM4 ----------

func TestSm4RoundTrip(t *testing.T) {
	key := []byte("0123456789abcdef") // 16字节密钥

	cases := [][]byte{
		[]byte("hello"),
		[]byte("中文内容"),
		bytes.Repeat([]byte("A"), 17), // 非整分组
	}
	for _, origin := range cases {
		cipherText, err := Sm4Encrypt(key, origin)
		if err != nil {
			t.Fatalf("Sm4Encrypt(%q) error: %v", origin, err)
		}
		plain, err := Sm4Decrypt(key, cipherText)
		if err != nil {
			t.Fatalf("Sm4Decrypt error: %v", err)
		}
		if !bytes.Equal(plain, origin) {
			t.Errorf("SM4 roundtrip mismatch: got %q want %q", plain, origin)
		}
	}
}

func TestYzCryptSM4(t *testing.T) {
	b64Key := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef"))

	cipher, err := YzCrypt.SM4.Encrypt(b64Key, "hello")
	if err != nil {
		t.Fatal(err)
	}
	if cipher == "" {
		t.Fatal("SM4.Encrypt 返回空")
	}
	text, err := YzCrypt.SM4.Decrypt(b64Key, cipher)
	if err != nil {
		t.Fatal(err)
	}
	if text != "hello" {
		t.Errorf("SM4 roundtrip = %q; want hello", text)
	}

	// 兼容旧接口
	cipher2, err := EncryptSM4(b64Key, "hello")
	if err != nil {
		t.Fatal(err)
	}
	if text2, err := DecryptSM4(b64Key, cipher2); err != nil || text2 != "hello" {
		t.Errorf("EncryptSM4/DecryptSM4 = %q, %v", text2, err)
	}
}

func TestSm4KeyLength(t *testing.T) {
	// 15 字节密钥(非法)应报错
	b64Bad := base64.StdEncoding.EncodeToString([]byte("0123456789abcde"))
	if _, err := YzCrypt.SM4.Encrypt(b64Bad, "x"); err == nil {
		t.Error("非法密钥长度应返回错误")
	}
	// 非法 Base64 密钥应报错
	if _, err := YzCrypt.SM4.Encrypt("!!!not-base64!!!", "x"); err == nil {
		t.Error("非法Base64密钥应返回错误")
	}
}

// ---------- 填充与旧接口边界 ----------

func TestPKCS7UnPaddingEdge(t *testing.T) {
	if _, err := PKCS7UnPadding([]byte{}); err == nil {
		t.Error("空数据应报错")
	}
	if _, err := PKCS7UnPadding([]byte{0x01, 0x02, 0x05}); err == nil {
		t.Error("填充长度超过数据长度应报错")
	}
	if got, err := PKCS7UnPadding([]byte("hello\x03\x03\x03")); err != nil || string(got) != "hello" {
		t.Errorf("正常去填充 = %q, %v", got, err)
	}
}

func TestSimplyAESEdge(t *testing.T) {
	if got := SimplyEncryptAES("0123456789abcdef", ""); got != "" {
		t.Errorf("空明文应返回空: %s", got)
	}
	if got := SimplyDecryptAES("0123456789abcdef", ""); got != "" {
		t.Errorf("空密文应返回空: %s", got)
	}
	// 非法 hex 应返回错误串而非 panic
	if got := SimplyDecryptAES("0123456789abcdef", "zz"); got == "" {
		t.Error("非法hex应返回错误串")
	}
	// 非法密钥应返回错误串
	if got := SimplyEncryptAES("short", "data"); got == "" {
		t.Error("非法密钥应返回错误串")
	}
}

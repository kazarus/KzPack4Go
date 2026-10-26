package KzPack4Go

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

// GenerateKey 生成的分量必须为定宽 32 字节,公钥为 04||X||Y 共 65 字节
func TestTEncConfigGenerateKey(t *testing.T) {
	var cfg TEncConfig
	if err := cfg.GenerateKey(); err != nil {
		t.Fatal(err)
	}

	for name, s := range map[string]string{"d": cfg.D, "x": cfg.X, "y": cfg.Y} {
		b, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			t.Fatalf("%s 分量不是合法 base64: %v", name, err)
		}
		if len(b) != 32 {
			t.Errorf("%s 分量长度 = %d; want 32", name, len(b))
		}
	}

	pub, err := base64.StdEncoding.DecodeString(cfg.EncryKey)
	if err != nil {
		t.Fatalf("公钥不是合法 base64: %v", err)
	}
	if len(pub) != 65 || pub[0] != 0x04 {
		t.Errorf("公钥格式非法: %d 字节, 头字节 %02x", len(pub), pub[0])
	}

	cfg.ClearAttr()
	if cfg.D != "" || cfg.X != "" || cfg.Y != "" {
		t.Error("ClearAttr 后分量未清空")
	}
}

// 分量重建私钥:与 GenerateKey 输出的 D/X/Y 一致
func TestGmKeyFromComponentsRoundTrip(t *testing.T) {
	var cfg TEncConfig
	if err := cfg.GenerateKey(); err != nil {
		t.Fatal(err)
	}

	priv, err := gmKeyFromComponents(cfg.D, cfg.X, cfg.Y)
	if err != nil {
		t.Fatal(err)
	}

	var src TEncConfig
	src.D = gmBigIntB64(priv.D)
	src.X = gmBigIntB64(priv.X)
	src.Y = gmBigIntB64(priv.Y)
	if src.D != cfg.D || src.X != cfg.X || src.Y != cfg.Y {
		t.Errorf("重建密钥分量不一致: got %s/%s/%s want %s/%s/%s",
			src.D, src.X, src.Y, cfg.D, cfg.X, cfg.Y)
	}
}

// 公钥分量与私钥不匹配(损坏)应报错;私钥分量非法也应报错
func TestGmKeyFromComponentsInvalid(t *testing.T) {
	var cfg TEncConfig
	if err := cfg.GenerateKey(); err != nil {
		t.Fatal(err)
	}

	// 篡改 y:与 D 推导的公钥不一致
	badY := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x01}, 32))
	_, err := gmKeyFromComponents(cfg.D, cfg.X, badY)
	if err == nil || !strings.Contains(err.Error(), "不匹配") {
		t.Errorf("公私钥分量不匹配应报错: %v", err)
	}

	// 私钥分量全 0:非法
	badD := base64.StdEncoding.EncodeToString(make([]byte, 32))
	if _, err := gmKeyFromComponents(badD, cfg.X, cfg.Y); err == nil {
		t.Error("全 0 私钥分量应报错")
	}

	// 分量长度非 32 字节
	short := base64.StdEncoding.EncodeToString([]byte{0x01})
	if _, err := gmKeyFromComponents(short, cfg.X, cfg.Y); err == nil {
		t.Error("非 32 字节分量应报错")
	}
}

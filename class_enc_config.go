package KzPack4Go

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/tjfoc/gmsm/sm2"
)

type TEncConfig struct {
	CryptTxt string `json:"cryptTxt"`
	EncryKey string `json:"encryKey"`
	D        string `json:"d"`
	X        string `json:"x"`
	Y        string `json:"y"`
}

func (this TEncConfig) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})

	result["cryptTxt"] = this.CryptTxt
	result["encryKey"] = this.EncryKey

	result["d"] = this.D
	result["x"] = this.X
	result["y"] = this.Y

	result["cipherBy"] = "SM2-C1C3C2/SM4-CBC"
	result["metaData"] = "随机生成 16 字节 SM4 密钥，用公钥 SM2 加密(C1C3C2)后 base64 放入 X-GM-KEY 请求头"

	return json.Marshal(result)
}

func (this *TEncConfig) ClearAttr() {
	this.D = ""
	this.X = ""
	this.Y = ""
}

func gmBigIntB64(n *big.Int) string {
	b := make([]byte, 32)
	n.FillBytes(b)
	return base64.StdEncoding.EncodeToString(b)
}

// gmDecodeComponent 密钥分量解码，兼容新旧编码：64 位纯 hex 视为旧
// hex 格式，其余按标准 base64；解码结果必须是定宽 32 字节。
func gmDecodeComponent(s string) (*big.Int, error) {
	var (
		b   []byte
		err error
	)
	b, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(b) != 32 {
		return nil, fmt.Errorf("密钥分量长度非法: %d 字节（应为 32）", len(b))
	}
	return new(big.Int).SetBytes(b), nil
}

// SM2PublicKeyBase64 返回非压缩格式（04||X||Y 共 65 字节）base64 编码的
// SM2 公钥，交给客户端做 SM4 会话密钥加密。不直接用 elliptic.Marshal
// （Go 已标记废弃），手工拼定点格式再编码。
func SM2PublicKeyBase64(priv *sm2.PrivateKey) string {
	x := make([]byte, 32)
	y := make([]byte, 32)
	priv.X.FillBytes(x)
	priv.Y.FillBytes(y)
	return base64.StdEncoding.EncodeToString(append([]byte{0x04}, append(x, y...)...))
}

func (this *TEncConfig) GenerateKey() error {

	PrivateKey, eror := sm2.GenerateKey(rand.Reader)
	if eror != nil {
		return eror
	}

	this.D = gmBigIntB64(PrivateKey.D)
	this.X = gmBigIntB64(PrivateKey.X)
	this.Y = gmBigIntB64(PrivateKey.Y)

	this.EncryKey = SM2PublicKeyBase64(PrivateKey)

	return eror
}

func gmKeyFromComponents(dStr, xStr, yStr string) (*sm2.PrivateKey, error) {
	d, err := gmDecodeComponent(dStr)
	if err != nil {
		return nil, err
	}
	x, err := gmDecodeComponent(xStr)
	if err != nil {
		return nil, err
	}
	y, err := gmDecodeComponent(yStr)
	if err != nil {
		return nil, err
	}
	curve := sm2.P256Sm2()
	if !curve.IsOnCurve(x, y) {
		return nil, fmt.Errorf("sm2 公钥点不在曲线上，密钥文件可能已损坏")
	}
	return &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{Curve: curve, X: x, Y: y},
		D:         d,
	}, nil
}

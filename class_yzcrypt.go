package KzPack4Go

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/tjfoc/gmsm/sm4"
)

// YzCrypt AES / SM4 加解密入口。
//
//	YzCrypt.AES.Encrypt(key, originData) // AES加密
//	YzCrypt.AES.Decrypt(key, crypted)    // AES解密
//	YzCrypt.SM4.Encrypt(key, data)       // SM4加密
//	YzCrypt.SM4.Decrypt(key, data)       // SM4解密
var YzCrypt = yzCrypt{
	AES: AES{},
	SM4: SM4{},
}

// yzCrypt AES / SM4 加解密容器
type yzCrypt struct {
	AES AES
	SM4 SM4
}

// AES AES加解密分组
type AES struct{}

// Encrypt AES加密
func (AES) Encrypt(key, originData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	originData = PKCS7DoPadding(originData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return crypted, nil
}

// Decrypt AES解密
func (AES) Decrypt(key, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(crypted))
	blockMode.CryptBlocks(originData, crypted)
	originData, err = PKCS7UnPadding(originData)
	if err != nil {
		return nil, err
	}
	return originData, nil
}

// SM4 SM4加解密分组
type SM4 struct{}

// Encrypt 加密字符串并返回Base64编码结果
func (SM4) Encrypt(key, data string) (string, error) {
	// 解码Base64格式的密钥
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("密钥Base64解码失败: %v", err)
	}

	// 检查密钥长度是否为16字节(128位)
	if len(keyBytes) != 16 {
		return "", fmt.Errorf("SM4密钥长度必须为16字节")
	}

	// 执行加密
	cipherText, err := Sm4Encrypt(keyBytes, []byte(data))
	if err != nil {
		return "", fmt.Errorf("加密失败: %v", err)
	}

	// 返回Base64编码的结果
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt 解密Base64编码的字符串
func (SM4) Decrypt(key, data string) (string, error) {
	// 解码Base64格式的密钥
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("密钥Base64解码失败: %v", err)
	}

	// 检查密钥长度是否为16字节(128位)
	if len(keyBytes) != 16 {
		return "", fmt.Errorf("SM4密钥长度必须为16字节")
	}

	// 解码Base64格式的密文
	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("密文Base64解码失败: %v", err)
	}

	// 执行解密
	plainText, err := Sm4Decrypt(keyBytes, cipherText)
	if err != nil {
		return "", fmt.Errorf("解密失败: %v", err)
	}

	return string(plainText), nil
}

func PKCS7DoPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padtext...)
}

func PKCS7UnPadding(originData []byte) ([]byte, error) {
	/*
		length := len(originData)
		unpadding := int(originData[length-1])
		return originData[:(length - unpadding)]
	*/
	length := len(originData)
	if length == 0 {
		return nil, fmt.Errorf("数据长度为0")
	}

	padding := int(originData[length-1])
	if padding > length {
		return nil, fmt.Errorf("无效的填充长度")
	}

	return originData[:length-padding], nil
}

// EncryptAES AES加密 (兼容旧接口, 等价于 YzCrypt.AES.Encrypt)
func EncryptAES(key, originData []byte) ([]byte, error) {
	return YzCrypt.AES.Encrypt(key, originData)
}

// DecryptAES AES解密 (兼容旧接口, 等价于 YzCrypt.AES.Decrypt)
func DecryptAES(key, crypted []byte) ([]byte, error) {
	return YzCrypt.AES.Decrypt(key, crypted)
}

// AES加密
func SimplyEncryptAES(key, originData string) string {

	if originData == "" {
		return ""
	}

	result, eror := EncryptAES([]byte(key), []byte(originData))
	if eror != nil {
		return eror.Error()
	}

	return hex.EncodeToString(result)
}

// AES解密
func SimplyDecryptAES(key, originData string) string {

	if originData == "" {
		return ""
	}

	byteData, eror := hex.DecodeString(originData)
	if eror != nil {
		return eror.Error()
	}

	result, eror := DecryptAES([]byte(key), byteData)
	if eror != nil {
		return eror.Error()
	}

	return string(result)
}

// SM4加密
func Sm4Encrypt(key, plaintext []byte) ([]byte, error) {
	// 创建SM4加密器
	cipher, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 对明文进行PKCS#7填充，使其长度为16的倍数
	paddedText := PKCS7DoPadding(plaintext, sm4.BlockSize)

	// 创建输出缓冲区
	ciphertext := make([]byte, len(paddedText))

	// 分组加密(ECB模式)
	for i := 0; i < len(paddedText); i += sm4.BlockSize {
		cipher.Encrypt(ciphertext[i:i+sm4.BlockSize], paddedText[i:i+sm4.BlockSize])
	}

	return ciphertext, nil
}

// SM4解密
func Sm4Decrypt(key, ciphertext []byte) ([]byte, error) {
	// 创建SM4解密器
	cipher, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建输出缓冲区
	plaintext := make([]byte, len(ciphertext))

	// 分组解密(ECB模式)
	for i := 0; i < len(ciphertext); i += sm4.BlockSize {
		cipher.Decrypt(plaintext[i:i+sm4.BlockSize], ciphertext[i:i+sm4.BlockSize])
	}

	// 去除PKCS#7填充
	return PKCS7UnPadding(plaintext)
}

// EncryptSM4 加密字符串并返回Base64编码结果 (兼容旧接口, 等价于 YzCrypt.SM4.Encrypt)
func EncryptSM4(key, data string) (string, error) {
	return YzCrypt.SM4.Encrypt(key, data)
}

// DecryptSM4 解密Base64编码的字符串 (兼容旧接口, 等价于 YzCrypt.SM4.Decrypt)
func DecryptSM4(key, data string) (string, error) {
	return YzCrypt.SM4.Decrypt(key, data)
}

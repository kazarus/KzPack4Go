package KzPack4Go

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/emmansun/gmsm/sm4"
)

// YzCrypt AES / SM4 加解密入口。
//
//	YzCrypt.AES.Encrypt(key, originData) // AES加密
//	YzCrypt.AES.Decrypt(key, crypted)    // AES解密
//	YzCrypt.SM4.Encrypt(key, data)       // SM4加密
//	YzCrypt.SM4.Decrypt(key, data)       // SM4解密
var YzCrypt = YzCryptClass{
	AES: AES{},
	SM4: SM4{},
}

// YzCrypt AES / SM4 加解密容器
type YzCryptClass struct {
	AES AES
	SM4 SM4
}

// AES AES加解密分组
type AES struct{}

// Encrypt AES加密(CBC/PKCS7)。
// 每次加密随机生成 16 字节 IV 并前置到密文头部,避免固定 IV 导致
// 相同明文前缀产生相同密文前缀(前缀泄露)。密文格式: IV(16B) || 密文本体。
func (AES) Encrypt(key, originData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()

	iv := make([]byte, blockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("生成随机IV失败: %v", err)
	}

	padded := PKCS7DoPadding(originData, blockSize)
	crypted := make([]byte, blockSize+len(padded))
	copy(crypted, iv)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(crypted[blockSize:], padded)
	return crypted, nil
}

// Decrypt AES解密(CBC/PKCS7)。
// 密文头部 16 字节为加密时前置的随机 IV,其余为密文本体。
func (AES) Decrypt(key, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(crypted) < 2*blockSize {
		return nil, fmt.Errorf("密文长度不足:至少需要 IV(%dB)+一个分组(%dB)", blockSize, blockSize)
	}
	// #CBC CryptBlocks 对非分组整数倍的输入会 panic,提前拦截
	if (len(crypted)-blockSize)%blockSize != 0 {
		return nil, fmt.Errorf("密文本体长度不是分组(%dB)的整数倍", blockSize)
	}

	iv := crypted[:blockSize]
	body := crypted[blockSize:]

	blockMode := cipher.NewCBCDecrypter(block, iv)
	originData := make([]byte, len(body))
	blockMode.CryptBlocks(originData, body)

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

// aesEncryptLegacy AES加密(旧格式:CBC,固定 IV=密钥前16字节)。
// ⚠️ 不安全:固定 IV 导致相同明文前缀产生相同密文前缀;仅用于兼容历史数据,
// 新数据请使用 YzCrypt.AES.Encrypt(随机 IV 前置)。
func aesEncryptLegacy(key, originData []byte) ([]byte, error) {
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

// aesDecryptLegacy AES解密(旧格式:固定 IV=密钥前16字节)。
// 仅用于解密旧接口加密的历史密文;新数据请使用 YzCrypt.AES.Decrypt。
func aesDecryptLegacy(key, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	// #CBC CryptBlocks 对非分组整数倍的输入会 panic,提前拦截
	if len(crypted)%blockSize != 0 {
		return nil, fmt.Errorf("密文长度不是分组(%dB)的整数倍", blockSize)
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(crypted))
	blockMode.CryptBlocks(originData, crypted)
	originData, err = PKCS7UnPadding(originData)
	if err != nil {
		return nil, err
	}
	return originData, nil
}

// EncryptAES AES加密 (兼容旧接口;⚠️ 旧格式:固定 IV=密钥前16字节,不安全,
// 新数据请用 YzCrypt.AES.Encrypt;本函数仅保证旧调用点与新格式互不干扰)
func EncryptAES(key, originData []byte) ([]byte, error) {
	return aesEncryptLegacy(key, originData)
}

// DecryptAES AES解密 (兼容旧接口;仅能解密旧格式密文,历史数据请用本函数,
// 新数据请用 YzCrypt.AES.Decrypt)
func DecryptAES(key, crypted []byte) ([]byte, error) {
	return aesDecryptLegacy(key, crypted)
}

// AES加密(旧格式,hex 字符串,与 EncryptAES 配套)
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

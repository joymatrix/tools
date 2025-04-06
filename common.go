package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"time"
)

func ConverTimeToTimestamp(timeStr string) (int64, error) {
	layout := "2006-01-02 15:04:05"
	// 解析时间字符串
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Printf("解析时间字符串时出错: %v\n", err)
		return -1, err
	}
	// 获取时间戳
	timestamp := parsedTime.Unix()
	return timestamp, nil
}

func EncryptURL(url string) string {
	return base64.StdEncoding.EncodeToString([]byte(url))
}

func DecryptURL(encryptedURL string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encryptedURL)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// PKCS7Padding 填充数据
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7UnPadding 去除填充数据
func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// AESEncrypt CBC 模式加密
func AESEncrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	ciphertext := make([]byte, len(plaintext))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt CBC 模式解密
func AESDecrypt(ciphertext string, key []byte) ([]byte, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	plaintext := make([]byte, len(ciphertextBytes))
	blockMode.CryptBlocks(plaintext, ciphertextBytes)
	plaintext = PKCS7UnPadding(plaintext)

	return plaintext, nil
}

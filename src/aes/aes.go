package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/coderant163/docSyncKit/src/logger"
)

// TODO 本文件暂时没被使用，后续删除

// =================== CBC ======================
func Encrypt(orig, secret string) (encrypted string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Sugar().Errorf("Encrypt [%s] fail, recover from %+v", orig, err)
		}
	}()
	origData := []byte(orig)
	key := []byte(secret)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encryptedByte := make([]byte, len(origData))                // 创建数组
	blockMode.CryptBlocks(encryptedByte, origData)              // 加密
	return base64.StdEncoding.EncodeToString(encryptedByte)
}
func Decrypt(encrypted string, secret string) string {
	defer func() {
		if err := recover(); err != nil {
			logger.Sugar().Errorf("Decrypt [%s] fail, recover from %+v", encrypted, err)
		}
	}()
	key := []byte(secret)
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(encrypted)
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decryptedByte := make([]byte, len(crytedByte))              // 创建数组
	blockMode.CryptBlocks(decryptedByte, crytedByte)            // 解密
	decryptedByte = pkcs5UnPadding(decryptedByte)               // 去除补全码
	return string(decryptedByte)
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

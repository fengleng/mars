package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/gososy/sorpc/log"
)

const (
	Pkcs5Padding = 0
	Pkcs7Padding = 1
)

func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Warn(err)
	}
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}

func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)     // 分组秘钥
	blockSize := block.BlockSize()     // 获取秘钥块的长度
	if len(encrypted)%blockSize != 0 { // 解密长度错误
		return []byte{}
		//不做判断 CryptBlocks 会触发 panic("crypto/cipher: input not full blocks")
	}

	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}

// ECB 加解密 pkcs7Padding 填充方式
func AesDecryptEcb(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}
	return pkcs7UnPadding(decrypted)
}
func AesEncryptEcb(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	data = pkcs7Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}
	return decrypted
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	}
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncryptEcbBySelectedPadding(data, key []byte, paddingType int) []byte {
	block, _ := aes.NewCipher(key)
	switch paddingType {
	case Pkcs5Padding:
		data = pkcs5Padding(data, block.BlockSize())
	case Pkcs7Padding:
		data = pkcs7Padding(data, block.BlockSize())
	default:
		panic("Invalid padding type.")
	}

	encrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(encrypted[bs:be], data[bs:be])
	}
	return encrypted
}

func AesDecryptEcbBySelectedPadding(data, key []byte, paddingType int) []byte {
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}

	switch paddingType {
	case Pkcs5Padding:
		return pkcs5UnPadding(decrypted)
	case Pkcs7Padding:
		return pkcs7UnPadding(decrypted)
	}

	panic("Invalid padding type.")
}

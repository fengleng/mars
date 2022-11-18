package encryAlgorithm

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"github.com/gososy/sorpc/log"
)

//3DES加密
func TripleEncrypt(origData []byte, key []byte) (string, error) {
	//通过3DES库产生分组模块
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	//补码
	sourceData := ZeroPadding(origData, block.BlockSize())
	//设置加密模式
	iv := make([]byte, 8)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//创建缓冲区
	crypted := make([]byte, len(sourceData))
	//加密
	blockMode.CryptBlocks(crypted, sourceData)
	return hex.EncodeToString(crypted), nil
}
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	if len(ciphertext)%blockSize != 0 {
		padding := blockSize - len(ciphertext)%blockSize
		padtext := bytes.Repeat([]byte{0}, padding)
		return append(ciphertext, padtext...)
	}
	return ciphertext
}

//3DES解密
func TrippleDesDecrypt(crypted []byte, key []byte) []byte {
	block, _ := des.NewTripleDESCipher(key)
	//设置解密方式
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	// 创建缓冲区
	sourceData := make([]byte, len(crypted))
	//解密
	blockMode.CryptBlocks(sourceData, crypted)
	//去码
	sourceData = PKCS5UnPadding(sourceData)
	return sourceData
}

//补码
func PKCS5Padding(sourceData []byte, blockSize int) []byte {
	padding := blockSize - len(sourceData)%8
	paddTxt := bytes.Repeat([]byte{byte(padding)}, padding)
	sourceData = append(sourceData, paddTxt...)
	return sourceData
}

//去码
func PKCS5UnPadding(destinateData []byte) []byte {
	length := len(destinateData)
	padding := int(destinateData[len(destinateData)-1])
	sourceData := destinateData[:length-padding]
	return sourceData
}

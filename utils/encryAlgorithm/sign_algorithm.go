package encryAlgorithm

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gososy/sorpc/log"
	"github.com/pkg/errors"
	"sort"
)

// 签名算法
const (
	SignTypeRsa    = 1
	SignTypeMd5    = 2
	SignTypeSha1   = 3
	SignTypeRsa256 = 4
)

// 对签名内容进行排序
func SortSignData(mapBody map[string]interface{}) string {
	sortedKeys := make([]string, 0)
	for k := range mapBody {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	var src string
	index := 0
	for _, k := range sortedKeys {
		value := fmt.Sprintf("%v", mapBody[k])
		if value != "" && value != "<nil>" {
			src = src + k + "=" + value
		}
		//最后一项后面不要&
		if value != "" && index < len(sortedKeys)-1 {
			src = src + "&"
		}
		index++
	}
	return src
}
func GetRsaSha256Sign(privateKey, signData string) (string, error) {
	return getSign(SignTypeRsa256, privateKey, signData, nil)
}

// signMap 可以为 nil ， 传入 signMap，会自动根据 key 进行排序
func GetRsaSign(publicKey, signData string, signMap map[string]interface{}) (string, error) {
	return getSign(SignTypeRsa, publicKey, signData, signMap)
}
func GetMd5Sign(signData string, signMap map[string]interface{}) (string, error) {
	return getSign(SignTypeMd5, "", signData, signMap)
}
func GetSha1Sign(secretKey, signData string, signMap map[string]interface{}) (string, error) {
	return getSign(SignTypeSha1, secretKey, signData, signMap)
}

// 获取签名, signData signMap 其中一个必须有值
func getSign(signType uint32, publicKey, signData string, signMap map[string]interface{}) (string, error) {
	if signMap != nil {
		signData = SortSignData(signMap)
	}
	if len(signData) == 0 {
		return "", errors.New("signData or signMap is nil")
	}
	switch signType {
	case SignTypeRsa:
		b, err := RsaEncrypt(publicKey, signData)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}
		return base64.StdEncoding.EncodeToString(b), nil
	case SignTypeRsa256:
		b := RsaSignWithSha256([]byte(signData), []byte(publicKey))
		return base64.StdEncoding.EncodeToString(b), nil
	case SignTypeMd5:
		return StrMd5(signData), nil
	case SignTypeSha1:
		mac := hmac.New(sha1.New, []byte(publicKey))
		mac.Write([]byte(signData))
		signData := mac.Sum(nil)
		return base64.StdEncoding.EncodeToString(signData), nil
	}
	return "", errors.New("invalid sign type")
}

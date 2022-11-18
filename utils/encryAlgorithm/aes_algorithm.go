package encryAlgorithm

import (
	"encoding/base64"
	"strings"

	"github.com/forgoer/openssl"
	"github.com/gososy/sorpc/log"
)

const (
	AesTypeECB    = 1
	AesTypeCBC    = 2
	AesPaddingCS5 = "PKCS5"
	AesPaddingCS7 = "PKCS7"
)

func AesDecrypt(secretKey, secretData, ivKey string, AesType uint32, padding string) ([]byte, error) {
	iv, _ := base64.StdEncoding.DecodeString(ivKey)
	key, _ := base64.StdEncoding.DecodeString(secretKey)
	data, _ := Base64URLDecode(secretData)
	switch AesType {
	case AesTypeECB:
		return openssl.AesECBDecrypt(data, key, padding)
	case AesTypeCBC:
		return openssl.AesCBCDecrypt(data, key, iv, padding)
	}
	return nil, nil
}
func AesEncrypt(secretKey, origData, ivKey string, AesType uint32, padding string) (string, error) {
	iv, _ := base64.StdEncoding.DecodeString(ivKey)
	key, _ := base64.StdEncoding.DecodeString(secretKey)
	data := []byte(origData)
	switch AesType {
	case AesTypeECB:
		str, err := openssl.AesECBEncrypt(data, key, padding)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}
		return base64.StdEncoding.EncodeToString(str), nil
	case AesTypeCBC:
		str, err := openssl.AesCBCEncrypt(data, key, iv, padding)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}
		return base64.StdEncoding.EncodeToString(str), nil
	}
	return "", nil
}
func Base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	//res,  := base64.URLEncoding.DecodeString(data)
	return base64.URLEncoding.DecodeString(data)
}
func Base64UrlSafeEncode(source []byte) string {
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
	byteArr := base64.StdEncoding.EncodeToString(source)
	safeUrl := strings.Replace(string(byteArr), "/", "_", -1)
	safeUrl = strings.Replace(safeUrl, "+", "-", -1)
	safeUrl = strings.Replace(safeUrl, "=", "", -1)
	return safeUrl
}

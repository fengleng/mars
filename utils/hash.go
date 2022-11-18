package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

func StrMd5(s string) string {
	m := md5.New()
	_, _ = m.Write([]byte(s))
	return fmt.Sprintf("%x", m.Sum(nil))
}

func StrSha256(dst string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(dst))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func StrSha256WithSalt(dst string, salt string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(dst + salt))
	return fmt.Sprintf("%x", h.Sum(nil))
}

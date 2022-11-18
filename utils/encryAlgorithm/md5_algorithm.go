package encryAlgorithm

import (
	"crypto/md5"
	"fmt"
)

func StrMd5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return fmt.Sprintf("%x", m.Sum(nil))
}

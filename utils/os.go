package utils

import (
	"os"
)

func GetPwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

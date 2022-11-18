package utils

import (
	"os"
	"runtime"
	"strings"
)

var sep string
var sepEnvPath string

func init() {
	if runtime.GOOS == "windows" {
		sep = `\`
		sepEnvPath = `;`
	} else {
		sep = "/"
		sepEnvPath = ":"
	}
}

func AdjPathSep(src string) string {
	var f, t string
	if runtime.GOOS == "windows" {
		f = `/`
		t = `\`
	} else {
		f = `\`
		t = `/`
	}
	return strings.Replace(src, f, t, -1)
}

func Basename(name string) string {
	i := len(name) - 1
	// Remove trailing slashes
	for ; i > 0 && name[i] == '/'; i-- {
		name = name[:i]
	}
	// Remove leading directory name
	for i--; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}

	return name
}

func FilePathSplit(path string) (dirPath string, fileName string) {
	i := strings.LastIndex(path, sep)
	return path[:i+1], path[i+1:]
}

func FileSuffix(fileName string) string {
	list := strings.Split(fileName, ".")
	idx := len(list) - 1
	if idx < 0 {
		idx = 0
	}

	return list[idx]
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

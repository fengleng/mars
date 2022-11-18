//go:build !linux
// +build !linux

package kfklog

import (
	"os"
	"path/filepath"
)

var logPath = filepath.Join(os.TempDir(), "brick", "kafka", "kafka.log")

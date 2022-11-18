package utils

import (
	"context"
	"net"
)

func IsTimeoutError(err error) bool {
	return err == context.DeadlineExceeded
}

func IsHttpTimeoutError(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}
	return false
}

package utils

import (
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func MustNotEmpty(args ...string) bool {
	for _, s := range args {
		if len(s) == 0 {
			return false
		}
	}
	return true
}
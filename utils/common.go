package utils

import (
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
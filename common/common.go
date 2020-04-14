package common

import (
	"strconv"
	"time"
)

func GenerateOrderId() string {
	now := time.Now().Format("20060102")
	unix := time.Now().UnixNano()
	now2 := strconv.FormatUint(uint64(unix), 10)
	return now + now2
}

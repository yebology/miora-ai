package clients

import (
	"strconv"
	"strings"
)

func hexToUint64(hex string) uint64 {
	hex = strings.TrimPrefix(hex, "0x")
	val, _ := strconv.ParseUint(hex, 16, 64)
	return val
}

package clients

import (
	"strconv"
	"strings"
)

// hexToUint64 converts a hex string (e.g. "0x1a4") to uint64.
// Returns 0 if the input is invalid.
func hexToUint64(hex string) uint64 {

	hex = strings.TrimPrefix(hex, "0x")
	val, _ := strconv.ParseUint(hex, 16, 64)
	return val

}

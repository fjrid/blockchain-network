package util

import (
	"encoding/binary"
)

func Uint64ToBytes(val uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, val)

	return bytes
}

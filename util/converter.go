package util

import (
	"encoding/binary"
	"math"
)

func Float64ToBytes(val float64) []byte {
	bits := math.Float64bits(val)

	bytes := make([]byte, 64)
	binary.BigEndian.PutUint64(bytes, bits)

	return bytes
}

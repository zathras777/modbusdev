package modbusdev

import (
	"math"
)

func unsigned16(vals []byte) uint16 {
	return uint16(vals[0])<<8 + uint16(vals[1])
}

func signed16(vals []byte) int16 {
	u := unsigned16(vals)
	if u > 32767 {
		return int16(u - 1 - 65535)
	}
	return int16(u)
}

func unsigned32(vals []byte) uint32 {
	return uint32(vals[0])<<24 + uint32(vals[1])<<16 + uint32(vals[2])<<8 + uint32(vals[3])
}

func signed32(vals []byte) int32 {
	u := unsigned32(vals)
	if u > 2147483647 {
		return int32(u - 1 - 4294967295)
	}
	return int32(u)
}

func bool16(vals []byte) bool {
	return vals[1]&0x01 == 0x01
}

func ieee32(vals []byte) float64 {
	u := unsigned32(vals)
	sign := u >> 31
	exp := float64(u>>23&0xff) - 0x7f
	rem := uint64(u & 0x7fffff)
	var bottom uint64
	if exp != 0 {
		bottom = 0x800000
	} else {
		bottom = 0x400000
	}
	mant := float64(rem)/float64(bottom) + 1

	if sign == 0 {
		return mant * math.Exp2(exp)
	}
	return -1 * mant * math.Exp2(exp)
}

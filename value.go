package modbusdev

import (
	"log"
)

// Value As there are a number of possible return values, we simply
// return this structure with the appropriate member set.
type Value struct {
	Unsigned16 uint16
	Signed16   int16
	Unsigned32 uint32
	Signed32   int32
	Coil       bool
	Ieee32     float64
}

func (val *Value) formatBytes(format string, value []byte) {
	switch format {
	case "u16":
		val.Unsigned16 = unsigned16(value)
	case "s16":
		val.Signed16 = signed16(value)
	case "u32":
		val.Unsigned32 = unsigned32(value)
	case "s32":
		val.Signed32 = signed32(value)
	case "ieee32":
		val.Ieee32 = ieee32(value)
	case "coil":
		val.Coil = bool16(value)
	}
}

func (val *Value) asBytes(format string) (result []byte) {
	var err error
	switch format {
	case "u16":
		result, err = formatIntAsBytes(format, int(val.Unsigned16))
	case "s16":
		result, err = formatIntAsBytes(format, int(val.Signed16))
	case "u32":
		result, err = formatIntAsBytes(format, int(val.Unsigned32))
	case "s32":
		result, err = formatIntAsBytes(format, int(val.Signed32))
	//	case "ieee32":
	//		val.Ieee32 = ieee32(value)
	case "coil":
		result = make([]byte, 2)
		if val.Coil {
			result[1] = 0x1
		}
	}
	if err != nil {
		log.Print(err)
	}
	return
}

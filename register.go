package modbusdev

// Register Structure that contains details of the register value available.
type Register struct {
	description string
	units       string
	register    uint16
	format      string
	factor      float64
}

func (r Register) registersRqd() uint16 {
	switch r.format {
	case "u16", "s16", "coil":
		return 1
	case "u32", "s32", "ieee32":
		return 2
	}
	return 2
}

func (r Register) applyFactor(val *Value) {
	switch r.format {
	case "u16":
		val.Ieee32 = r.factor * float64(val.Unsigned16)
	case "s16":
		val.Ieee32 = r.factor * float64(val.Signed16)
	case "u32":
		val.Ieee32 = r.factor * float64(val.Unsigned32)
	case "s32":
		val.Ieee32 = r.factor * float64(val.Signed32)
	case "ieee32":
		val.Ieee32 = r.factor * val.Ieee32
	}
}

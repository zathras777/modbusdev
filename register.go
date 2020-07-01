package modbusdev

// Register Structure that contains details of the register value available.
type Register struct {
	description string
	units       string
	register    uint16
	format      string
	factor      float64
}

type registerCache struct {
	start        uint16
	qty          uint16
	registerData []byte
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

func (r Register) maxRegister() uint16 {
	return r.register + r.registersRqd()
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

func (rc *registerCache) init() {
	rc.start = 65535
}

func (rc *registerCache) update(reg Register) {
	if reg.register < rc.start {
		rc.start = reg.register
	}
	if reg.maxRegister() > rc.start+rc.qty {
		rc.qty = reg.maxRegister() - rc.start
	}
}

func (rc *registerCache) getValue(reg Register) Value {
	idx := (reg.register - rc.start) * 2
	var val Value
	val.FormatBytes(reg.format, rc.registerData[idx:idx+reg.registersRqd()*2])
	return val
}

package modbusdev

// Register Structure that contains details of the register value available.
type Register struct {
	Description string
	Units       string
	Register    uint16
	Format      string
	Factor      float64
}

type registerCache struct {
	start        uint16
	qty          uint16
	registerData map[int]byte
}

func (r Register) registersRqd() uint16 {
	switch r.Format {
	case "u16", "s16", "coil":
		return 1
	case "u32", "s32", "ieee32":
		return 2
	}
	return 2
}

func (r Register) maxRegister() uint16 {
	return r.Register + r.registersRqd()
}

func (r Register) applyFactor(val *Value) {
	switch r.Format {
	case "u16":
		val.Ieee32 = r.Factor * float64(val.Unsigned16)
	case "s16":
		val.Ieee32 = r.Factor * float64(val.Signed16)
	case "u32":
		val.Ieee32 = r.Factor * float64(val.Unsigned32)
	case "s32":
		val.Ieee32 = r.Factor * float64(val.Signed32)
	case "ieee32":
		val.Ieee32 = r.Factor * val.Ieee32
	}
}

func (rc *registerCache) init() {
	rc.registerData = make(map[int]byte)
	rc.start = 65535
}

func (rc *registerCache) update(reg Register) {
	if reg.Register < rc.start {
		rc.start = reg.Register
	}
	if reg.maxRegister() > rc.start+rc.qty {
		rc.qty = reg.maxRegister() - rc.start
	}
}

func (rc *registerCache) updateBytes(offset uint16, newBytes []byte) {
	idx := int(rc.start+offset) * 2
	for i, bb := range newBytes {
		rc.registerData[idx+i] = bb
	}
}

func (rc *registerCache) getValue(reg Register) Value {
	idx := int(reg.Register-rc.start) * 2
	sz := int(reg.registersRqd() * 2)
	rawBytes := make([]byte, sz)
	for i := 0; i < sz; i++ {
		rawBytes[i] = rc.registerData[idx+i]
	}
	var val Value
	val.FormatBytes(reg.Format, rawBytes)
	return val
}

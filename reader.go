package modbusdev

import (
	"fmt"
	"sort"
	"strings"

	"github.com/goburrow/modbus"
)

const (
	ieeeFmt = "%12.2f"
	numFmt  = "%15d"
	baseFmt = "  %5d: %-40s "
)

// Reader A reader structure allows us to tie a client to a device register map
type Reader struct {
	client    modbus.Client
	registers map[int]Register
}

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

// NewReader Return a configured Reader with the correct register mappings.
// Device names are converted to lower case for matching, so case provided is irrelevant.
func NewReader(client modbus.Client, device string) (rdr Reader, err error) {
	rdr.client = client
	switch strings.ToLower(device) {
	case "sdm230":
		rdr.registers = sdm230
	case "solaxx1hybrid":
		rdr.registers = solaxX1Hybrid
	default:
		err = fmt.Errorf("Device '%s' is not known. Add the details and then update reader.go to include it", device)
	}
	return
}

// ReadRegister Read the register specified by the code.
func (rdr *Reader) ReadRegister(code int, factored bool) (val Value, err error) {
	reg, ck := rdr.registers[code]
	if !ck {
		return val, fmt.Errorf("Code %d is not available", code)
	}
	var results []byte
	nRqd := reg.registersRqd()

	var typ int
	for typ = code; typ >= 10; typ = typ / 10 {
	}
	switch typ {
	case 3:
		results, err = rdr.client.ReadInputRegisters(reg.register, nRqd)
	case 4:
		results, err = rdr.client.ReadHoldingRegisters(reg.register, nRqd)
	}

	if err != nil {
		return val, err
	}
	switch reg.format {
	case "u16":
		val.Unsigned16 = unsigned16(results)
	case "s16":
		val.Signed16 = signed16(results)
	case "u32":
		val.Unsigned32 = unsigned32(results)
	case "s32":
		val.Signed32 = signed32(results)
	case "ieee32":
		val.Ieee32 = ieee32(results)
	case "coil":
		val.Coil = bool16(results)
	}
	if factored {
		reg.applyFactor(&val)
	}
	return val, nil
}

// Units For the given register code, return the units specified
func (rdr *Reader) Units(code int) string {
	reg, ck := rdr.registers[code]
	if !ck {
		return ""
	}
	return reg.units
}

// Map Return a map object of the registers. If getting a register returns a value it is
// simply omitted from the map.
func (rdr *Reader) Map(factored bool) map[int]Value {
	mapValues := make(map[int]Value)
	var keys []int
	for k := range rdr.registers {
		keys = append(keys, k)
	}

	for _, code := range keys {
		val, err := rdr.ReadRegister(code, factored)
		if err != nil {
			continue
		}
		mapValues[code] = val
	}
	return mapValues
}

// Dump Query all defined registers and print the results to stdout.
func (rdr *Reader) Dump(factored bool) {
	var keys []int
	for k := range rdr.registers {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, code := range keys {
		reg := rdr.registers[code]
		val, err := rdr.ReadRegister(code, factored)
		if err != nil {
			fmt.Printf(baseFmt+"ERROR %s\n", code, reg.description, err)
			continue
		}

		if factored {
			fmt.Printf(baseFmt+ieeeFmt+" %s\n", code, reg.description, val.Ieee32, reg.units)
			continue
		}

		switch reg.format {
		case "u16":
			fmt.Printf(baseFmt+numFmt+" %s\n", code, reg.description, val.Unsigned16, reg.units)
		case "s16":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.description, val.Signed16, reg.units)
		case "u32":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.description, val.Unsigned32, reg.units)
		case "s32":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.description, val.Signed32, reg.units)
		case "ieee32":
			fmt.Printf(baseFmt+ieeeFmt+"%s\n", code, reg.description, val.Ieee32, reg.units)
		case "coil":
			fmt.Printf(baseFmt+"%t\n", code, reg.description, val.Coil)
		}
	}
}

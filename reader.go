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
	unsigned16 uint16
	signed16   int16
	unsigned32 uint32
	signed32   int32
	coil       bool
	ieee32     float64
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
	if code > 29999 && code < 39999 {
		results, err = rdr.client.ReadInputRegisters(reg.register, nRqd)
	}

	if err != nil {
		return val, err
	}
	switch reg.format {
	case "u16":
		val.unsigned16 = unsigned16(results)
	case "s16":
		val.signed16 = signed16(results)
	case "u32":
		val.unsigned32 = unsigned32(results)
	case "s32":
		val.signed32 = signed32(results)
	case "ieee32":
		val.ieee32 = ieee32(results)
	case "coil":
		val.coil = bool16(results)
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
			fmt.Printf(baseFmt+ieeeFmt+" %s\n", code, reg.description, val.ieee32, reg.units)
			continue
		}

		switch reg.format {
		case "u16":
			fmt.Printf(baseFmt+numFmt+" %s\n", code, reg.description, val.unsigned16, reg.units)
		case "s16":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.description, val.signed16, reg.units)
		case "u32":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.description, val.unsigned32, reg.units)
		case "s32":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.description, val.signed32, reg.units)
		case "ieee32":
			fmt.Printf(baseFmt+ieeeFmt+"%s\n", code, reg.description, val.ieee32, reg.units)
		case "coil":
			fmt.Printf(baseFmt+"%t\n", code, reg.description, val.coil)
		}
	}
}

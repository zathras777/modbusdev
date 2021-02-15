package modbusdev

import (
	"fmt"
	"log"
	"sort"

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
	holding   registerCache
	input     registerCache
}

// NewReader Return a configured Reader with the correct register mappings.
// Device names are converted to lower case for matching, so case provided is irrelevant.
func NewReader(client modbus.Client, device string) (rdr Reader, err error) {
	rdr.client = client
	regs, err := RegistersByName(device)
	if err != nil {
		return
	}
	rdr.registers = regs
	rdr.input.init()
	rdr.holding.init()

	for code, reg := range rdr.registers {
		switch getRegisterType(code) {
		case 3:
			rdr.input.update(reg)
		case 4:
			rdr.holding.update(reg)
		}
	}

	return
}

// ReadRegister Read the register specified by the code. This always causes the device to be
// queried.
func (rdr *Reader) ReadRegister(code int, factored bool) (val Value, err error) {
	reg, ck := rdr.registers[code]
	if !ck {
		return val, fmt.Errorf("Code %d is not available", code)
	}
	var results []byte
	nRqd := reg.registersRqd()
	switch getRegisterType(code) {
	case 3:
		results, err = rdr.client.ReadInputRegisters(reg.Register, nRqd)
	case 4:
		results, err = rdr.client.ReadHoldingRegisters(reg.Register, nRqd)
	}

	if err != nil {
		return val, err
	}
	val.FormatBytes(reg.Format, results)
	if factored {
		reg.applyFactor(&val)
	}
	return val, nil
}

func min(a, b uint16) uint16 {
	if a > b {
		return b
	}
	return a
}

// Read Read the registers that are required to provide data for the configured device. This
// attempts a single call to the device.
func (rdr *Reader) Read() error {
	regsRead := uint16(0)
	totalRead := regsRead

	if rdr.input.start != 65535 {
		for {
			toRead := min(125, rdr.input.qty-regsRead+1)
			results, err := rdr.client.ReadInputRegisters(rdr.input.start+regsRead, toRead)
			if err != nil {
				return err
			}

			rdr.input.updateBytes(regsRead, results)
			regsRead += toRead
			if regsRead >= rdr.input.qty {
				break
			}
		}
	}
	totalRead += regsRead
	if rdr.holding.start != 65535 {
		regsRead = 0
		for {
			toRead := min(125, rdr.holding.qty-regsRead+1)
			results, err := rdr.client.ReadHoldingRegisters(rdr.holding.start+regsRead, toRead)
			if err != nil {
				return err
			}
			rdr.holding.updateBytes(regsRead, results)
			regsRead += toRead
			if regsRead >= rdr.holding.qty {
				break
			}
		}
	}
	totalRead += regsRead
	if totalRead == 0 {
		return fmt.Errorf("Read no data. Do you need to configure registers?")
	}
	return nil
}

// Units For the given register code, return the units specified
func (rdr *Reader) Units(code int) string {
	reg, ck := rdr.registers[code]
	if !ck {
		return ""
	}
	return reg.Units
}

// Get Return the data stored following a Read() call.
func (rdr *Reader) Get(code int, factored bool) (rValue Value, err error) {
	reg, ck := rdr.registers[code]
	if !ck {
		err = fmt.Errorf("Code %d was not registered", code)
		return
	}

	switch getRegisterType(code) {
	case 3:
		rValue = rdr.input.getValue(reg)
	case 4:
		rValue = rdr.holding.getValue(reg)
	}
	if factored {
		reg.applyFactor(&rValue)
	}
	return
}

// Map Return a map object of the registers. If getting a register returns a value it is
// simply omitted from the map.
func (rdr *Reader) Map(factored bool) map[int]Value {
	mapValues := make(map[int]Value)
	if err := rdr.Read(); err != nil {
		log.Printf("Error reading values: %s", err)
		return mapValues
	}

	for code, reg := range rdr.registers {
		var val Value
		switch getRegisterType(code) {
		case 3:
			val = rdr.input.getValue(reg)
		case 4:
			val = rdr.holding.getValue(reg)
		}
		if factored {
			reg.applyFactor(&val)
		}
		//		log.Printf("reg: %s => %f", reg.description, val.Ieee32)
		mapValues[code] = val
	}
	return mapValues
}

// Dump Query all defined registers and print the results to stdout.
func (rdr *Reader) Dump(factored bool) {
	if err := rdr.Read(); err != nil {
		fmt.Printf("Unable to read register data from device.\n%s\n", err)
		return
	}

	var keys []int
	for k := range rdr.registers {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, code := range keys {
		reg := rdr.registers[code]

		var val Value
		switch getRegisterType(code) {
		case 3:
			val = rdr.input.getValue(reg)
		case 4:
			val = rdr.holding.getValue(reg)
		}

		if factored {
			reg.applyFactor(&val)
			fmt.Printf(baseFmt+ieeeFmt+" %s\n", code, reg.Description, val.Ieee32, reg.Units)
			continue
		}

		switch reg.Format {
		case "u16":
			fmt.Printf(baseFmt+numFmt+" %s\n", code, reg.Description, val.Unsigned16, reg.Units)
		case "s16":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.Description, val.Signed16, reg.Units)
		case "u32":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.Description, val.Unsigned32, reg.Units)
		case "s32":
			fmt.Printf(baseFmt+numFmt+"%s\n", code, reg.Description, val.Signed32, reg.Units)
		case "ieee32":
			fmt.Printf(baseFmt+ieeeFmt+"%s\n", code, reg.Description, val.Ieee32, reg.Units)
		case "coil":
			fmt.Printf(baseFmt+"%t\n", code, reg.Description, val.Coil)
		}
	}
}

// ScanHolding Given a start and stop register, scan the holding registers. Added as a
// convenience.
func (rdr *Reader) ScanHolding(start, stop uint16) {
	qty := stop - start + 1
	results, err := rdr.client.ReadHoldingRegisters(start, qty)
	if err != nil {
		fmt.Printf("Unable to read registers %d to %d\n%s\n", start, stop, err)
		return
	}
	for n := uint16(0); n < qty; n++ {
		reg := start + n
		val := uint16(results[n*2])<<8 + uint16(results[n*2+1])
		fmt.Printf("Register %d [%04X] : %X [%d]\n", reg, reg, val, val)
	}
}

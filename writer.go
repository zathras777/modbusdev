package modbusdev

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/goburrow/modbus"
)

// Reader A reader structure allows us to tie a client to a device register map
type Writer struct {
	client    modbus.Client
	registers map[int]Register
}

// NewWriter Return a configured Writer with the correct register mappings.
// Device names are converted to lower case for matching, so case provided is irrelevant.
func NewWriter(client modbus.Client, device string) (wrt Writer, err error) {
	wrt.client = client
	switch strings.ToLower(device) {
	case "sdm230":
		wrt.addRegisters(sdm230)
	case "solaxx1hybrid":
		wrt.addRegisters(solaxX1Hybrid)
	default:
		err = fmt.Errorf("Device '%s' is not known. Add the details and then update reader.go & writer.go to include it", device)
	}
	return
}

func (wrt *Writer) addRegisters(possible map[int]Register) {
	wrt.registers = make(map[int]Register, len(possible))
	for num, reg := range possible {
		if getRegisterType(num) == 4 {
			wrt.registers[num] = reg
		}
	}
}

// WriteSimple Write a given int value to a register after converting the type (if possible)
func (wrt *Writer) WriteSimple(code, value int) error {
	reg, ck := wrt.registers[code]
	if !ck {
		return fmt.Errorf("Register %d unknown", code)
	}
	switch reg.format {
	case "u16", "s16", "u32", "s32":
		bytes, err := formatIntAsBytes(reg.format, value)
		if err != nil {
			return err
		}
		return wrt.writeSingle(reg, bytes)
	default:
		return fmt.Errorf("Cannot convert int to %s", reg.format)
	}
}

// WriteRegister Write a given value to a register
func (wrt *Writer) WriteRegister(code int, val Value) error {
	reg, ck := wrt.registers[code]
	if !ck {
		return fmt.Errorf("Register %d unknown", code)
	}
	return wrt.writeSingle(reg, val.asBytes(reg.format))
}

// WriteDirect Write the given values to the specified register
func (wrt *Writer) WriteDirect(address, value uint16) error {
	rrr, err := wrt.client.WriteSingleRegister(address, value)
	//	rrr, err := wrt.client.WriteMultipleRegisters(address, uint16(len(byts)), byts)
	//	if !bytes.Equal(rrr, byts) {
	//		fmt.Printf("Did not get expected return value: %v != %v\n", rrr, byts)
	//	}
	fmt.Printf("rrr = %v\n", rrr)
	return err
}

func (wrt *Writer) writeSingle(reg Register, byts []byte) error {
	switch reg.format {
	case "u16", "s16":
		uval := uint16(byts[0])<<8 + uint16(byts[1])
		rrr, err := wrt.client.WriteSingleRegister(reg.register, uval)
		if err != nil {
			return err
		}
		fmt.Printf("Return: %v vs %v\n", rrr, byts)
		if !bytes.Equal(rrr, byts) {
			fmt.Printf("WriteSingle did not return identical values. %v != %v\n", rrr, byts)
			return fmt.Errorf("Incorrect return from write. %v != %v", rrr, byts)
		}
	}
	return nil
}

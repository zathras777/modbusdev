package modbusdev

import (
	"fmt"
	"testing"
)

func TestBytesUnsigned16(t *testing.T) {
	var val Value
	val.Unsigned16 = 44609
	ck := val.asBytes("u16")
	expectedValue := []byte{0xAE, 0x41}
	fmt.Printf("%X\n", ck)
	if ck[0] != expectedValue[0] || ck[1] != expectedValue[1] {
		t.Fatalf("Incorrect value. Got %X expected %X", ck, expectedValue)
	}
}

func TestBytesSigned16(t *testing.T) {
	var val Value
	val.Signed16 = -20927
	ck := val.asBytes("s16")
	expectedValue := []byte{0xAE, 0x41}
	fmt.Printf("%X\n", ck)
	if ck[0] != expectedValue[0] || ck[1] != expectedValue[1] {
		t.Fatalf("Incorrect value. Got %X expected %X", ck, expectedValue)
	}
}

func TestByteUnsigned32(t *testing.T) {
	var val Value
	val.Unsigned32 = 2923517522
	ck := val.asBytes("u32")
	expectedValue := []byte{0xAE, 0x41, 0x56, 0x52}
	fmt.Printf("%X\n", ck)
	if ck[0] != expectedValue[0] || ck[1] != expectedValue[1] || ck[2] != expectedValue[2] || ck[3] != expectedValue[3] {
		t.Fatalf("Incorrect value. Got %X expected %X", ck, expectedValue)
	}
}

func TestByteSigned32(t *testing.T) {
	var val Value
	val.Signed32 = -1371449774
	ck := val.asBytes("s32")
	expectedValue := []byte{0xAE, 0x41, 0x56, 0x52}
	fmt.Printf("%X\n", ck)
	if ck[0] != expectedValue[0] || ck[1] != expectedValue[1] || ck[2] != expectedValue[2] || ck[3] != expectedValue[3] {
		t.Fatalf("Incorrect value. Got %X expected %X", ck, expectedValue)
	}
}

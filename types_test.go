package modbusdev

import (
	"fmt"
	"testing"
)

func TestUnsigned16(t *testing.T) {
	testVals := []byte{0xAE, 0x41}
	if v := unsigned16(testVals); v != 44609 {
		t.Fatalf("Incorrect value. Got %d expected 44609", v)
	}
}

func TestSigned16(t *testing.T) {
	testVals := []byte{0xAE, 0x41}
	if v := signed16(testVals); v != -20927 {
		t.Fatalf("Incorrect value. Got %d expected -20927", v)
	}

}

func TestUnsigned32(t *testing.T) {
	testVals := []byte{0xAE, 0x41, 0x56, 0x52}
	if v := unsigned32(testVals); v != 2923517522 {
		t.Fatalf("Incorrect value. Got %d expected 2,923,517,522", v)
	}
}

func TestSigned32(t *testing.T) {
	testVals := []byte{0xAE, 0x41, 0x56, 0x52}
	if v := signed32(testVals); v != -1371449774 {
		t.Fatalf("Incorrect value. Got %d expected -1,371,449,774", v)
	}
}

func TestBool16(t *testing.T) {
	testVals := []byte{0x00, 0x01}
	if v := bool16(testVals); !v {
		t.Fatalf("Incorrect value. Got %t expected true", v)
	}
}

func TestIeee32(t *testing.T) {
	testVals := []byte{0x40, 0x49, 0x0f, 0xdb}
	v := ieee32(testVals)
	if vs := fmt.Sprintf("%.6f", v); vs != "3.141593" {
		t.Fatalf("Incorrect value. Got %s expected 3.141593", vs)
	}
}

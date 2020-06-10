package modbusdev

import (
	"testing"
)

func TestRegister_1(t *testing.T) {
	r := Register{"Test", "", 1, "u16", 1}
	if r.registersRqd() != 1 {
		t.Fatalf("Incorrect registersRqd() value, %d vs expected 1", r.registersRqd())
	}
	if r.maxRegister() != 2 {
		t.Fatalf("Invalid maxRegister() of %d vs expected 1", r.maxRegister())
	}
	r = Register{"Test", "", 1, "u32", 1}
	if r.registersRqd() != 2 {
		t.Fatalf("Incorrect registersRqd() value, %d vs expected 1", r.registersRqd())
	}
	if r.maxRegister() != 3 {
		t.Fatalf("Invalid maxRegister() of %d vs expected 1", r.maxRegister())
	}
}

func TestRegisterCache(t *testing.T) {
	r := Register{"Test", "", 1, "u16", 1}
	rc := registerCache{}
	rc.init()
	rc.update(r)
	if rc.start != 1 {
		t.Fatalf("Invalid start point in registerCache, %d vs expected 1", rc.start)
	}
	if rc.qty != 1 {
		t.Fatalf("Invalid register quantity in registerCache, %d vs expected 1", rc.qty)
	}
	r = Register{"Test", "", 5, "u32", 1}
	rc.update(r)
	if rc.start != 1 {
		t.Fatalf("Invalid start point in registerCache, %d vs expected 1", rc.start)
	}
	if rc.qty != 6 {
		t.Fatalf("Invalid register quantity in registerCache, %d vs expected 1", rc.qty)
	}

}

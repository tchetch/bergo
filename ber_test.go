package ber

import "testing"

func TestNull(t *testing.T) {
	var bval BerVal

	bval.Null()

	if 0x05 != bval[0] && 0x00 != bval[1] {
		t.Fail()
	}
}

func TestBool(t *testing.T) {
	var b BerVal

	b.Bool(true)
	if 0x01 != b[0] && 0xFF != b[1] {
		t.Fail()
	}
	b.Bool(false)
	if 0x01 != b[0] && 0x00 != b[1] {
		t.Fail()
	}
}

func TestEOC(t *testing.T) {
	var b BerVal

	b.EOC()

	if 0x00 != b[0] && 0x00 != b[1] {
		t.Fail()
	}
}

func TestSet(t *testing.T) {
	var b BerVal

	b.Set()

	if 0x31 != b[0] && 0x80 != b[1] {
		t.Fail()
	}
}

func TestSequence(t *testing.T) {
	var b BerVal

	b.Sequence()

	if 0x30 != b[0] && 0x80 != b[1] {
		t.Fail()
	}
}

func TestInteger64ToBer(t *testing.T) {
	var res []byte
	min := []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	max := []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	zero := []byte{0x00}
	minus1 := []byte{0xFF}
	one := []byte{0x01}
	vr1 := []byte{0xFF, 0x04}                   // -252
	vr2 := []byte{0xFD, 0x6D, 0x83, 0x7E}       // -43154562
	vr3 := []byte{0x00, 0xA5, 0x65, 0xA2, 0x76} // 2774901366

	res = Integer64ToBer(-9223372036854775808)
	if len(res) != len(min) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != min[i] {
			t.Fail()
		}
	}
	res = Integer64ToBer(9223372036854775807)
	if len(res) != len(max) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != max[i] {
			t.Fail()
		}
	}
	res = Integer64ToBer(0)
	if len(res) != len(zero) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != zero[i] {
			t.Fail()
		}
	}
	res = Integer64ToBer(-1)
	if len(res) != len(minus1) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != minus1[i] {
			t.Fail()
		}
	}
	res = Integer64ToBer(1)
	if len(res) != len(one) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != one[i] {
			t.Fail()
		}
	}
	res = Integer64ToBer(-252)
	if len(res) != len(vr1) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != vr1[i] {
			t.Fail()
		}
	}
	res = Integer64ToBer(-43154562)
	if len(res) != len(vr2) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != vr2[i] {
			t.Fail()
		}
	}
	res = Integer64ToBer(2774901366)
	if len(res) != len(vr3) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != vr3[i] {
			t.Fail()
		}
	}
}

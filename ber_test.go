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

func TestMakeBERLen(t *testing.T) {
	var res []byte
	max := []byte{0x88, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} // 18446744073709551615
	min := []byte{0x00}
	ord := []byte{0x88, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08} // 72623859790382856
	x609 := []byte{0x81, 0xC9}                                          // example value in X.609 document -> 201
	res = MakeBERLen(18446744073709551615)
	if len(res) != len(max) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != max[i] {
			t.Fail()
		}
	}
	res = MakeBERLen(0)
	if len(res) != len(min) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != min[i] {
			t.Fail()
		}
	}
	res = MakeBERLen(72623859790382856)
	if len(res) != len(ord) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != ord[i] {
			t.Fail()
		}
	}
	res = MakeBERLen(201)
	if len(res) != len(x609) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != x609[i] {
			t.Fail()
		}
	}
}

func TestOctetstringC(t *testing.T) {
	var b BerVal

	b.OctetstringC()

	if 0x24 != b[0] && 0x80 != b[1] {
		t.Fail()
	}
}

func TestOctetstringP(t *testing.T) {
	var b BerVal
	val := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F}

	b.OctetstringP(val)

	if 0x04 != b[0] && 0x81 != b[1] && 0xC0 != b[2] {
		t.Fail()
	}

	for i := 0; i < len(val); i++ {
		if b[i+3] != val[i] {
			t.Fail()
		}
	}
}

func TestMakeBigTag(t *testing.T) {
	var res []byte
	v1 := []byte{0x1F, 0x36}                   // 54
	v2 := []byte{0x1F, 0x81, 0x1A}             // 154
	v3 := []byte{0x1F, 0x81, 0xB7, 0x8D, 0x40} //3000000

	res = MakeBigTag(54)
	if len(res) != len(v1) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != v1[i] {
			t.Fail()
		}
	}
	res = MakeBigTag(154)
	if len(res) != len(v2) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != v2[i] {
			t.Fail()
		}
	}
	res = MakeBigTag(3000000)
	if len(res) != len(v3) {
		t.Fail()
	}
	for i := 0; i < len(res); i++ {
		if res[i] != v3[i] {
			t.Fail()
		}
	}
}

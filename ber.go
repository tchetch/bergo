package ber

type BerVal []byte
type Ber struct {
	value       BerVal
	constructed bool
	tag         uint64
	class       byte
	content     []byte
}

/* Class bits (to be OR'd) */
const (
	Universal   byte = 0x00
	Application byte = 0x40
	Context     byte = 0x80
	Private     byte = 0xC0
)

/* P/C (to be OR'd) */
const (
	Primitive   byte = 0x00
	Constructed byte = 0x20
)

/* Some Universal Class Tag */
const (
	EOC         byte = 0x00
	Integer     byte = 0x02
	Bitstring   byte = 0x03
	Octetstring byte = 0x04
	Bool        byte = 0x01
	Null        byte = 0x05
	OID         byte = 0x06
	Real        byte = 0x09
	Enumerated  byte = 0x0A
	Sequence    byte = 0x10
	Set         byte = 0x11
)

func (bval *BerVal) Bool(v bool) {
	b := []byte{0x00}
	if v {
		b[0] = 0xFF
	}
	bval.Any(Universal, uint64(Bool), false, b)
}

func (bval *BerVal) Null() {
	bval.Any(Universal, uint64(Null), false, nil)
}

func (bval *BerVal) EOC() {
	bval.Any(Universal, uint64(EOC), false, nil)
}

func (bval *BerVal) Set() {
	bval.Any(Universal, uint64(Set), true, nil)
}

func (bval *BerVal) Sequence() {
	bval.Any(Universal, uint64(Sequence), true, nil)
}

/* Go store int64 as 2's complement exactly as needed for BER */
func Integer64ToBer(v int64) (b []byte) {
	b = make([]byte, 8)

	for j := 0; j < 8; j++ {
		b[j] = byte(v >> uint((7-j)*8))
	}

	/* Remove forbidden leading octets */
	skip := 0
	for j := 0; j < 7; j++ {
		if (b[j] == 0x00) && !(b[j+1]&0x80 == 0x80) {
			skip++
		} else if (b[j] == 0xFF) && (b[j+1]&0x80 == 0x80) {
			skip++
		} else {
			break
		}
	}

	return b[skip:8]
}

func (bval *BerVal) Integer64(v int64) {
	bval.Any(Universal, uint64(Integer), false, Integer64ToBer(v))
}

func (bval *BerVal) Enumerated64(v int64) {
	bval.Any(Universal, uint64(Enumerated), false, Integer64ToBer(v))
}

func MakeBERLen(v uint64) (b []byte) {
	if v <= 127 {
		b = append(b, byte(v))
	} else {
		ocount := 0

		b = append(b, 0x00) // place for length
		for i := 7; i >= 0; i-- {
			if byte(v>>uint(i*8)) != 0x00 {
				b = append(b, byte(v>>uint(i*8)))
				ocount++
			}
		}

		b[0] = 0x80 | (byte(ocount) & 0x3F)
	}

	return b
}

/* Primitive octetstring, if size is known in advance 
 * 0 sized octetstring are allowed.
 */
func (bval *BerVal) OctetstringP(v []byte) {
	bval.Any(Universal, uint64(Octetstring), false, v)
}

/* If size is not known in advance, use OctetstringC then use OctetstringP to
 * send segment as they come. A segment can be zero sized.
 * X.609 allow to send constructed octetstring within a constructed octetstring
 * so you may restart a constructed octetstring ... but you should prefer
 * primitive octetstring.
 * Finish with EOC.
 */
func (bval *BerVal) OctetstringC() {
	bval.Any(Universal, uint64(Octetstring), true, nil)
}

func MakeBigTag(tag uint64) (b []byte) {
	var tmp []byte
	for tag != 0 {
		tmp = append(tmp, (byte(tag) & 0x7F))
		tag = tag >> 7
	}

	leading := true
	b = append(b, 0x1F) // bit 5 to 1 set to 1
	for i := len(tmp) - 1; i >= 0; i-- {
		if tmp[i] != 0x00 {
			leading = false
		}
		if !leading {
			b = append(b, 0x80|tmp[i])
		}
	}
	b[len(b)-1] = b[len(b)-1] & 0x7F // last byte 8th bit to 0

	return b
}

func (bval *BerVal) Any(class byte, tag uint64, constructed bool, val []byte) {
	var tmp []byte
	if tag > 30 {
		tmp = MakeBigTag(tag)
	} else {
		tmp = append(tmp, byte(tag)&0x1F)
	}

	if constructed {
		tmp[0] |= 0x20
	}
	tmp[0] |= byte((class & 0x03) << 6)

	*bval = tmp

	if len(val) > 0 && !constructed {
		tmp = MakeBERLen(uint64(len(val)))
		for i := 0; i < len(val); i++ {
			tmp = append(tmp, val[i])
		}
	} else {
		tmp = append(tmp, 0x80)
	}

	for i := 0; i < len(tmp); i++ {
		*bval = append(*bval, tmp[i])
	}
}

func (ber *Ber) Constructed() {
	ber.constructed = true
}

func (ber *Ber) PutInt64(v int64) {
	copy(ber.content, Integer64ToBer(v))
}

func (ber *Ber) Build() {
	ber.value.Any(ber.class, ber.tag, ber.constructed, ber.content)
}

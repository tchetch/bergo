package ber

type BerVal []byte

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
	var b byte
	if v {
		b = 0xFF
	} else {
		b = 0x00
	}
	*bval = append(*bval, Universal|Primitive|Bool, 0x01, b)
}

func (bval *BerVal) Null() {
	*bval = append(*bval, Universal|Primitive|Null, 0x00)
}

func (bval *BerVal) EOC() {
	*bval = append(*bval, Universal|Primitive|EOC, 0x00)
}

func (bval *BerVal) Set() {
	*bval = append(*bval, Universal|Constructed|Set, 0x80)
}

func (bval *BerVal) Sequence() {
	*bval = append(*bval, Universal|Constructed|Sequence, 0x80)
}

/* Go store int64 as 2's complement exactly as needed for BER */
func Integer64ToBer(v int64) []byte {
	var b [8]byte

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
	var b []byte

	b = Integer64ToBer(v)
	*bval = append(*bval, Universal|Primitive|Integer, byte(len(b)))
	for i := 0; i < len(b); i++ {
		*bval = append(*bval, b[i])
	}
}

func (bval *BerVal) Enumerated64(v int64) {
	var b []byte

	b = Integer64ToBer(v)
	*bval = append(*bval, Universal|Primitive|Enumerated, byte(len(b)))
	for i := 0; i < len(b); i++ {
		*bval = append(*bval, b[i])
	}
}

func MakeBERLen(v uint64) []byte {
	var b []byte

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
	var length []byte

	length = MakeBERLen(uint64(len(v)))

	*bval = append(*bval, Universal|Primitive|Octetstring)
	for i := 0; i < len(length); i++ {
		*bval = append(*bval, length[i])
	}
	for i := 0; i < len(v); i++ {
		*bval = append(*bval, v[i])
	}
}

/* If size is not known in advance, use OctetstringC then use OctetstringP to
 * send segment as they come. A segment can be zero sized.
 * X.609 allow to send constructed octetstring within a constructed octetstring
 * so you may restart a constructed octetstring ... but you should prefer
 * primitive octetstring.
 * Finish with EOC.
 */
func (bval *BerVal) OctetstringC() {
	*bval = append(*bval, Universal|Constructed|Octetstring, 0x80)
}

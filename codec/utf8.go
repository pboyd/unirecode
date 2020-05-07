package codec

import (
	"errors"
	"io"
)

func init() {
	registerCodec("UTF-8", NewUTF8Decoder, NewUTF8Encoder)
}

var _ Decoder = &UTF8Decoder{}

// UTF8Decoder implements Decoder for UTF-8.
type UTF8Decoder struct {
	buf []byte
}

// NewUTF8Decoder creates a new instance of UTF8Decoder
func NewUTF8Decoder() Decoder {
	return &UTF8Decoder{
		buf: make([]byte, 0, 4),
	}
}

// Decode satifies the Decoder interface for UTF-8.
func (d *UTF8Decoder) Decode(br io.ByteReader) (rune, error) {
	b, err := br.ReadByte()
	if err != nil {
		return 0, err
	}

	l := utf8Len(b)
	if l <= 0 || l > 4 {
		return 0, errors.New("invalid character")
	}
	if l == 1 {
		return rune(b), nil
	}

	// 2 byte characters have three bits for the length, so use the 5 low bits.
	// 3 byte character, use the low 4 bits.
	// 4 byte character, use the low 3 bits.
	r := rune(b) & (0x7f >> l)

	for i := 1; i < l; i++ {
		b, err = br.ReadByte()
		if err != nil {
			return 0, err
		}

		// Make sure the two high bits are 1 and 0 respectively.
		if b>>6 != 2 {
			return 0, errors.New("invalid character")
		}

		// There are 6 bits of the code point in this byte. So shift the
		// everything we've gotten so far to make room.
		r <<= 6

		// Put the low 6 bits of b on the low 6 bits of r.
		r |= rune(b & (0xff >> 2))
	}

	return r, nil
}

// utf8Len returns the number total bytes used for the UTF-8 character based on the information in the first byte.
func utf8Len(b byte) int {
	// Under 128 is ASCII
	if b < 0x80 {
		return 1
	}

	// 0b10xxxxxx is invalid
	// 0b110xxxxx is 2 bytes
	// 0b1110xxxx is 3 bytes
	// 0b11110xxx is 4 bytes
	for i := 1; i <= 4; i++ {
		var m byte = 0x80 >> i
		if b&m == 0 {
			return i
		}
	}

	return 0
}

var _ Encoder = &UTF8Encoder{}

// UTF8Encoder implements Encoder for UTF-8.
type UTF8Encoder struct {
}

// NewUTF8Encoder creates a new instance of UTF8Encoder
func NewUTF8Encoder() Encoder {
	return &UTF8Encoder{}
}

// Encode satifies the Decoder interface for UTF-8.
func (*UTF8Encoder) Encode(w io.Writer, r rune) error {
	buf := make([]byte, 0, 4)
	switch {
	case r < 0:
		return errors.New("invalid character")
	case r < 0x80:
		buf = append(buf, byte(r))
	case r < 0x800:
		// 11 bits available, 5 bits in the first byte
		buf = buf[:2]
		buf[0] = 0x80 | 0x40 | byte(r>>6)
		buf[1] = 0x80 | byte(r&0x3f)
	case r < 0xffff:
		// 16 bits available, 4 in the first byte
		buf = buf[:3]
		buf[0] = 0x80 | 0x40 | 0x20 | byte(r>>12)
		buf[1] = 0x80 | byte(r>>6&0x3f)
		buf[2] = 0x80 | byte(r&0x3f)
	case r < 0x10ffff:
		// 21 bits available, 3 in the first byte
		buf = buf[:4]
		buf[0] = 0x80 | 0x40 | 0x20 | 0x10 | byte(r>>18)
		buf[1] = 0x80 | byte(r>>12&0x3f)
		buf[2] = 0x80 | byte(r>>6&0x3f)
		buf[3] = 0x80 | byte(r&0x3f)
	default:
		return errors.New("invalid character")
	}

	_, err := w.Write(buf)
	return err
}

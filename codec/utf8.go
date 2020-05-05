package codec

import (
	"errors"
	"io"
)

func init() {
	registerDecoder("UTF-8", NewUTF8Decoder)
	registerEncoder("UTF-8", NewUTF8Encoder)
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
func (d *UTF8Decoder) Decode(r io.ByteReader) ([]byte, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	l := utf8Len(b)
	if l <= 0 {
		return nil, errors.New("invalid character")
	}

	d.buf = d.buf[:l]
	d.buf[0] = b
	if l == 1 {
		return d.buf, nil
	}

	for i := 1; i < l; i++ {
		d.buf[i], err = r.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	return d.buf, nil
}

// utf8Len returns the number total bytes used for the UTF-8 character based on the information in the first byte.
func utf8Len(b byte) int {
	// Under 128 is ASCII
	if b < 128 {
		return 1
	}

	// 0b10xxxxxx is invalid
	// 0b110xxxxx is 2 bytes
	// 0b1110xxxx is 3 bytes
	// 0b11110xxx is 4 bytes
	for i := 1; i <= 4; i++ {
		var m byte = 128 >> i
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
func (*UTF8Encoder) Encode(w io.Writer, buf []byte) error {
	_, err := w.Write(buf)
	return err
}

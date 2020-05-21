package codec

import (
	"errors"
	"io"
)

func init() {
	registerCodec("ASCII", NewASCIIDecoder, NewASCIIEncoder)
}

var _ Decoder = &ASCIIDecoder{}

// ASCIIDecoder implements Decoder for ASCII.
type ASCIIDecoder struct {
}

// NewASCIIDecoder creates a new instance of ASCIIDecoder
func NewASCIIDecoder() Decoder {
	return &ASCIIDecoder{}
}

// Decode satifies the Decoder interface for ASCII.
func (d *ASCIIDecoder) Decode(r io.Reader) (rune, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	return rune(buf[0]), nil
}

var _ Encoder = &ASCIIEncoder{}

// ASCIIEncoder implements Encoder for ASCII.
type ASCIIEncoder struct {
}

// NewASCIIEncoder creates a new instance of ASCIIEncoder
func NewASCIIEncoder() Encoder {
	return &ASCIIEncoder{}
}

// Encode satifies the Decoder interface for ASCII.
func (*ASCIIEncoder) Encode(w io.Writer, r rune) error {
	if r > 127 {
		return errors.New("character out of range")
	}

	buf := make([]byte, 1)
	buf[0] = byte(r)
	_, err := w.Write(buf)
	return err
}

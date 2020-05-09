package codec

import (
	"errors"
	"io"
)

func init() {
	registerCodec("UCS-2", NewUCS2Decoder, NewUCS2Encoder)
	registerCodec("UCS-2BE", NewUCS2BEDecoder, NewUCS2BEEncoder)
	registerCodec("UCS-2LE", NewUCS2LEDecoder, NewUCS2LEEncoder)
}

// UCS2Decoder reads UCS-2 characters. UCS-2 is a character encoding where each
// code point takes exactly 2 bytes. It can only encode characters up to
// U+FFFF.
type UCS2Decoder struct {
	byteOrder byteOrder
}

// NewUCS2Decoder returns a UCS-2 decoder.
//
// If the encoded text begins with a byte order mark (U+FEFF) that will
// determine the endianness used. Otherwise, it defaults to little-endian.
func NewUCS2Decoder() Decoder {
	return &UCS2Decoder{}
}

// NewUCS2LEDecoder returns a UCS-2 decoder with a little-endian byte order.
func NewUCS2LEDecoder() Decoder {
	return &UCS2Decoder{
		byteOrder: littleEndian,
	}
}

// NewUCS2BEDecoder returns a UCS-2 decoder with a big-endian byte order.
func NewUCS2BEDecoder() Decoder {
	return &UCS2Decoder{
		byteOrder: bigEndian,
	}
}

// Decode satifies the Decoder interface for UCS-2.
func (d *UCS2Decoder) Decode(r io.Reader) (rune, error) {
	var buf = make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}

	if d.byteOrder == unknownByteOrder {
		if buf[0] == 0xfe && buf[1] == 0xff {
			d.byteOrder = bigEndian
			_, err = io.ReadFull(r, buf)
		} else if buf[0] == 0xff && buf[1] == 0xfe {
			d.byteOrder = littleEndian
			_, err = io.ReadFull(r, buf)
		} else {
			d.byteOrder = littleEndian
		}

		if err != nil {
			return 0, err
		}
	}

	switch d.byteOrder {
	case bigEndian:
		return (rune(buf[0]) << 8) | rune(buf[1]), nil
	case littleEndian:
		return (rune(buf[1]) << 8) | rune(buf[0]), nil
	default:
		return 0, errors.New("unknown byte order")
	}
}

// UCS2Encoder encodes unicode code points using exactly two bytes.
// It can only encode characters up to U+FFFF.
type UCS2Encoder struct {
	byteOrder byteOrder
	writeBOM  bool
}

// NewUCS2Encoder returns a UCS-2 encoder with a little-endian byte order.
//
// This is identical to NewUCS2LEEncoder except this one does not write a byte
// order mark.
func NewUCS2Encoder() Encoder {
	return &UCS2Encoder{
		byteOrder: littleEndian,
		writeBOM:  false,
	}
}

// NewUCS2LEEncoder returns a UCS-2 encoder with a little-endian byte order.
//
// It will write a byte order mark with the first character.
func NewUCS2LEEncoder() Encoder {
	return &UCS2Encoder{
		byteOrder: littleEndian,
		writeBOM:  true,
	}
}

// NewUCS2BEEncoder returns a UCS-2 encoder with a big-endian byte order.
//
// It will write a byte order mark with the first character.
func NewUCS2BEEncoder() Encoder {
	return &UCS2Encoder{
		byteOrder: bigEndian,
		writeBOM:  true,
	}
}

// Encode satifies the Encoder interface for UCS-2.
func (d *UCS2Encoder) Encode(w io.Writer, r rune) error {
	if r > 0xffff {
		return errors.New("character out of range")
	}

	buf := make([]byte, 2)

	if d.writeBOM {
		d.writeBOM = false
		if d.byteOrder == bigEndian {
			buf[0] = 0xfe
			buf[1] = 0xff
		} else {
			buf[1] = 0xfe
			buf[0] = 0xff
		}
		_, err := w.Write(buf)
		if err != nil {
			return err
		}
	}

	if d.byteOrder == bigEndian {
		buf[0] = byte(r >> 8)
		buf[1] = byte(r)
	} else {
		buf[0] = byte(r)
		buf[1] = byte(r >> 8)
	}
	_, err := w.Write(buf)
	return err
}

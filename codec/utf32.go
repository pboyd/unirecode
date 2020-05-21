package codec

import (
	"errors"
	"io"
)

func init() {
	registerCodec("UTF-32", NewUTF32Decoder, NewUTF32Encoder)
	registerCodec("UTF-32BE", NewUTF32BEDecoder, NewUTF32BEEncoder)
	registerCodec("UTF-32LE", NewUTF32LEDecoder, NewUTF32LEEncoder)

	registerCodec("UCS-4", NewUTF32Decoder, NewUTF32Encoder)
	registerCodec("UCS-4BE", NewUTF32BEDecoder, NewUTF32BEEncoder)
	registerCodec("UCS-4LE", NewUTF32LEDecoder, NewUTF32LEEncoder)
}

// UTF32Decoder reads UTF-32 encoded Unicode characters. UTF-32 is a character
// encoding where each code point takes exactly 4 bytes.
type UTF32Decoder struct {
	byteOrder byteOrder
}

// NewUTF32Decoder returns a UTF-32 decoder.
//
// If the encoded text begins with a byte order mark (U+FEFF) that will
// determine the endianness used. Otherwise it will derive the byte order by
// looking for the 0-byte in the first code point.
func NewUTF32Decoder() Decoder {
	return &UTF32Decoder{}
}

// NewUTF32LEDecoder returns a UTF-32 decoder with a little-endian byte order.
func NewUTF32LEDecoder() Decoder {
	return &UTF32Decoder{
		byteOrder: littleEndian,
	}
}

// NewUTF32BEDecoder returns a UTF-32 decoder with a big-endian byte order.
func NewUTF32BEDecoder() Decoder {
	return &UTF32Decoder{
		byteOrder: bigEndian,
	}
}

// Decode satifies the Decoder interface for UTF-32.
func (d *UTF32Decoder) Decode(r io.Reader) (rune, error) {
	var buf = make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}

	if d.byteOrder == unknownByteOrder {
		// Check for a BE or LE byte order mark.
		//
		// UTF-32 doesn't really need a BOM, because the most
		// significant byte is always zero (max assigned Unicode code
		// point is 0x10ffff which never takes four bytes), so check
		// which end is zero if there's no BOM.
		if buf[0] == 0xfe && buf[1] == 0xff && buf[2] == 0 && buf[3] == 0 {
			d.byteOrder = bigEndian
			_, err = io.ReadFull(r, buf)
		} else if buf[0] == 0xff && buf[1] == 0xfe && buf[2] == 0 && buf[3] == 0 {
			d.byteOrder = littleEndian
			_, err = io.ReadFull(r, buf)
		} else if buf[0] == 0 {
			d.byteOrder = bigEndian
		} else if buf[3] == 0 {
			d.byteOrder = littleEndian
		} else {
			return 0, errors.New("invalid UTF-32 character")
		}

		if err != nil {
			return 0, err
		}
	}

	switch d.byteOrder {
	case bigEndian:
		return (rune(buf[0]) << 24) | (rune(buf[1]) << 16) | (rune(buf[2]) << 8) | (rune(buf[3])), nil
	case littleEndian:
		return (rune(buf[3]) << 24) | (rune(buf[2]) << 16) | (rune(buf[1]) << 8) | (rune(buf[0])), nil
	default:
		return 0, errors.New("unknown byte order")
	}
}

// UTF32Encoder encodes unicode code points using exactly four bytes.
type UTF32Encoder struct {
	byteOrder byteOrder
}

// NewUTF32Encoder returns a UTF-32 encoder with a little-endian byte order.
//
// This is identical to NewUTF32LEEncoder.
func NewUTF32Encoder() Encoder {
	return &UTF32Encoder{
		byteOrder: littleEndian,
	}
}

// NewUTF32LEEncoder returns a UTF-32 encoder with a little-endian byte order.
func NewUTF32LEEncoder() Encoder {
	return &UTF32Encoder{
		byteOrder: littleEndian,
	}
}

// NewUTF32BEEncoder returns a UTF-32 encoder with a big-endian byte order.
func NewUTF32BEEncoder() Encoder {
	return &UTF32Encoder{
		byteOrder: bigEndian,
	}
}

// Encode satifies the Encoder interface for UTF-32.
func (d *UTF32Encoder) Encode(w io.Writer, r rune) error {
	buf := make([]byte, 4)

	if d.byteOrder == bigEndian {
		buf[0] = byte(r >> 24)
		buf[1] = byte(r >> 16)
		buf[2] = byte(r >> 8)
		buf[3] = byte(r)
	} else {
		buf[0] = byte(r)
		buf[1] = byte(r >> 8)
		buf[2] = byte(r >> 16)
		buf[3] = byte(r >> 24)
	}
	_, err := w.Write(buf)
	return err
}

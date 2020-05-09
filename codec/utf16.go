package codec

import (
	"errors"
	"io"
)

func init() {
	registerCodec("UTF-16", NewUTF16Decoder, NewUTF16Encoder)
	registerCodec("UTF-16BE", NewUTF16BEDecoder, NewUTF16BEEncoder)
	registerCodec("UTF-16LE", NewUTF16LEDecoder, NewUTF16LEEncoder)
}

// UTF16Decoder reads UTF-16 characters. UTF-16 is identical to UCS-2 for
// characters U+FFFF and below. Characters above U+FFFF are encoded in two
// 16-bit words, called surrogate pairs.
type UTF16Decoder struct {
	ucs2 Decoder
}

// NewUTF16Decoder returns a UTF-16 decoder.
//
// If the encoded text begins with a byte order mark (U+FEFF) that will
// determine the endianness used. Otherwise, it defaults to little-endian.
func NewUTF16Decoder() Decoder {
	return &UTF16Decoder{
		ucs2: NewUCS2Decoder(),
	}
}

// NewUTF16LEDecoder returns a UTF-16 decoder with a little-endian byte order.
func NewUTF16LEDecoder() Decoder {
	return &UTF16Decoder{
		ucs2: NewUCS2LEDecoder(),
	}
}

// NewUTF16BEDecoder returns a UTF-16 decoder with a big-endian byte order.
func NewUTF16BEDecoder() Decoder {
	return &UTF16Decoder{
		ucs2: NewUCS2BEDecoder(),
	}
}

const (
	// Mask the six high bits of a 16 bit number
	utf16SurrogateMask = 0x3f << 10
	utf16HighSurrogate = 0xd800
	utf16LowSurrogate  = 0xdc00
)

// Decode reads one UTF-16 encoded character from the reader.
func (d *UTF16Decoder) Decode(r io.Reader) (rune, error) {
	w1, err := d.ucs2.Decode(r)
	if err != nil {
		return 0, err
	}

	// The high six bits of UTF-16 surrogate pairs are 0xd800
	if w1&utf16SurrogateMask != utf16HighSurrogate {
		// Character under 0x10000 that's not a surrogate, just return.
		return w1, nil
	}

	w2, err := d.ucs2.Decode(r)
	if err != nil {
		return 0, err
	}

	if w2&utf16SurrogateMask != utf16LowSurrogate {
		return 0, errors.New("invalid UTF-16 surrogate pair")
	}

	u := rune(w1&0x3ff) << 10
	u |= rune(w2 & 0x3ff)
	u |= 0x10000

	return u, nil
}

// UTF16Encoder encodes unicode code points using exactly two bytes.
// It can only encode characters up to U+FFFF.
type UTF16Encoder struct {
	ucs2 Encoder
}

// NewUTF16Encoder returns a UCS-2 encoder with a little-endian byte order.
//
// This is identical to NewUTF16LEEncoder except this one does not write a byte
// order mark.
func NewUTF16Encoder() Encoder {
	return &UTF16Encoder{
		ucs2: NewUCS2Encoder(),
	}
}

// NewUTF16LEEncoder returns a UTF-16 encoder with a little-endian byte order.
//
// It will write a byte order mark with the first character.
func NewUTF16LEEncoder() Encoder {
	return &UTF16Encoder{
		ucs2: NewUCS2LEEncoder(),
	}
}

// NewUTF16BEEncoder returns a UTF-16 encoder with a big-endian byte order.
//
// It will write a byte order mark with the first character.
func NewUTF16BEEncoder() Encoder {
	return &UTF16Encoder{
		ucs2: NewUCS2BEEncoder(),
	}
}

// Encode writes one UTF-16 encoded character to the writer.
func (d *UTF16Encoder) Encode(w io.Writer, r rune) error {
	if r < 0x10000 {
		return d.ucs2.Encode(w, r)
	}

	if r > 0x10ffff {
		return errors.New("character out of range")
	}

	// Split the character into two words. Subtract 0x10000, the largest
	// possible code point will be 20 bits. The first ten bits go in the
	// first word, the second ten bits go in the second word.
	r ^= 0x10000
	r1 := utf16HighSurrogate | (r >> 10)
	r2 := utf16LowSurrogate | (r & 0x3ff)

	err := d.ucs2.Encode(w, r1)
	if err != nil {
		return err
	}

	return d.ucs2.Encode(w, r2)
}

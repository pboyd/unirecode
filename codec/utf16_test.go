package codec

import (
	"bytes"
	"testing"
)

func TestUTF16Decoder(t *testing.T) {
	cases := []struct {
		decoder  Decoder
		in       []byte
		expected string
	}{
		{
			decoder:  NewUTF16LEDecoder(),
			in:       []byte{0x44, 0x00, 0x6f, 0x00, 0x77, 0x00, 0x6e, 0x00},
			expected: "Down",
		},
		{
			decoder:  NewUTF16BEDecoder(),
			in:       []byte{0x00, 0x74, 0x00, 0x68, 0x00, 0x65},
			expected: "the",
		},
		{
			decoder:  NewUTF16Decoder(),
			in:       []byte{0x52, 0x00, 0x61, 0x00, 0x62, 0x00, 0x62, 0x00, 0x69, 0x00, 0x74, 0x00},
			expected: "Rabbit",
		},
		{
			decoder:  NewUTF16Decoder(),
			in:       []byte{0xfe, 0xff, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x6C, 0x00, 0x65},
			expected: "Hole",
		},
		{
			decoder:  NewUTF16LEDecoder(),
			in:       []byte{0x20, 0x22},
			expected: "∠",
		},
		{
			decoder:  NewUTF16BEDecoder(),
			in:       []byte{0x22, 0x20},
			expected: "∠",
		},
		{
			decoder:  NewUTF16LEDecoder(),
			in:       []byte{0x3d, 0xd8, 0x07, 0xdc},
			expected: "🐇",
		},
		{
			decoder:  NewUTF16BEDecoder(),
			in:       []byte{0xd8, 0x3d, 0xdc, 0x07},
			expected: "🐇",
		},
	}

	encoder := NewUTF8Encoder()

	for _, c := range cases {
		actual := &bytes.Buffer{}

		err := Recode(bytes.NewReader(c.in), actual, c.decoder, encoder)
		if err != nil {
			t.Errorf("recode error: %v", err)
			continue
		}

		if actual.String() != c.expected {
			t.Errorf("got %q, want %q", actual.String(), c.expected)
		}
	}
}

func TestUTF16Encoder(t *testing.T) {
	cases := []struct {
		encoder  Encoder
		in       string
		expected []byte
	}{
		{
			encoder:  NewUTF16LEEncoder(),
			in:       "Down",
			expected: []byte{0xff, 0xfe, 0x44, 0x00, 0x6f, 0x00, 0x77, 0x00, 0x6e, 0x00},
		},
		{
			encoder:  NewUTF16BEEncoder(),
			in:       "the",
			expected: []byte{0xfe, 0xff, 0x00, 0x74, 0x00, 0x68, 0x00, 0x65},
		},
		{
			encoder:  NewUTF16LEEncoder(),
			in:       "Rabbit",
			expected: []byte{0xff, 0xfe, 0x52, 0x00, 0x61, 0x00, 0x62, 0x00, 0x62, 0x00, 0x69, 0x00, 0x74, 0x00},
		},
		{
			encoder:  NewUTF16BEEncoder(),
			in:       "Hole",
			expected: []byte{0xfe, 0xff, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x6C, 0x00, 0x65},
		},
		{
			encoder:  NewUTF16LEEncoder(),
			in:       "∠",
			expected: []byte{0xff, 0xfe, 0x20, 0x22},
		},
		{
			encoder:  NewUTF16BEEncoder(),
			in:       "∠",
			expected: []byte{0xfe, 0xff, 0x22, 0x20},
		},
		{
			encoder:  NewUTF16LEEncoder(),
			in:       "🐇",
			expected: []byte{0xff, 0xfe, 0x3d, 0xd8, 0x07, 0xdc},
		},
		{
			encoder:  NewUTF16BEEncoder(),
			in:       "🐇",
			expected: []byte{0xfe, 0xff, 0xd8, 0x3d, 0xdc, 0x07},
		},
	}

	decoder := NewUTF8Decoder()

	for _, c := range cases {
		actual := &bytes.Buffer{}

		err := Recode(bytes.NewReader([]byte(c.in)), actual, decoder, c.encoder)
		if err != nil {
			t.Errorf("recode error: %v", err)
			continue
		}

		if !bytes.Equal(actual.Bytes(), c.expected) {
			t.Errorf("got %v, want %v", actual.Bytes(), c.expected)
		}
	}
}

package codec

import (
	"bytes"
	"testing"
)

func TestUCS2Decoder(t *testing.T) {
	cases := []struct {
		decoder  Decoder
		in       []byte
		expected string
	}{
		{
			decoder:  NewUCS2LEDecoder(),
			in:       []byte{0x44, 0x00, 0x6f, 0x00, 0x77, 0x00, 0x6e, 0x00},
			expected: "Down",
		},
		{
			decoder:  NewUCS2BEDecoder(),
			in:       []byte{0x00, 0x74, 0x00, 0x68, 0x00, 0x65},
			expected: "the",
		},
		{
			decoder:  NewUCS2Decoder(),
			in:       []byte{0x52, 0x00, 0x61, 0x00, 0x62, 0x00, 0x62, 0x00, 0x69, 0x00, 0x74, 0x00},
			expected: "Rabbit",
		},
		{
			decoder:  NewUCS2Decoder(),
			in:       []byte{0xfe, 0xff, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x6C, 0x00, 0x65},
			expected: "Hole",
		},
		{
			decoder:  NewUCS2LEDecoder(),
			in:       []byte{0x20, 0x22},
			expected: "∠",
		},
		{
			decoder:  NewUCS2BEDecoder(),
			in:       []byte{0x22, 0x20},
			expected: "∠",
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

func TestUCS2Encoder(t *testing.T) {
	cases := []struct {
		encoder  Encoder
		in       string
		expected []byte
	}{
		{
			encoder:  NewUCS2LEEncoder(),
			in:       "Down",
			expected: []byte{0xff, 0xfe, 0x44, 0x00, 0x6f, 0x00, 0x77, 0x00, 0x6e, 0x00},
		},
		{
			encoder:  NewUCS2BEEncoder(),
			in:       "the",
			expected: []byte{0xfe, 0xff, 0x00, 0x74, 0x00, 0x68, 0x00, 0x65},
		},
		{
			encoder:  NewUCS2LEEncoder(),
			in:       "Rabbit",
			expected: []byte{0xff, 0xfe, 0x52, 0x00, 0x61, 0x00, 0x62, 0x00, 0x62, 0x00, 0x69, 0x00, 0x74, 0x00},
		},
		{
			encoder:  NewUCS2BEEncoder(),
			in:       "Hole",
			expected: []byte{0xfe, 0xff, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x6C, 0x00, 0x65},
		},
		{
			encoder:  NewUCS2LEEncoder(),
			in:       "∠",
			expected: []byte{0xff, 0xfe, 0x20, 0x22},
		},
		{
			encoder:  NewUCS2BEEncoder(),
			in:       "∠",
			expected: []byte{0xfe, 0xff, 0x22, 0x20},
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

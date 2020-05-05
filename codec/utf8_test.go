package codec

import (
	"bytes"
	"testing"
)

func TestUTF8Decode(t *testing.T) {
	cases := []struct {
		buf []byte
		xp  []byte
	}{
		{
			buf: []byte{'x', 0x0a},
			xp:  []byte{'x'},
		},
		{
			buf: []byte{0xd8, 0x89, 0xa},
			xp:  []byte{0xd8, 0x89},
		},
		{
			buf: []byte{0xd8, 0x89, 0xa},
			xp:  []byte{0xd8, 0x89},
		},
		{
			buf: []byte{0xe2, 0x82, 0xa1, 0xa},
			xp:  []byte{0xe2, 0x82, 0xa1},
		},
		{
			buf: []byte{0xf0, 0x9f, 0x8c, 0x8e, 0x0a},
			xp:  []byte{0xf0, 0x9f, 0x8c, 0x8e},
		},
	}

	decoder := GetDecoder("UTF-8")

	for _, c := range cases {
		r := bytes.NewReader(c.buf)
		actual, err := decoder.Decode(r)
		if err != nil {
			t.Errorf("error: %v", err)
			continue
		}

		if string(actual) != string(c.xp) {
			t.Errorf("got %q, want %q", string(actual), string(c.xp))
		}
	}
}

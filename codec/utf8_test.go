package codec

import (
	"bytes"
	"testing"
)

func TestUTF8Decode(t *testing.T) {
	cases := []struct {
		buf []byte
		xp  rune
	}{
		{
			buf: []byte{'x', 0x0a},
			xp:  'x',
		},
		{
			buf: []byte{0xd8, 0x89, 0xa},
			xp:  '؉',
		},
		{
			buf: []byte{0xe2, 0x82, 0xa1, 0xa},
			xp:  '₡',
		},
		{
			buf: []byte{0xf0, 0x9f, 0x8c, 0x8e, 0x0a},
			xp:  '🌎',
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

func TestUTF8Encoder(t *testing.T) {
	cases := []struct {
		r        rune
		expected []byte
	}{
		{
			r:        'x',
			expected: []byte{'x'},
		},
		{
			r:        '؉',
			expected: []byte{0xd8, 0x89},
		},
		{
			r:        '₡',
			expected: []byte{0xe2, 0x82, 0xa1},
		},
		{
			r:        '🌎',
			expected: []byte{0xf0, 0x9f, 0x8c, 0x8e},
		},
	}

	encoder := GetEncoder("UTF-8")

	for _, c := range cases {
		actual := &bytes.Buffer{}
		err := encoder.Encode(actual, c.r)
		if err != nil {
			t.Errorf("error: %v", err)
			continue
		}

		if actual.String() != string(c.expected) {
			t.Errorf("got %q, want %q", actual.String(), string(c.expected))
		}
	}
}

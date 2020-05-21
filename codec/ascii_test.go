package codec

import (
	"bytes"
	"testing"
)

func TestASCIIDecode(t *testing.T) {
	cases := []struct {
		buf []byte
		xp  rune
	}{
		{
			buf: []byte{'x', 0x0a},
			xp:  'x',
		},
	}

	decoder := GetDecoder("ASCII")

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

func TestASCIIEncoder(t *testing.T) {
	cases := []struct {
		r        rune
		expected []byte
	}{
		{
			r:        'x',
			expected: []byte{'x'},
		},
	}

	encoder := GetEncoder("ASCII")

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

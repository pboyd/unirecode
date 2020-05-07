package codec

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

// Decoder reads a single character and returns the Unicode code point as a rune.
type Decoder interface {
	Decode(io.ByteReader) (rune, error)
}

// Encoder writes an encoded Unicode code point to the writer.
type Encoder interface {
	Encode(io.Writer, rune) error
}

type registeredCodec struct {
	initDecoder func() Decoder
	initEncoder func() Encoder
}

var (
	codecRegistryMu sync.RWMutex
	codecRegistry   = map[string]registeredCodec{}
)

func registerCodec(name string, initDecoder func() Decoder, initEncoder func() Encoder) {
	codecRegistryMu.Lock()
	defer codecRegistryMu.Unlock()

	codecRegistry[name] = registeredCodec{
		initDecoder: initDecoder,
		initEncoder: initEncoder,
	}
}

// GetDecoder looks up a decoder by name. Returns nil if no decoder is found
// with the given name.
func GetDecoder(name string) Decoder {
	codecRegistryMu.RLock()
	defer codecRegistryMu.RUnlock()

	cr, ok := codecRegistry[name]
	if !ok {
		return nil
	}
	return cr.initDecoder()
}

// GetEncoder looks up an encoder by name. Returns nil if no encoder is found
// with the given name.
func GetEncoder(name string) Encoder {
	codecRegistryMu.RLock()
	defer codecRegistryMu.RUnlock()

	cr, ok := codecRegistry[name]
	if !ok {
		return nil
	}
	return cr.initEncoder()
}

// Recode decodes data from the reader with decoder, then writes it back out to
// w with the encoder.
func Recode(r io.Reader, w io.Writer, decoder Decoder, encoder Encoder) error {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = bufio.NewReader(r)
	}

	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
		defer bw.Flush()
	}

	for {
		char, err := decoder.Decode(br)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("error decoding character: %w", err)
			}
			break
		}

		err = encoder.Encode(bw, char)
		if err != nil {
			return fmt.Errorf("error encoding character: %w", err)
		}
	}

	return nil
}

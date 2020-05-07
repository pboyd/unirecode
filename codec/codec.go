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

var (
	decoderRegistryMu sync.RWMutex
	decoderRegistry   map[string]func() Decoder
)

func registerDecoder(name string, init func() Decoder) {
	decoderRegistryMu.Lock()
	defer decoderRegistryMu.Unlock()

	if decoderRegistry == nil {
		decoderRegistry = map[string]func() Decoder{}
	}
	decoderRegistry[name] = init
}

// GetDecoder looks up a decoder by name. Returns nil if no decoder is found
// with the given name.
func GetDecoder(name string) Decoder {
	decoderRegistryMu.RLock()
	defer decoderRegistryMu.RUnlock()

	init, ok := decoderRegistry[name]
	if !ok {
		return nil
	}
	return init()
}

var (
	encoderRegistryMu sync.RWMutex
	encoderRegistry   map[string]func() Encoder
)

func registerEncoder(name string, init func() Encoder) {
	encoderRegistryMu.Lock()
	defer encoderRegistryMu.Unlock()

	if encoderRegistry == nil {
		encoderRegistry = map[string]func() Encoder{}
	}
	encoderRegistry[name] = init
}

// GetEncoder looks up an encoder by name. Returns nil if no encoder is found
// with the given name.
func GetEncoder(name string) Encoder {
	encoderRegistryMu.RLock()
	defer encoderRegistryMu.RUnlock()

	init, ok := encoderRegistry[name]
	if !ok {
		return nil
	}
	return init()
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

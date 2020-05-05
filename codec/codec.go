package codec

import (
	"io"
	"sync"
)

// Decoder reads a single character and converts it to UTF-8.
type Decoder interface {
	Decode(io.ByteReader) ([]byte, error)
}

// Encoder converts a character from UTF-8 to some other format and writes it
// to the io.Writer.
type Encoder interface {
	Encode(io.Writer, []byte) error
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

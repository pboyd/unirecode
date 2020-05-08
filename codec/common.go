package codec

type byteOrder int

const (
	unknownByteOrder byteOrder = iota
	littleEndian
	bigEndian
)

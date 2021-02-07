package maybe

import "github.com/spaolacci/murmur3"

// AsBytes is an abstraction needed by the maybe library component to indicate anytype that can be transformed
// to a series of bytes, so that we can apply a hash function to it.
type AsBytes interface {
	Bytes() []byte
}

// HashFunc is any type that can transform an `AsBytes` value into a uint64
type HashFunc func(AsBytes) uint64

// MurmurFrom constructs a HashFunc from a murmur3 hash function and a multiplication factor
// This is a trick that can be used to create multiple hash function from the same one.
func MurmurFrom(value uint64) HashFunc {
	return func(ab AsBytes) uint64 {
		m := murmur3.New64()
		m.Write(ab.Bytes())
		h64 := m.Sum64()
		hi := h64 >> 32
		lo := uint32(h64)
		return hi + uint64(lo)*value
	}
}

func MurMurHash(value AsBytes) uint64 {
	m := murmur3.New64()
	m.Write(value.Bytes())
	h64 := m.Sum64()
	return h64
}
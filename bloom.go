package maybe

import (
	"errors"
	"github.com/spaolacci/murmur3"
	"github.com/willf/bitset"
)
// BloomFilter is a probabilistic data structure to test weather a given value is a member of a  given  set
// (represented by this bloom filter instance).
//
// The main advantage of this data structure is the low memory footprint, it can be used to do set membership on
// a very large multiset with very low memory footprint (in the order of `Kb`) as opposed to storing this in a set
// which can have a memory footprint in the order of `Mg` or `Gb`.
// Of course this comes with the additional tradeoff of certainty, if a bloom filter `Has` method returns `true`
// this means that the element is already a member of this set with a certainty of `x%` as opposed to the `100%` given
// by a set-like data structure.
// however if a `Has` call on a bloom filter returns `false` then the element is not in the set `100%`.
//
// The certainty factor is determined by the size of the bitset in the internal implementation, and the number of hash
// functions you apply to the value, the corollary is this: the bigger the bitset and the more hash functions you have
// the more accurate the result of the bloom filter (or the less false positives you will get).
//
// You can read the original paper for the math behind it.
type BloomFilter struct {
 	set *bitset.BitSet
 	hashes []HashFunc
 	bits uint32
 	times uint32
 	ownHashFn bool
}

// NewBloomFilter creates a new BloomFilter with a given size and given a number of hash functions.
// The numFunc will indicate the number of times the default hash function is going to be used to create
// additional hash functions.
// The default hash base hash function is the murmur3, any other hash function would've done the job however.
func NewBloomFilter(size uint, numFunc uint) (*BloomFilter, error) {
	return bloomWithHashes(size, []HashFunc {
		func(i AsBytes) uint64 {
			murmur := murmur3.New64()
			murmur.Write(i.Bytes())
			return murmur.Sum64()
		},
	}, numFunc, true)
}

// NewBloomFilterWithHashes creates a bloom filter with user defined hash functions
func NewBloomFilterWithHashes(size uint, hashes []HashFunc) (*BloomFilter, error) {
	return bloomWithHashes(size, hashes, 1, false)
}

func bloomWithHashes(size uint, hashes[]HashFunc, numFunc uint, ownHash bool) (*BloomFilter, error) {
	if len(hashes) == 0 {
		return nil, errors.New("you should provide at least one hash function")
	}
	return &BloomFilter{
		set: bitset.New(size),
		hashes: hashes,
		bits: uint32(size),
		times: uint32(numFunc),
		ownHashFn: ownHash,
	}, nil
}

// Add adds the value to the bloom filter.
//
// The value should implement the AsBytes interface for the hash functions to work.
func (bf *BloomFilter) Add(value AsBytes) {
	if bf.ownHashFn {
		bf.hashMany(value, bf.hashes[0])
		return
	}
	for _, fn := range bf.hashes {
		index := uint(fn(value) % uint64(bf.bits))
		bf.set.Set(index)
	}
}

// Has returns weather the value element exists in the bloom filter.
//
// If the return value is true, than the value might have been added to the bloom filter
// otherwise the value has **definitely** never been added
func (bf *BloomFilter) Has(value AsBytes) bool {
	if bf.ownHashFn {
		return bf.findMany(value, bf.hashes[0])
	}
	for _, fn := range bf.hashes {
		index := uint(fn(value) % uint64(bf.bits))
		if !bf.set.Test(index) {
			return false
		}
	}
	return true
}

// hasMany will the formula: hash(i) = low(h(0)) + (hi(hash(0)) * i) to generate all of the of the
// `times` hash functions bit positions.
func(bf *BloomFilter) hashMany(value AsBytes, hashFn HashFunc) {
	h64 := hashFn(value)
	hlo := uint32(h64)
	hhi := uint32(h64 >> 32)
	for i := uint32(1); i <= bf.times; i = i + 1 {
		comb := hlo + (i * hhi)
		bf.set.Set(uint(comb % bf.bits))
	}
}

// findMany applies the same formula as hasMany to find the weather the value's hashes are all set.
func(bf *BloomFilter) findMany(value AsBytes, hashFn HashFunc) bool {
	h64 := hashFn(value)
	hlo := uint32(h64)
	hhi := uint32(h64 >> 32)
	for i := uint32(1); i <= bf.times; i = i + 1 {
		comb := hlo + (i * hhi)
		if !bf.set.Test(uint(comb % bf.bits)) {
			return false
		}
	}
	return true
}

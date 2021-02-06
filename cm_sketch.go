package maybe

import "math"

// CountMinSketch is a probabilistic data structure that is used to count the number of occurrences of a given value in
// a stream of events
//
// The data structure's bias is determined by the width and depth of the table, the bigger they are the better precision
// you get in point count queries.
// The we use a linear hash function based on the murmur3 algorithm in this implementation, the algorithm used is not highly
// important in the functioning of the data structure, that being said of course, using sha1 or some cryptographic hashing
// algorithm (usually slow) will impact the performance of adding and querying.
// A count min sketch is biased estimator data structure, it can over estimate (report more elements then the ones added)
// but it will never under estimate.
type CountMinSketch struct {
	width uint32
	depth uint32
	counters [][]uint64
	hashes []HashFunc
}

// NewCountMinSketch will create a count min sketch based of the width and depth given.
//
// The values will determine the size of the underlying table, and the number of hash functions generated, bigger values
// will give you better estimates, but it will also impact your memory / cpu footprint.
func NewCountMinSketch (width, depth uint32) *CountMinSketch {
	arr := make([][]uint64, depth)
	for i := range arr {
		arr[i] = make([]uint64, width)
	}
	hashes := make([]HashFunc, depth)

	for i := uint32(1); i <= depth; i = i + 1 {
		hashes[i - 1] = MurmurFrom(uint64(i))
	}

	return &CountMinSketch{
		width: width,
		depth: depth,
		counters: arr,
		hashes: hashes,
	}
}

// Increment is equivalent to adding one element of type represented by value.
func(cm *CountMinSketch) Increment(value AsBytes)  {
	cm.Add(value, 1)
}

// Add will increment the buckets corresponding to the value by count.
func(cm *CountMinSketch) Add(value AsBytes, count uint64)  {
	for i := uint32(0); i < cm.depth; i = i + 1 {
		id := cm.hashes[i](value) % uint64(cm.width)
		cm.counters[i][id] += count
	}
}

// Count will return an estimate of the number of values of type value that has been added to the sketch.
//
// The count estimate is always greater than or equal the actual count.
func(cm *CountMinSketch) Count(value AsBytes) uint64 {
	res := uint64(math.MaxUint64)
	for i := uint32(0); i < cm.depth; i = i + 1 {
		id := cm.hashes[i](value) % uint64(cm.width)
		res = min(cm.counters[i][id], res)
	}
	return res
}

func min(a uint64, res uint64) uint64 {
	if a > res {
		return res
	}
	return a
}



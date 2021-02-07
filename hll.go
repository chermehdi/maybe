package maybe

import (
	"errors"
	"math"
	"math/bits"
)

// HyperLogLog is data structure for estimating cardinality with a accuracy of ` 1 - 1.04 / sqrt(m)` with
// m defined as `2^b`, the reference is from this paper: http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf
// This being said the accuracy of the data structure is proportional to the value of b.
type HyperLogLog struct {
	counters []uint
	b        uint64
	fn       HashFunc
	alpha    float64
}

// NewHyperLogLog creates a new HLL based on the the number of counters
func NewHyperLogLog(b uint) (*HyperLogLog, error) {
	if b > 32 {
		return nil, errors.New("cannot create more the `1 << 32` counters array")
	}
	counters := make([]uint, 1<<b)
	alpha := computeAlpha(b, 1<<b)
	return &HyperLogLog{
		counters: counters,
		b:        uint64(b),
		fn:       MurMurHash,
		alpha:    alpha,
	}, nil
}

// Add adds the given observable (value) to the HLL instance
func (hll *HyperLogLog) Add(value AsBytes) {
	x := hll.fn(value)
	j := uint(x >> (64 - hll.b))
	// everything is offset by b so the r should be leading zeros of the first 64 - b bits portion of the hash
	// minus b.
	leading := uint(bits.LeadingZeros64(x & (uint64(1) <<(64 - hll.b) - 1)))
	r := leading - uint(hll.b) + 1
	hll.counters[j] = max(hll.counters[j], r)
}

func (hll *HyperLogLog) Cardinality() uint64 {
	var sum float64 = 0
	var zeros float64 = 0
	for i := 0; i < len(hll.counters); i = i + 1 {
		v := hll.counters[i]
		sum += 1.0 / float64(uint64(1)<<v)
		if v == 0 {
			zeros += 1
		}
	}
	estimate := hll.alpha * (1.0 / sum)
	m := float64(len(hll.counters))
	if estimate < 2.5*m {
		// Small range correction
		if zeros != 0 {
			return uint64(math.Round(m * math.Log(m/zeros)))
		}
	}

	return uint64(math.Round(estimate))
}

// Values from the paper mentioned in the type documentation
func computeAlpha(p, m uint) float64 {
	mf := float64(m)
	switch p {
	case 4:
		return 0.673 * mf * mf
	case 5:
		return 0.697 * mf * mf
	case 6:
		return 0.709 * mf * mf
	default:
		return (0.7213 / (1 + 1.079/mf)) * mf * mf
	}
}

func max(a, b uint) uint {
	if a < b {
		return b
	}
	return a
}

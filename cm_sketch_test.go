package maybe

import (
	"testing"
)

func TestCountMinSketch_Invariant(t *testing.T) {
	cm := NewCountMinSketch(2014, 5)
	counter := make(map[string]uint64)
	for i := 0; i < 1000000; i = i + 1 {
		word := randWord(3)
		counter[word] += 1
		cm.Increment(String(word))
		// CountMinSketch data structure never underestimates
		if counter[word] > cm.Count(String(word)) {
			t.Fatalf("test: %d The count of the count min sketch is bigger than the size of the actual counter, which violates the invariant, expected: '%d' found '%d'", i, counter[word], cm.Count(String(word)))
		}
	}
}

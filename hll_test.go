package maybe

import (
	"math"
	"testing"
)

type Int int

func (i Int) Bytes() []byte {
	return []byte{
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

func TestHyperLogLog_Add(t *testing.T) {
	hll, err := NewHyperLogLog(8)
	if err != nil {
		t.Fatal(err)
	}
	hll.Add(String("a"))
	hll.Add(String("b"))
	hll.Add(String("c"))
	hll.Add(String("d"))
	hll.Add(String("e"))
	if hll.Cardinality() != 5 {
		t.Fatalf("Expected 5 found %d", hll.Cardinality())
	}
}

func TestNewHyperLogLog_Cardinality(t *testing.T) {
	hll, _ := NewHyperLogLog(12)
	size := 10000000
	for i := 0; i < size; i += 1 {
		hll.Add(Int(i))
	}
	estimate := hll.Cardinality()
	err := math.Abs(float64(estimate)-float64(size)) / float64(size)
	if err > .1 {
		t.Fatalf("Expected an error smaller than .1 got %.5f", err)
	}
}

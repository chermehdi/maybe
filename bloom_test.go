package maybe

import (
	"bytes"
	"math/rand"
	"testing"
)

type String string

type Set struct {
	st map[string]bool
}

func (s *Set) Add(value string) {
	s.st[value] = true
}

func (s *Set) Contains(value string) bool {
	_, has := s.st[value]
	return has
}

func (s String) Bytes() []byte {
	return []byte(s)
}

func TestBloomFilter_Add(t *testing.T) {
	bf, err := NewBloomFilter(100, 7)
	if err != nil {
		t.Fatal(err)
	}
	set := &Set{
		st: make(map[string]bool),
	}

	for i := 0; i <= 100000; i += 1 {
		word := randWord(5)
		if !bf.Has(String(word)) && set.Contains(word) {
			t.Fatalf("Could not find word '%s' while it should have", word)
		}
		set.Add(word)
		bf.Add(String(word))
	}
}

func randWord(size uint) string {
	var buff bytes.Buffer
	for i := uint(0); i < size; i = i + 1 {
		buff.WriteRune(rune('a' + rand.Int()%26))
	}
	return buff.String()
}

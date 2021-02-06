# Maybe

A library containing a set of probabilstic data structures that I found interesting while researching the subject.

## Contents

### Bloom Filter

A bloom filter is a great data structure to do set membership queries on large number of elements with very little
memory footprint.

An example usage:

```go
package main

import (
	"fmt"
	"github.com/chermehdi/maybe"
)

type String string

func (i String) Bytes() []byte {
	return []byte(i)
}

func main() {
	bf, _ := maybe.NewBloomFilter(1024, 3)
	bf.Add(String("hello"))

	if bf.Has(String("hello")) {
		fmt.Println("Found Hello in bloom filter")
	}
}
```


### Count min sketch

A count min sketch is a data structure to estimate the count of elements in a large set of events with a very little memory 
footprint.

An example usage:

```go
package main

import (
	"fmt"
	"github.com/chermehdi/maybe"
)

type String string

func (i String) Bytes() []byte {
	return []byte(i)
}

func main() {
	sk := maybe.NewCountMinSketch(1024, 3)
	sk.Increment(String("hello"))
	sk.Add(String("hello1"), 2)

    fmt.Println(sk.Count(String("hello")))
}
```

## Resources 
- [Bloom filters](https://en.wikipedia.org/wiki/Bloom_filter)
- [Count min sketch](https://en.wikipedia.org/wiki/Count%E2%80%93min_sketch)
- [Count min sketch paper](https://sites.google.com/site/countminsketch/cm-latin.pdf)

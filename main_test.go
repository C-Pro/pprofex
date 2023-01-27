package main

import (
	"strings"
	"testing"
)

func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// genesis block with randomized hash
		genesis := &Block{
			data: []byte("genesis" + strings.Repeat("#", blockDataSize-len("genesis"))),
			hash: genData(hashSize),
		}

		genBlockchain(genesis, 1000)
	}
}

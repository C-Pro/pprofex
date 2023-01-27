package main

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"hash"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
)

var (
	hashFn        hash.Hash
	hashSize      int
	blockDataSize = 128
	blockSize     int
)

func init() {
	hashFn = sha512.New()
	hashSize = hashFn.Size()
	blockSize = blockDataSize + hashSize
}

type Block struct {
	data []byte
	hash []byte
}

func (b *Block) Bytes() []byte {
	return append(b.data, b.hash...)
}

func (b *Block) NewBlock(data []byte) *Block {
	hashFn.Reset()
	hashFn.Write(data)
	hashFn.Write(b.hash)
	return &Block{
		data: data,
		hash: hashFn.Sum(nil),
	}
}

func genData(n int) []byte {
	var d []byte
	for i := 0; i < n; i++ {
		d = append(d, byte(rand.Intn(255)))
	}
	return d
}

func genBlockchain(genesis *Block, n int) []byte {
	var blockChain []byte
	blockChain = append(blockChain, genesis.Bytes()...)

	block := genesis

	// Append new blocks with random data
	for i := 0; i < n; i++ {
		block = block.NewBlock(genData(blockDataSize))
		blockChain = append(blockChain, block.Bytes()...)
	}

	return blockChain
}

func validateBlockchain(blockChain []byte) error {
	var prevBlock *Block
	for i := 0; i < len(blockChain)/blockSize; i++ {
		startOffset := i * blockSize
		block := &Block{
			data: blockChain[startOffset : startOffset+blockDataSize],
			hash: blockChain[startOffset+blockDataSize : startOffset+blockSize],
		}

		// Check block validity
		if prevBlock != nil {
			hashFn.Reset()
			hashFn.Write(block.data)
			hashFn.Write(prevBlock.hash)
			expected := hashFn.Sum(nil)
			if !bytes.Equal(expected, block.hash) {
				return fmt.Errorf("block %d hash is invalid:\nGOT: %x\nEXP: %x", i, block.hash, expected)
			}
		}

		prevBlock = block
	}
	return nil
}

func main() {
	f, err := os.Create("pprof.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}

	genesis := &Block{
		data: []byte("genesis" + strings.Repeat("#", blockDataSize-len("genesis"))),
		hash: genData(hashSize),
	}

	blockChain := genBlockchain(genesis, 1_000_000)
	f2, err := os.Create("mem.out")
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	pprof.WriteHeapProfile(f2)

	pprof.StopCPUProfile()

	if err := validateBlockchain(blockChain); err != nil {
		panic(err)
	}
}

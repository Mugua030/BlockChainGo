package main

import (
	"fmt"
)

type Block struct {
	PreBlockHash []byte
	Hash         []byte
	Data         []byte
}

func NewBlock(data string, preBlock []byte) *Block {

	block := Block{
		PreBlockHash: preBlock,
		Hash:         nil,
		Data:         []byte(data),
	}
	return &block
}

func main() {
	block := NewBlock("Love BlockChain", []byte{})

	fmt.Printf("PreBlock: %s\n", block.PreBlockHash)
	fmt.Printf("Hash: %s\n", block.Hash)
	fmt.Printf("Data: %s\n", block.Data)

}

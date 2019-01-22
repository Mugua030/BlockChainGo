package main

import (
	"crypto/sha256"
	"encoding/hex"
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
	block.SetHash()

	return &block
}
func (b *Block) SetHash() {

	var blockInfo []byte
	blockInfo = append(blockInfo, b.PreBlockHash...)
	blockInfo = append(blockInfo, b.Hash...)
	blockInfo = append(blockInfo, b.Data...)

	hash := sha256.Sum256(blockInfo)
	b.Hash = hash[:]
}

func main() {
	block := NewBlock("Love BlockChain", []byte{})

	fmt.Printf("PreBlock: %s\n", block.PreBlockHash)
	//fmt.Printf("Hash: %s\n", block.Hash)
	encodeStr := hex.EncodeToString(block.Hash)
	fmt.Printf("Hash: %s\n", encodeStr)
	fmt.Printf("Data: %s\n", block.Data)

}

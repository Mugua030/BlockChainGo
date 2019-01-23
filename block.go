package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

const BITS = 16

type Block struct {
	Version      uint64
	PreBlockHash []byte
	Hash         []byte
	MerkelRoot   []byte
	TimeStamp    uint64
	Bits         uint64 //难度值 difficultly
	Nonce        uint64 //随机数
	Data         []byte
}

func NewBlock(data string, preBlock []byte) *Block {

	block := Block{
		Version:      0,
		PreBlockHash: preBlock, //record pre block hashValue :: block info => to hash
		Hash:         nil,      // current blockinfo hash
		MerkelRoot:   nil,
		TimeStamp:    uint64(time.Now().Unix()),
		Bits:         BITS,
		Nonce:        0,
		Data:         []byte(data), // infos: msg
	}
	//block.SetHash()
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash

	return &block
}

//SetHash set current block hash
func (b *Block) SetHash() {

	var blockInfo []byte
	blockInfo = append(blockInfo, uint2bytes(b.Version)...)
	blockInfo = append(blockInfo, b.PreBlockHash...)
	blockInfo = append(blockInfo, b.Hash...)
	blockInfo = append(blockInfo, b.MerkelRoot...)
	blockInfo = append(blockInfo, uint2bytes(b.TimeStamp)...)
	blockInfo = append(blockInfo, uint2bytes(b.Bits)...)
	blockInfo = append(blockInfo, uint2bytes(b.Nonce)...)
	blockInfo = append(blockInfo, b.Data...)

	hash := sha256.Sum256(blockInfo)
	fmt.Printf("%x", hash)
	b.Hash = hash[:]
}

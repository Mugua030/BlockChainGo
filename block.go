package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

const BITS = 10

type Block struct {
	Version      uint64
	PreBlockHash []byte
	Hash         []byte
	MerkelRoot   []byte
	TimeStamp    uint64
	Bits         uint64 //难度值 difficultly
	Nonce        uint64 //随机数
	//Data         []byte
	Transactions []*Transaction
}

func NewBlock(txs []*Transaction, preBlock []byte) *Block {

	block := Block{
		Version:      0,
		PreBlockHash: preBlock, //record pre block hashValue :: block info => to hash
		Hash:         nil,      // current blockinfo hash
		MerkelRoot:   nil,
		TimeStamp:    uint64(time.Now().Unix()),
		Bits:         BITS,
		Nonce:        0,
		//Data:         []byte(data), // infos: msg
		Transactions: txs,
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
	//blockInfo = append(blockInfo, b.Data...)

	hash := sha256.Sum256(blockInfo)
	//fmt.Printf("%x", hash)
	b.Hash = hash[:]
}
func (b *Block) toBytes() []byte {
	return []byte("222222")
}
func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil {
		//return nil
		panic(err)
	}
	return buffer.Bytes()
}

func Deserialize(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

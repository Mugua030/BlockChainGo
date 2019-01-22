package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

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

var genesisBlock string = "genesisBlock begin"

func NewBlock(data string, preBlock []byte) *Block {

	block := Block{
		Version:      0,
		PreBlockHash: preBlock, //record pre block hashValue :: block info => to hash
		Hash:         nil,      // current blockinfo hash
		MerkelRoot:   nil,
		TimeStamp:    uint64(time.Now().Unix()),
		Bits:         0,
		Nonce:        0,
		Data:         []byte(data), // infos: msg
	}
	block.SetHash()

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
	b.Hash = hash[:]
}
func uint2bytes(num uint64) []byte {

	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

//BlockChain  用户切片来当存储容器
type BlockChain struct {
	Blocks []*Block
}

//AddBlock 添加区块 hash= hash blockinfo
func (bc *BlockChain) AddBlock(data string) {

	len := len(bc.Blocks)
	lastBlock := bc.Blocks[len-1]
	newBlock := NewBlock(data, lastBlock.Hash)

	bc.Blocks = append(bc.Blocks, newBlock)

}

func NewBlockChain() *BlockChain {

	genesisBlock := NewBlock(genesisBlock, []byte{})

	bc := BlockChain{
		Blocks: []*Block{genesisBlock},
	}

	return &bc
}

func main() {

	blockChain := NewBlockChain()
	blockChain.AddBlock("love blockchain")
	blockChain.AddBlock("bb block")

	for _, b := range blockChain.Blocks {
		//fmt.Printf("preBlock : %s\n", hex.EncodeToString(b.PreBlockHash))
		//fmt.Printf("Hash : %s\n", hex.EncodeToString(b.Hash))
		fmt.Printf("preBlockHash : %x\n", b.PreBlockHash)
		fmt.Printf("Hash : %x\n", string(b.Hash))
		fmt.Printf("Data: %s\n", string(b.Data))
		fmt.Println()
	}
	//Todo:: 随机数  难度值 没有用上 2.hash没有规律

}

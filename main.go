package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var genesisBlock string = "genesisBlock begin"

func uint2bytes(num uint64) []byte {

	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func main() {
	//type Block struct {
	//Version      uint64
	//PreBlockHash []byte
	//Hash         []byte
	//MerkelRoot   []byte
	//TimeStamp    uint64
	//Bits         uint64 //难度值 difficultly
	//Nonce        uint64 //随机数
	//Data         []byte
	//}

	blockChain := NewBlockChain()
	blockChain.AddBlock("love blockchain")
	blockChain.AddBlock("bb block")

	for _, b := range blockChain.Blocks {
		//fmt.Printf("preBlock : %s\n", hex.EncodeToString(b.PreBlockHash))
		//fmt.Printf("Hash : %s\n", hex.EncodeToString(b.Hash))
		fmt.Printf("version: %d\n", b.Version)
		fmt.Printf("PreBlockHash: %x\n", b.PreBlockHash)
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Printf("MerkelRoot: %x\n", b.MerkelRoot)
		fmt.Printf("TimeStamp: %d\n", b.TimeStamp)
		fmt.Printf("Bits: %d\n", b.Bits)
		fmt.Printf("Nonce: %d\n", b.Nonce)
		fmt.Printf("Data: %s\n", string(b.Data))
		fmt.Println()
	}
	//Todo:: 随机数  难度值 没有用上 2.hash没有规律

}

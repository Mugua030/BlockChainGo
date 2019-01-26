package main

import (
	"bytes"
	"encoding/binary"
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

	//	blockChain, err := NewBlockChain()
	//if err != nil {
	//panic(err)
	//}
	//defer blockChain.db.Close()

	//blockChain.AddBlock("love blockchain")
	//blockChain.AddBlock("bb block")

	//blockChain.PrintBC()
	cli := CLI{}
	cli.Run()

}

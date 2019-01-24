package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//ProofOfWork 工作量证明 to generate block
type ProofOfWork struct {
	block  Block
	target big.Int
}

func NewProofOfWork(block Block) *ProofOfWork {

	//	targetStr :=
	//"0001000000000000000000000000000000000000000000000000000000000000"

	//bigIntTmp := big.Int{}
	//bigIntTmp.SetString(targetStr, 16)
	bigIntTmp := big.NewInt(1)
	bigIntTmp.Lsh(bigIntTmp, 256-BITS) // 返回一个指针

	pow := ProofOfWork{
		block:  block,
		target: *bigIntTmp,
	}

	return &pow
}

//Run 挖矿...
func (p *ProofOfWork) Run() (uint64, []byte) {
	fmt.Println("挖矿中...")

	var nonce uint64
	var hash [32]byte
	for {
		fmt.Printf("Run::  %x\r", hash)

		//对拼接好的数据进行256hash
		hash = sha256.Sum256(p.prepareData(nonce))

		bigIntTmp := big.Int{}
		bigIntTmp.SetBytes(hash[:])

		if bigIntTmp.Cmp(&p.target) == -1 {
			fmt.Printf("挖矿成功, hash: %x, nonce: %d\n\r", hash, nonce)
			fmt.Println()
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}
func (p *ProofOfWork) prepareData(nonce uint64) []byte {

	b := p.block

	//var blockInfo []byte
	//blockInfo = append(blockInfo, uint2bytes(b.Version)...)
	//blockInfo = append(blockInfo, b.PreBlockHash...)
	//blockInfo = append(blockInfo, b.Hash...)
	//blockInfo = append(blockInfo, b.MerkelRoot...)
	//blockInfo = append(blockInfo, uint2bytes(b.TimeStamp)...)
	//blockInfo = append(blockInfo, uint2bytes(b.Bits)...)

	tmp := [][]byte{
		uint2bytes(b.Version),
		b.PreBlockHash,
		b.Hash,
		b.MerkelRoot,
		uint2bytes(b.TimeStamp),
		uint2bytes(b.Bits),
		uint2bytes(nonce),
		b.Data,
	}

	//blockInfo = append(blockInfo, uint2bytes(b.Nonce)...)
	//blockInfo = append(blockInfo, uint2bytes(nonce)...) // 用b.Nonce  产生的数据不对
	//blockInfo = append(blockInfo, b.Data...)
	blockInfo := bytes.Join(tmp, []byte{})

	return blockInfo
}

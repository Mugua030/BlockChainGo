package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

const CON_reward = 12.2

//TXInput 交易输入
type TXInput struct {
	TXID      []byte
	Index     int64
	ScriptSig string
}

//TXOutput 交易输出
type TXOutput struct {
	Value        float64
	ScriptPubKey string
}

type Transaction struct {
	TxId      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
	TimeStamp uint64
}

func (tx *Transaction) SetTXHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(buffer.Bytes())

	tx.TxId = hash[:]
}

func NewCoinbaseTX(data string, miner /*矿工地址*/ string) *Transaction {

	txinput := TXInput{
		TXID:      nil,
		Index:     -1,
		ScriptSig: data,
	}
	txoutput := TXOutput{
		Value:        CON_reward,
		ScriptPubKey: miner,
	}

	tx := Transaction{
		nil,
		[]TXInput{txinput},
		[]TXOutput{txoutput},
		uint64(time.Now().Unix()),
	}

	tx.SetTXHash()

	return &tx
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	var returnedUtxos map[string][]int64
	var calcMoney float64

	returnedUtxos, calcMoney = bc.FindNeedUtxos(from, amount)

	if calcMoney < amount {
		//余额不足
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//余额充足，拼接inputs
	for txid, indexArray := range returnedUtxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(txid), i, from}
			inputs = append(inputs, input)
		}
	}

	//拼接output
	output := TXOutput{amount, to}
	outputs = append(outputs, output)
	if calcMoney > amount {
		balanceOutput := TXOutput{calcMoney - amount, from}
		outputs = append(outputs, balanceOutput)
	}

	tx := Transaction{
		TxId:      nil,
		TXInputs:  inputs,
		TXOutputs: outputs,
		TimeStamp: uint64(time.Now().Unix()),
	}
	tx.SetTXHash()

	return &tx
}

func (tx *Transaction) IsCoinbaseTx() bool {

	if len(tx.TXInputs) == 1 && tx.TXInputs[0].TXID == nil && tx.TXInputs[0].Index == -1 {
		fmt.Printf("这个是挖矿交易，不统计!")
		return true
	}

	return false
}

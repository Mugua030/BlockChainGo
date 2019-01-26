package main

import (
	"fmt"
)

func (c *CLI) createBC() {
	fmt.Println("Invoke CreateBC ... ")
	bc := CreateBlockChain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	fmt.Println("yes run hear..")
}
func (c *CLI) send(from, to string, amount float64, miner string, data string) {
	fmt.Println("Invoke send ... ")

	bc := GetBlockChain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	txs := []*Transaction{}
	//挖矿交易
	coinbaseTx := NewCoinbaseTX(data, miner)
	txs = append(txs, coinbaseTx)
	//普通交易

	ok := bc.AddBlock(txs)
	if !ok {
		fmt.Println("add block fail !")
		return
	}
	fmt.Println("add block success !!")
}
func (c *CLI) printBC() {
	fmt.Println("Invoke printchain ... ")
	bc := GetBlockChain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	bc.PrintBC()
}

func (c *CLI) getBalance(address string) {
	fmt.Println("Invoke getBalance...")
	bc := GetBlockChain()
	if bc == nil {
		return
	}
	defer bc.db.Close()

	bc.GetBalance(address)
}

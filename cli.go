package main

import (
	"fmt"
	"os"
	"strconv"
)

const Usage = `
	./blk createBC  "Create BlockChain DB"
	./blk printchain, print blockchain
	./blk send FROM TO AMOUNT MINER DATA
	./blk getBalance address "获取指定地址的余额"
`

type CLI struct {
	bc *BlockChain
}

func (cli *CLI) Run() {

	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Println(Usage)
		os.Exit(1)
	}

	switch cmds[1] {
	case "createBC":
		cli.createBC()
	case "send":
		if len(cmds) != 7 {
			fmt.Println("params error")
			os.Exit(1)
		}
		from := cmds[1]
		to := cmds[2]
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		data := cmds[6]
		cli.send(from, to, amount, miner, data)
	case "printchain":
		cli.printBC()
	case "getBalance":
		//balance := cli.getBalance()
		//fmt.Println("balance: ", balance)
		address := cmds[2]
		cli.getBalance(address)
	default:
		fmt.Println(Usage)
		os.Exit(1)
	}

}

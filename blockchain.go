package main

import (
	"errors"
	"fmt"
	"time"

	"BlockChainGo/lib/bolt"
)

//BlockChain  用户切片来当存储容器
type BlockChain struct {
	//Blocks []*Block
	db   *bolt.DB
	tail []byte
}
type Iterator struct {
	db          *bolt.DB
	currentHash []byte
}

type UTXOInfo struct {
	TxId   []byte
	Index  int64
	Output TXOutput
}

const CON_DBFILE = "blockChain.db"
const CON_BUCKET = "blockBucket"
const CON_LASTHASHKEY = "lastHashKey"

func (bc *BlockChain) NewIterator() *Iterator {
	it := Iterator{
		db:          bc.db,
		currentHash: bc.tail,
	}
	return &it
}

//AddBlock 添加区块 hash= hash blockinfo
func (bc *BlockChain) AddBlock(txs []*Transaction) bool {

	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CON_BUCKET))
		if b == nil {
			return errors.New("no bucket")
		}
		preHash := bc.tail

		newBlock := NewBlock(txs, preHash)
		newBlockBy := newBlock.Serialize()
		err := b.Put(newBlock.Hash, newBlockBy)
		if err != nil {
			return err
		}
		err = b.Put([]byte(CON_LASTHASHKEY), newBlock.Hash)
		if err != nil {
			return err
		}
		bc.tail = newBlock.Hash

		//打印一下
		blockInfo := b.Get(bc.tail)
		block, _ := Deserialize(blockInfo)
		fmt.Printf("%x", block)

		return nil
	})

	if err != nil {
		fmt.Printf("addNewBlock err: %x ", err)
		return false
	}

	return true
}

func (it *Iterator) Next() *Block {

	var block Block
	err := it.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(CON_BUCKET))
		if bkt == nil {
			return errors.New("IteratorNext:: Get Bucket fail")
		}
		fmt.Printf("NextCurrentHash: %x\n", it.currentHash)
		blockBytesInfo := bkt.Get(it.currentHash)
		blockp, _ := Deserialize(blockBytesInfo)
		block = *blockp

		return nil
	})
	if err != nil {
		panic(err)
	}
	//fmt.Printf("next:: preblockhSH, %x\n", block.PreBlockHash)

	//block = Block{}
	it.currentHash = block.PreBlockHash

	return &block
}

func (bc *BlockChain) PrintBC() {

	it := bc.NewIterator()

	for {
		block := it.Next()
		fmt.Printf("Version: %d\n", block.Version)
		fmt.Printf("preBlockHash: %x\n", block.PreBlockHash)
		fmt.Printf("curHash: %x\n", block.Hash)
		fmt.Printf("merkelRoot: %x\n", block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("timeStamp: %s\n", timeFormat)
		fmt.Printf("Bits: %d\n", block.Bits)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Data: %s\n", block.Transactions[0].TXInputs[0].ScriptSig)

		fmt.Println()
		if len(block.PreBlockHash) == 0 {
			fmt.Printf("区块链遍历完成 !\n")
			break
		}
	}
}

//CreateBlockChain 创建区块链
func CreateBlockChain() *BlockChain {

	//判断是否已存在db file
	if isFileExist(CON_DBFILE) {
		fmt.Println("BlockChain isExist")
		return nil
	}

	var bc BlockChain
	//db save data
	db, err := bolt.Open(CON_DBFILE, 0600, nil)
	if err != nil {
		//return nil, errors.New("NewBlockChain:: fail,open bolt db fail")
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(CON_BUCKET))
		if bkt == nil {
			//return errors.New("NewBlockChain:: no bucket ")
			b, err := tx.CreateBucket([]byte(CON_BUCKET))
			if err != nil {
				return errors.New("NewBlockChain:: create bucket fail")
			}

			//genesisBlock := NewBlock(genesisBlock, []byte{})
			coinbaseTx := NewCoinbaseTX(genesisBlock, "中本本")
			genesisBlock := NewBlock([]*Transaction{coinbaseTx}, []byte{})

			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize()) // ?? todo
			if err != nil {
				return err
			}
			err2 := b.Put([]byte(CON_LASTHASHKEY), genesisBlock.Hash)
			if err2 != nil {
				return err2
			}
			bc.tail = genesisBlock.Hash
		} else {
			lastHash := bkt.Get([]byte(CON_LASTHASHKEY))
			bc.tail = lastHash
		}
		return nil //commit transaction
	})

	bc.db = db

	return &bc
}
func GetBlockChain() *BlockChain {

	if !isFileExist(CON_DBFILE) {
		fmt.Println("Please to create blockChain")
		return nil
	}
	var bc BlockChain
	db, err := bolt.Open(CON_DBFILE, 0600, nil)
	if err != nil {
		panic("Please to create blockchain first!!")
	}
	db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(CON_BUCKET))
		if bkt == nil {
			panic("bucket can not empty")
		}
		lastHash := bkt.Get([]byte(CON_LASTHASHKEY))
		bc.tail = lastHash

		return nil
	})
	bc.db = db

	return &bc
}

func (bc *BlockChain) GetBalance(address string) float64 {
	var totalMoney float64
	utxos := bc.FindMyUtxos(address)
	for _, utxo := range utxos {
		totalMoney += utxo.Output.Value
	}

	return totalMoney
}
func (bc *BlockChain) FindNeedUtxos(from string, amount float64) (map[string][]int64, float64) {

	var returnedUtxos = make(map[string][]int64)
	var calcMoney float64

	utxoinfos := bc.FindMyUtxos(from)

	for _, utxoinfo := range utxoinfos {
		calcMoney += utxoinfo.Output.Value
		key := string(utxoinfo.TxId)
		returnedUtxos[key] = append(returnedUtxos[key], int64(utxoinfo.Index))
		if calcMoney > amount {
			break
		}
	}

	return returnedUtxos, calcMoney
}

//FindMyUtxos  查找我的未消费支出
func (bc *BlockChain) FindMyUtxos(address string) []UTXOInfo {
	var utxoInfos []UTXOInfo
	spentOutputs := make(map[string][]int64)

	it := bc.NewIterator()

	for {
		block := it.Next()

		//遍历交易
		for _, tx := range block.Transactions {
			//遍历input
			for _, input := range tx.TXInputs {
				if !tx.IsCoinbaseTx() {
					if input.ScriptSig == address {
						key := string(input.TXID)

						spentOutputs[key] = append(spentOutputs[key], input.Index)
						fmt.Println("找到历史消耗过的utxo")
					}
				}
			}
		}

		//遍历output
		for _, tx := range block.Transactions {
			key := string(tx.TxId)
		LAB_OUTPUT:
			for currindex, output := range tx.TXOutputs {
				if output.ScriptPubKey == address {
					if len(spentOutputs) != 0 {
						for _, spentIndex := range spentOutputs[key] {
							if int64(currindex) == spentIndex {
								fmt.Println("这个个output被消耗了!")
								continue LAB_OUTPUT
							}

						}
					}
					utxosinfo := UTXOInfo{tx.TxId, int64(currindex), output}
					utxoInfos = append(utxoInfos, utxosinfo)
				}
			}
		}

	} //for

	return utxoInfos
}

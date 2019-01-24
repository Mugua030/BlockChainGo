package main

import (
	"errors"
	"fmt"
	"time"

	"./lib/bolt"
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
func NewBlockChain() (*BlockChain, error) {

	//genesisBlock := NewBlock(genesisBlock, []byte{})

	//bc := BlockChain{
	//	Blocks: []*Block{genesisBlock},
	//}

	var bc BlockChain
	//db save data
	db, err := bolt.Open(CON_DBFILE, 0600, nil)
	if err != nil {
		return nil, errors.New("NewBlockChain:: fail,open bolt db fail")
	}
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(CON_BUCKET))
		if bkt == nil {
			//return errors.New("NewBlockChain:: no bucket ")
			b, err := tx.CreateBucket([]byte(CON_BUCKET))
			if err != nil {
				return errors.New("NewBlockChain:: create bucket fail")
			}
			genesisBlock := NewBlock(genesisBlock, []byte{})
			err = b.Put(genesisBlock.Hash, genesisBlock.toBytes()) // ?? todo
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

	return &bc, nil
}

//AddBlock 添加区块 hash= hash blockinfo
func (bc *BlockChain) AddBlock(data string) (bool, error) {

	//len := len(bc.Blocks)
	//lastBlock := bc.Blocks[len-1]
	//newBlock := NewBlock(data, lastBlock.Hash)

	//bc.Blocks = append(bc.Blocks, newBlock)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CON_BUCKET))
		if b == nil {
			return errors.New("no bucket")
		}
		preHash := bc.tail
		newBlock := NewBlock(data, preHash)
		newBlockBy, err := newBlock.Serialize()
		if err != nil {
			return err
		}
		err = b.Put(newBlock.Hash, newBlockBy)
		if err != nil {
			return err
		}
		err = b.Put([]byte(CON_LASTHASHKEY), newBlock.Hash)
		if err != nil {
			return err
		}
		bc.tail = newBlock.Hash

		//打印一下
		//blockInfo := b.Get(bc.tail)
		//block, _ := Deserialize(blockInfo)
		//fmt.Printf("%x", block)

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
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
		fmt.Printf("Data: %s\n", block.Data)

		fmt.Println()
		if len(block.PreBlockHash) == 0 {
			fmt.Printf("区块链遍历完成 !\n")
			break
		}
	}
}

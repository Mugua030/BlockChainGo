package main

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

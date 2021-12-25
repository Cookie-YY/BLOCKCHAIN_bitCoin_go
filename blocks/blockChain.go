package blocks

import (
	"blcokChain/utils"
	"log"
)

type BlockChain struct {
	Blocks   []*Block `json:"blocks"`
	HeadHash []byte   `json:"head_hash"` // the first hash value of the block in the  block-chain
	TailHash []byte   `json:"tail_hash"` // the hash of the last block in the block-chain
}

// GetOrCreateBlockChain : return an existed block-chain if there is one in the db else create one and save
func GetOrCreateBlockChain() *BlockChain {
	jdb := utils.GetJsonDB()
	bc := &BlockChain{}
	if jdb.IsExist() {
		err := jdb.ReadFromDB(bc)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		gBlock := newGenesisBlock()
		bc.Blocks = []*Block{gBlock}
		bc.TailHash = bc.Blocks[len(bc.Blocks)-1].Hash
		bc.HeadHash = gBlock.Hash
		err := utils.GetJsonDB().WriteToDB(bc)
		if err != nil {
			log.Fatal(err)
		}
	}

	return bc
}

// newGenesisBlock : return the first block
func newGenesisBlock() *Block {
	return NewBlock(
		"first block",
		[]byte{}, // there's no prev block
	)
}

// AddBlock : will save the blocks in db
func (bc *BlockChain) AddBlock(data string) {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash
	bc.Blocks = append(bc.Blocks, NewBlock(data, prevHash))
	bc.TailHash = bc.Blocks[len(bc.Blocks)-1].Hash
	err := utils.GetJsonDB().WriteToDB(bc)
	if err != nil {
		log.Fatal(err)
	}
}

// Iterator : gets a channel which contains the blocks in right order from TailHash
func (bc *BlockChain) Iterator() <-chan Block {
	c := make(chan Block)
	go func() {
		currentHash := bc.TailHash
		for {
			if len(currentHash) == 0 {
				break
			}
			for _, v := range bc.Blocks {
				if utils.GetBigIntWrapperFromBytes(v.Hash).EqualToAnotherBigIntWrapper(utils.GetBigIntWrapperFromBytes(currentHash)) {
					c <- *v
					currentHash = v.PrevHash
				}
			}
		}
		close(c)
	}()
	return c
}

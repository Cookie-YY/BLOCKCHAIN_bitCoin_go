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
//   The address is required when you want to create a block-chain
//   The address indicates the ADDRESS in GenesisBlock's only one transaction.
func GetOrCreateBlockChain(address string) *BlockChain {
	jdb := utils.GetJsonDB()
	bc := &BlockChain{}
	if jdb.IsExist() {
		err := jdb.ReadFromDB(bc)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		gBlock := newGenesisBlock(address, "cookie is so niubi")
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
func newGenesisBlock(address string, data string) *Block {
	coinBaseTX := NewCoinBaseTX(address, data)
	return NewBlock(
		[]*Transaction{coinBaseTX},
		[]byte{}, // there's no prev block
	)
}

// AddBlock : will save the blocks in db
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash
	bc.Blocks = append(bc.Blocks, NewBlock(txs, prevHash))
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

// GetBalanceOf : calculate all the rest money of this address
func (bc *BlockChain) GetBalanceOf(address string) float64 {
	_, amount := bc.FindUTXOsAndAmountsOf(address, -1) // -1 means get all money of the address
	return amount
}

// FindUTXOsAndAmountsOf : get all utxos of the address (if the amount > 0, it will return when the money >= amount)
func (bc *BlockChain) FindUTXOsAndAmountsOf(address string, amount float64) ([]UTXO, float64) {
	current := 0.0
	utxos := make([]UTXO, 0)
	txIDOutputIndexListMap := make(map[string][]int64) // txID -> txo's index
	for block := range bc.Iterator() {
		for _, tx := range block.Txs {
			// find all outputs: it means the rest money of the address
			for txOutputIndex, txOutput := range tx.TXOutputs {
				// if the current output has been spent
				if outputIndexList, ok := txIDOutputIndexListMap[string(tx.TXId)]; ok {
					if utils.IsExistItem(int64(txOutputIndex), outputIndexList) {
						continue
					}
				}
				if txOutput.PubKeyHash == address {
					utxos = append(utxos, UTXO{
						TXOutput: txOutput, TXId: tx.TXId, OutputIndex: int64(txOutputIndex),
					})
					current += txOutput.Amount
					if amount > 0 && current >= amount {
						return utxos, current
					}
				}
			}
			// check the input: it will spend the corresponding output
			if !tx.IsCoinBaseTX() { // there's no need to check coinBase transaction's input
				for _, txInput := range tx.TXInputs {
					if txInput.Sig == address {
						utils.AddIntElemIntoStringMap(txIDOutputIndexListMap, string(txInput.TXId), txInput.TXOutputIndex)
					}
				}
			}
		}
	}
	return utxos, current
}

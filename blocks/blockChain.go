package blocks

import (
	"blcokChain/utils"
	. "blcokChain/utils/consts"
	"bytes"
)

// BlockChain is the definition of a block-chain
// notice:
//   - Blocks: all the block in the block-chain. Notice, it could not be in the right order. When more than one node mining, it can't be in the right order
//  		we order the block in the right way through the hash of the last block
//   - HeadHash: the hash of the GenesisBlock
//   - TailHash: the hash of the last block: we use the hash value as the start to get the right order of the blocks
//          the TailHash may be more than one. When competition in mining occurs.
type BlockChain struct {
	Blocks   []*Block `json:"blocks"`
	HeadHash []byte   `json:"head_hash"` // the first hash value of the block in the  block-chain
	TailHash []byte   `json:"tail_hash"` // the hash of the last block in the block-chain
}

// GetOrCreateBlockChain : return an existed block-chain if there is one in the db else create one and save
//   The address is required when you want to create a block-chain
//   The address indicates the ADDRESS in GenesisBlock's only one transaction.
func GetOrCreateBlockChain(address string) (*BlockChain, error) {
	jdb := utils.GetJsonDB(JsonFileNameBlockChain, utils.NewJsonSerializer())
	bc := &BlockChain{}
	if jdb.IsExist() {
		err := jdb.ReadFromDB(bc)
		if err != nil {
			return nil, err
		}
	} else {
		if len(address) == 0 {
			return nil, ErrBlockChainNoGod()
		}
		gBlock := newGenesisBlock(address, MinerMsgOfGenesisBlock)
		bc.Blocks = []*Block{gBlock}
		bc.TailHash = bc.Blocks[len(bc.Blocks)-1].Hash
		bc.HeadHash = gBlock.Hash
		err := jdb.WriteToDB(bc)
		if err != nil {
			return nil, err
		}
	}

	return bc, nil
}

// newGenesisBlock : return the first block of the block-chain
func newGenesisBlock(address string, data string) *Block {
	coinBaseTX := NewCoinBaseTX(address, data)
	return NewBlock(
		[]*Transaction{coinBaseTX},
		[]byte{}, // there's no prev block
		nil,      // no need to pass block-chain
	)
}

// AddBlock : will save the blocks in db
//   about prevHash: when competition in mining occurs, the num of tailHash may be more than one.
func (bc *BlockChain) AddBlock(txs []*Transaction) error {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash
	block := NewBlock(txs, prevHash, bc)
	if block != nil && len(block.Txs) > 1 { // will not get in block-chain, if block is nil or only coinBase tx
		bc.Blocks = append(bc.Blocks, block)           // this block's prevHash should be chosen in the future.
		bc.TailHash = bc.Blocks[len(bc.Blocks)-1].Hash // tailHash may be more than one.
		return utils.GetJsonDB(JsonFileNameBlockChain, utils.NewJsonSerializer()).WriteToDB(bc)
	}
	return ErrBlockGetInBlockChainFailed()
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
				if bytes.Equal(v.Hash, currentHash) {
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
	_, amount := bc.FindUTXOsAndAmountsOf(utils.GetPubKeyHashFromAddress(address), -1) // -1 means get all money of the address
	return amount
}

// FindUTXOsAndAmountsOf : get all utxos of the sender's pubKeyHash
// Usage: if the amount > 0, it will return when the money >= amount
// FIXME: when counting the balance of the address. It will go wrong when every time finding in the block-chain (more than one tx packed in one block)
func (bc *BlockChain) FindUTXOsAndAmountsOf(senderPubKeyHash []byte, amount float64) ([]UTXO, float64) {
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
				if bytes.Equal(txOutput.PubKeyHash, senderPubKeyHash) {
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
					if bytes.Equal(utils.GetPubKeyHashFromPubKey(txInput.PubKey), senderPubKeyHash) {
						utils.AddIntElemIntoStringMap(txIDOutputIndexListMap, string(txInput.TXId), txInput.TXOutputIndex)
					}
				}
			}
		}
	}
	return utxos, current
}

func (bc *BlockChain) FindOutputByTXIdAndIndex(txId []byte, txOutputIndex int64) *TXOutput {
	for block := range bc.Iterator() {
		for _, blockTx := range block.Txs {
			if bytes.Equal(blockTx.TXId, txId) {
				return &blockTx.TXOutputs[txOutputIndex]
			}
		}
	}
	return nil
}

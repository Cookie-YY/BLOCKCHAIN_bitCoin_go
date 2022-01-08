package blocks

import (
	"blcokChain/utils"
	. "blcokChain/utils/consts"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Block is the definition of a block in block-chain
// notice:
//   - Hash: we add the hash of current block in the struct
//  		while it should be calculated by the node in block-chain every time
//   - Difficulty: the ceil hash value should be calculated by difficulty. just for demo here
//   - Data: It will be replaced by transaction(Used in Stage2. Stage3 has already replaced it with txs)
//   - Txs: Transactions of this block.
//          When hash the block, it should not be included.
//          get the Txs through miner's pack
type Block struct {
	Version    uint64 // 00 for main-chain 01 for test-chain
	MerkelRoot []byte // to represent all the txs
	TimeStamp  uint64
	Difficulty uint64 // to calculate the ceil hash value: no bigger than this value
	Nonce      uint64 // produce in mining procedure
	PrevHash   []byte // the hash of prev block
	Hash       []byte // it shouldn't be here. Just for demo
	//Data       []byte // use byte to save data
	Txs []*Transaction
}

// NewBlock : get a block(it could be GenesisBlock or commonBlock)
//  - GenesisBlock: It doesn't need to be mined. It gets its hash though simple calculation
//  - CommonBlock: It should be mined. It gets its hash though find the proper nonce.
//  NOTICE:
// 		block-chain: it needs block-chain: when verify txs, it needs to scan the whole block-chain, when creating a GenesisBlock, it could pass nil
//      prevBlockHash: the hash of prev block may not be the tailHash of the block-chain.
//     		The block need to choose the tailHash represented the longest chain as prevHash. And replace that tailHash
func NewBlock(txs []*Transaction, prevBlockHash []byte, bc *BlockChain) *Block {
	block := Block{
		Version:    00,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0, // just for demo
		Nonce:      0, // It will be filled after mining
		PrevHash:   prevBlockHash,
		Hash:       []byte{}, // It will be filled after calculation
		//Data:       []byte(data), // byte stream
		//Txs: txs,
	}
	if len(prevBlockHash) != 0 { // common block need to mine
		block.Txs = block.packTxs(txs, bc) // TODO: 打包交易和挖矿的顺序
		if len(block.Txs) <= 1 {           // no block if only coinBaseTX
			return nil
		}
		hash, nonce := block.mining()
		block.Hash = hash
		block.Nonce = nonce
	} else { // prevBlockHash == nil -> no prev block -> the first block doesn't need to mine
		block.Txs = txs
		block.Hash = block.getBlockHash()
	}
	block.MerkelRoot = block.getMerkelRoot()

	return &block
}

// packTxs : package the txs in the txs pool
//   1. chose the highest tx fee of all the txs(we have not considered the fee for now)
//   2. verify the sig of each tx.
//   TODO: 3. when packing txs, it needs to modify the uxtos' pool in block-chain.
func (b *Block) packTxs(txs []*Transaction, bc *BlockChain) []*Transaction {
	// FIXME: not all txs can be included in this block, makeTransactions in the same block may go wrong
	verifiedTxs := b.verifyTxs(txs, bc)
	return verifiedTxs
}

func (b *Block) verifyTxs(txs []*Transaction, bc *BlockChain) []*Transaction {
	verifiedTxs := make([]*Transaction, 0)
	for _, tx := range txs {
		if tx.IsCoinBaseTX() || tx.Verify(bc) { // coinBaseTX doesn't need to verify
			verifiedTxs = append(verifiedTxs, tx)
		} else {
			bytes, _ := json.Marshal(tx)
			fmt.Printf(TXVerifyFailed, string(bytes))
		}
	}
	return verifiedTxs
}

// getMerkelRoot: get MerkelRoot from txId(simply concat the transactionId to get bytes to hash)
func (b *Block) getMerkelRoot() []byte {
	txIdBytes := make([]byte, 0)
	for _, tx := range b.Txs {
		txIdBytes = append(txIdBytes, tx.TXId...)
	}
	hash := sha256.Sum256(txIdBytes)
	return hash[:]
}

// getBlockInfoToCalculateHash: when calculate hash, it will not consider real transactions.
//    Only consider MerkelRoot to represent all txs
func (b *Block) getBlockInfoToCalculateHash() *Block {
	// !! afterCopy := *b   doesn't work.
	// there is slice in the block struct. when modify the slice of copy the original could also be modified
	afterCopy := &Block{}
	err := utils.DeepCopy(afterCopy, b)
	if err != nil {
		log.Fatal(err)
	}
	afterCopy.Txs = nil
	return afterCopy
}

// getBlockHash: calculate and return the block's hash, only used in mining and genesisBlock
func (b *Block) getBlockHash() []byte {
	waitForHash := b.getBlockInfoToCalculateHash()
	blockInfo := []byte(fmt.Sprintf("%v", *waitForHash))
	hash := sha256.Sum256(blockInfo)
	return hash[:]
}

// mining: trying to find a nonce, which enable the block' hash smaller than(by bigInt) the ceil hash
func (b *Block) mining() ([]byte, uint64) {
	var nonce uint64 = 0
	for {
		b.Nonce = nonce // trying to find a nonce ==> hash(block header + nonce) < ceilHash
		currentHash := b.getBlockHash()
		currentHashBigIntWrapper := utils.GetBigIntWrapperFromBytes(currentHash)
		ceilHashBigIntWrapper := b.getCeilHashBigIntWrapper()
		if currentHashBigIntWrapper.SmallerThanAnotherBigIntWrapper(ceilHashBigIntWrapper) { // if current hash < ceil hash
			fmt.Printf(MineTip, currentHash, nonce, ceilHashBigIntWrapper.Value)
			return b.getBlockHash(), nonce
		} else {
			nonce++
		}
	}
}

// getCeilHashBigIntWrapper: it should be calculated by `difficulty` for demo, it is a const
// it can be changed to modify the difficulty
func (b *Block) getCeilHashBigIntWrapper() *utils.BigIntWrapper {
	return utils.GetBigIntWrapperFromStr(MineCeilHash, 16)
}

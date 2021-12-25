package blocks

import (
	"blcokChain/utils"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

// Block is the definition of a block in block-chain
// notice:
//   - Hash: we add the hash of current block in the struct
//  		while it should be calculated by the node in block-chain every time
//   - Difficulty: the ceil hash value should be calculated by difficulty. just for demo here
//   - Data: It will be replaced by transaction
type Block struct {
	Version    uint64 // 00 for main-chain 01 for test-chain
	MerkelRoot []byte
	TimeStamp  uint64
	Difficulty uint64 // to calculate the ceil hash value: no bigger than this value
	Nonce      uint64 // produce in mining procedure
	PrevHash   []byte // the hash of prev block
	Hash       []byte // it shouldn't be here. Just for demo
	Data       []byte // use byte to save data
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0, // just for demo
		Nonce:      0, // It will be filled after mining
		PrevHash:   prevBlockHash,
		Hash:       []byte{},     // It will be filled after calculation
		Data:       []byte(data), // byte stream
	}
	if len(prevBlockHash) != 0 { // common block need to mine
		hash, nonce := block.mining()
		block.Hash = hash
		block.Nonce = nonce
	} else { // prevBlockHash == nil -> no prev block -> the first block doesn't need to mine
		block.Hash = block.getBlockHash()
	}

	return &block
}

// getBlockHash: calculate and return the block's hash, only used in mining and genesisBlock
func (b *Block) getBlockHash() []byte {
	blockInfo := []byte(fmt.Sprintf("%v", *b))
	hash := sha256.Sum256(blockInfo)
	return hash[:]
}

// mining: trying to find a nonce, which enable the block' hash smaller than(by bigInt) the ceil hash
func (b *Block) mining() ([]byte, uint64) {
	var nonce uint64 = 0
	for {
		b.Nonce = nonce
		currentHash := b.getBlockHash()
		currentHashBigIntWrapper := utils.GetBigIntWrapperFromBytes(currentHash)
		ceilHashBigIntWrapper := b.getCeilTargetHashBigIntWrapper()
		if currentHashBigIntWrapper.SmallerThanAnotherBigIntWrapper(ceilHashBigIntWrapper) { // if current hash < ceil hash
			fmt.Printf("** Mining Sucess! hash: %x, nonce: %v, ceil hash: %x ** \n",
				currentHash, nonce, ceilHashBigIntWrapper.Value)
			return b.getBlockHash(), nonce
		} else {
			nonce++
		}
	}
}

// getCeilTargetHash: it should be calculated by `difficulty`
func (b *Block) getCeilTargetHashBigIntWrapper() *utils.BigIntWrapper {
	targetStr := "00011" + strings.Repeat("0", 59)
	return utils.GetBigIntWrapperFromStr(targetStr, 16)
}
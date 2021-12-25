package blocks

import (
	"fmt"
	"strings"
	"testing"
)

// Test_addBlock: test add block in a block-chain
func Test_addBlock(t *testing.T) {
	bc := GetOrCreateBlockChain()
	bc.AddBlock("hello-world-1")
	bc.AddBlock("hello-world-2")
	bc.AddBlock("hello-world-3")
	fmt.Print("\n######################### Current Block Chain: ######################### \n\n")
	for block := range bc.Iterator() {
		//fmt.Printf("current block height: %v\n", index)
		fmt.Printf("current block hash: %x\n", block.Hash) // []byte -> hexadecimal
		fmt.Printf("current block nonce: %d\n", block.Nonce)
		fmt.Printf("current block data: %s\n", block.Data)  // []byte -> string
		fmt.Printf("prev block hash: %x\n", block.PrevHash) // []byte -> hexadecimal

		fmt.Print("" + strings.Repeat("=", 100) + "\n")
	}

}

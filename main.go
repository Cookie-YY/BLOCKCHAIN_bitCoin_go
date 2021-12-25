package main

import (
	"blcokChain/blocks"
	"fmt"
	"strings"
)

/* Stage1: add blocks in block-chain
1. init the block-chain: it will save to db
2. add some blocks: use data string to imitate transaction in real chain.
3. for range the blocks: through iterator to print every block.
*/
func main() {
	bc := blocks.GetOrCreateBlockChain()
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

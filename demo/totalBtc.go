package main

import "fmt"

// To demonstrate the total number of bitcoin
//
// 1. Firstly, 50 bitcoins/block
// 2. after 21w blocks, bitcoins become half
const (
	firstRewordBTC = 50.0
	halfInterval   = 210000
)

func main() {
	totalBTC := 0.0
	currentRewordBTC := firstRewordBTC
	numOfBatch := 0
	for currentRewordBTC > 0 {
		numOfBatch += 1
		batchBTC := halfInterval * currentRewordBTC
		currentRewordBTC /= 2 // currentReword = currentReword / 2
		totalBTC += batchBTC
		fmt.Printf("numOfBatch: %v, currentRewordBTC%v, totalBTC: %v\n", numOfBatch, currentRewordBTC, totalBTC)
	}
	fmt.Printf("======= Conclusion =======\nnumOfBatch: %v, currentRewordBTC%v, totalBTC: %v\n", numOfBatch, currentRewordBTC, totalBTC)
	// numOfBatch: 1072, currentRewordBTC9.9e-322, totalBTC: 2.1e+07     2.1*10^7 = 2100w BTC
}

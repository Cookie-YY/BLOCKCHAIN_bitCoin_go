package main

import (
	"blcokChain/blocks"
	"blcokChain/utils"
	"fmt"
	"log"
	"strings"
)

/*
handler: defines several operations in the block-chain
	- makeTransaction : from ---amount---> to with miner
	- makeTransactions :  several  from ---amount---> to with one miner which could reduce the times of mining
	- getBalanceOf : get the rest of money of the given address
*/

// godAddress : The address of god, it can get the first coinbase transaction of the first block
const godAddress = "cookie"

var bc = blocks.GetOrCreateBlockChain(godAddress)

// makeTransaction : transfer money from ... to ...
//   notice: there will be only two transaction: coinBase and the given one
func makeTransaction(from, to string, amount float64, miner string) {
	fmt.Printf("【Running】: address [%v] want to transfer addresss [%v], amount: [%v]\n", from, to, amount)
	coinBaseTX := blocks.NewCoinBaseTX(miner, fmt.Sprintf("%v's pool name", miner))
	tx, err := blocks.NewTX(from, to, amount, bc)
	if err != nil {
		log.Fatalf("【Error】: address [%v] doesn't have enough amount of money to transfer", from)
		return
	}
	bc.AddBlock([]*blocks.Transaction{coinBaseTX, tx})

	// result
	fmt.Printf("【Done】: address [%v] sucessfully transfered money to addresss [%v], amount: [%v]\n", from, to, amount)
	fmt.Printf("\nResult: \n")
	getBalanceOf(from)
	getBalanceOf(to)
	fmt.Printf(strings.Repeat("=", 100) + "\n")
}

// makeTransactions : can have more than one transaction in one block(Reduce the times of mining)
func makeTransactions(transferMapList []map[string]string, miner string) {
	// check params
	for _, transferMap := range transferMapList {
		_, fromOk := transferMap["from"]
		_, toOk := transferMap["to"]
		amount, amountOk := transferMap["amount"]
		_, err := utils.Str2Float(amount)
		if !(fromOk && toOk && amountOk) || err != nil {
			log.Fatalf("【Error】 : usage of makeTransactions: transferMap: {'from': 'demo1', 'to': 'demo2', amount: '100'}")
		}
	}
	// make transactions
	txs := make([]*blocks.Transaction, 0)
	txs = append(txs, blocks.NewCoinBaseTX(miner, fmt.Sprintf("%v's pool name", miner)))
	for _, transferMap := range transferMapList {
		// get params
		from := transferMap["from"]
		to := transferMap["to"]
		amount, _ := utils.Str2Float(transferMap["amount"])
		fmt.Printf("【Running】: address [%v] want to transfer addresss [%v], amount: [%v]\n", from, to, amount)

		// transaction
		tx, err := blocks.NewTX(from, to, amount, bc)
		if err != nil {
			log.Fatalf("【Error】: address [%v] doesn't have enough amount of money to transfer", from)
		}
		txs = append(txs, tx)
	}
	// submit result
	bc.AddBlock(txs)

	// report result
	for _, transferMap := range transferMapList {
		// get params
		from := transferMap["from"]
		to := transferMap["to"]
		amount, _ := utils.Str2Float(transferMap["amount"])
		fmt.Printf("【Done】: address [%v] sucessfully transfered money to addresss [%v], amount: [%v]\n", from, to, amount)
		fmt.Printf(strings.Repeat("-", 100) + "\n")
	}
}

// getBalanceOf : query for the rest of money of the given address
func getBalanceOf(address string) float64 {
	money := bc.GetBalanceOf(address)
	fmt.Printf("【Result】: address [%v] has amount: %v\n", address, money)
	return money
}

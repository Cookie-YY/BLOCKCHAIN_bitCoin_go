package client

import (
	"blcokChain/blocks"
	"blcokChain/utils"
	. "blcokChain/utils/consts"
	"blcokChain/wallets"
	"fmt"
	"log"
)

/*
CMDCli: Wraps several operations in the block-chain
	- addWallets : add given num of wallets in to db
	- listAddresses : list all the addresses in the wallets
	- getOrCreateBlockChain : using for init the block-chain
	- makeTransaction : from ---amount---> to with miner
	- makeTransactions :  several  from ---amount---> to with one miner which could reduce the times of mining
	- getBalanceOf : get the rest of money of the given address
*/

func (c *CMDCli) AddWallets(num int64) {
	ws, err := wallets.GetOrCreateEmptyWallets()
	if err != nil {
		log.Fatal(err)
	}
	err = ws.AddWallets(num)
	if err != nil {
		log.Fatal(err)
	}

}

func (c *CMDCli) ListAddresses() []string {
	ws, err := wallets.GetOrCreateEmptyWallets()
	if err != nil {
		log.Fatal(err)
	}
	addresses := ws.ListAllAddresses()
	fmt.Println(addresses)
	return addresses
}

func (c *CMDCli) GetOrCreateBlockChain(godAddress string) {
	bc, err := blocks.GetOrCreateBlockChain(godAddress)
	if err != nil {
		log.Fatal(err)
	}
	c.bc = bc
}

// MakeTransaction : transfer money from ... to ...
// NOTICE: there will be only two transaction: coinBase and the given one
func (c *CMDCli) MakeTransaction(from, to string, amount float64, miner string) {
	c.GetOrCreateBlockChain("")
	// check address
	checkAddress(from, to, miner)
	fmt.Printf(CMDCliTipBeforeTransfer, from, to, amount)
	err := blocks.NewTXAndAndBlock(from, to, amount, miner, c.bc)
	if err != nil {
		log.Fatal(err)
	}
	// result
	fmt.Printf(CMDCliTipAfterTransfer, from, to, amount)
	c.GetBalanceOf(from)
	c.GetBalanceOf(to)
	fmt.Print(CMDCliTipBlockBoundary)
}

// GetBalanceOf : query for the rest of money of the given address
func (c *CMDCli) GetBalanceOf(address string) float64 {
	c.GetOrCreateBlockChain("")
	// check address
	checkAddress(address)
	// query balance
	if c.bc == nil {
		log.Fatal(ErrBlockChainNotInit())
	}
	money := c.bc.GetBalanceOf(address)
	fmt.Printf(CMDCliTipRestOfMoney, address, money)
	return money
}

// makeTransactions : can have more than one transaction in one block(Reduce the times of mining)
// NOTICE: it is unavailable now !!
func (c *CMDCli) makeTransactions(transferMapList []map[string]string, miner string) {
	if c.bc == nil {
		log.Fatal(ErrBlockChainNotInit())
	}
	// check params
	for _, transferMap := range transferMapList {
		from, fromOk := transferMap["from"]
		to, toOk := transferMap["to"]
		amount, amountOk := transferMap["amount"]
		_, err := utils.Str2Float(amount)
		if !(fromOk && toOk && amountOk) || err != nil {
			log.Fatalf(CMDCliUsageMakeTransactions)
		}
		checkAddress(from, to, miner)
	}
	// make transactions
	txs := make([]*blocks.Transaction, 0)
	txs = append(txs, blocks.NewCoinBaseTX(miner, fmt.Sprintf(MinerMsg, miner)))
	for _, transferMap := range transferMapList {
		// get params
		from := transferMap["from"]
		to := transferMap["to"]
		amount, _ := utils.Str2Float(transferMap["amount"])
		fmt.Printf(CMDCliTipBeforeTransfer, from, to, amount)

		// transaction
		tx, err := blocks.NewTX(from, to, amount, c.bc)
		if err != nil {
			log.Fatal(err)
		}
		txs = append(txs, tx)
	}
	// submit result
	err := c.bc.AddBlock(txs)
	if err != nil {
		log.Fatal(err)
	}

	// report result
	for _, transferMap := range transferMapList {
		// get params
		from := transferMap["from"]
		to := transferMap["to"]
		amount, _ := utils.Str2Float(transferMap["amount"])
		fmt.Printf(CMDCliTipAfterTransfer, from, to, amount)
		fmt.Print(CMDCliTipTXBoundary)
	}
}

func checkAddress(addresses ...string) {
	if pass, InvalidAddress := wallets.IsValidAddress(addresses); !pass {
		log.Fatal(ErrAddressInvalid(InvalidAddress))
	}
}

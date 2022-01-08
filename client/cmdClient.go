package client

import (
	"blcokChain/blocks"
	"blcokChain/utils"
	"log"
	"os"
	"strings"
)

/* const here defines some msg showing in the cmd*/
const (
	CMDCliUsageMakeTransactions = "【usage】：makeTransactions: transferMap: {'from': 'demo1', 'to': 'demo2', amount: '100'} 【Warn: no available for now】"
	CMDCliUsageMakeTransaction  = "【usage】：makeTransaction [from] [to] [amount] [miner]"
	CMDCliUsageAddWallets       = "【usage】：addWallets <num>                 the num param is optional"
	CMDCliUsageInitBlockChain   = "【usage】：initBlockChain [address]     must give a god address.It will receive the fist reward of the first block So you must run addWallets before init block chain"
	CMDCliUsageListAddress      = "【usage】：listAddress                     list all the address in the wallets"
	CMDCliUsageGetBalanceOf     = "【usage】：getBalanceOf [address]         get the balance of the given address"
	CMDCliUsageGetBalances      = "【usage】：getBalances                    get the balance of all addresses"

	CMDCliTipWelCome = `
========= Welcome to the demo block-chain =========
  @author: cookie
  @course: Block Chain of PKU
  @date: 2021-11 ~ 2022-01
	`
	CMDCliTipBeforeTransfer = "【Running】: address [%v] want to transfer address [%v], amount: [%v]\n"
	CMDCliTipAfterTransfer  = "【Done】: address [%v] successfully transferred money to address [%v], amount: [%v]\n"
	CMDCliTipRestOfMoney    = "【Result】: address [%v] has amount: %v\n"
	CMDCliTipBlockBoundary  = "====================================================================================================\n"
	CMDCliTipTXBoundary     = "----------------------------------------------------------------------------------------------------\n"
)

var CMDCliUsage = "\n" + strings.Join([]string{CMDCliTipWelCome, CMDCliUsageMakeTransactions, CMDCliUsageMakeTransaction, CMDCliUsageAddWallets, CMDCliUsageInitBlockChain, CMDCliUsageListAddress, CMDCliUsageGetBalanceOf, CMDCliUsageGetBalances}, "\n")

// CMDCli : Wraps several operations in the block-chain. Using in the cmd-line
type CMDCli struct {
	bc *blocks.BlockChain
}

func GetCMDCli() *CMDCli {
	return &CMDCli{}
}

func (c *CMDCli) Run() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal(CMDCliUsage)
	}

	op, params := args[1], args[2:]
	switch op {
	case "addWallets", "aw":
		var num int64 = 1
		var err error = nil
		if len(params) == 1 {
			num, err = utils.Str2Int(params[0])
			if err != nil {
				log.Fatal(CMDCliUsage)
			}
		}
		c.AddWallets(num)
	case "listWallets", "lw":
		//c.AddWallet()
	case "listBlockChain", "lbc":
		//c.AddWallet()
	case "listAddress", "la":
		c.ListAddresses()
	case "getBalanceOf", "gbo":
		if len(params) != 1 {
			log.Fatal(CMDCliUsage)
		}
		c.GetBalanceOf(params[0])
	case "getBalances", "gbs":
		for _, address := range c.ListAddresses() {
			c.GetBalanceOf(address)
		}
	case "makeTransaction", "mt":
		if len(params) != 4 {
			log.Fatal(CMDCliUsage)
		}
		from, to, amountStr, miner := params[0], params[1], params[2], params[3]
		amount, err := utils.Str2Float(amountStr)
		if err != nil {
			log.Fatal(CMDCliUsage)
		}
		c.MakeTransaction(from, to, amount, miner)
	case "initBlockChain", "initbc":
		if len(params) != 1 {
			log.Fatal(CMDCliUsage)
		}
		c.GetOrCreateBlockChain(params[0])
	default:
		log.Fatal(CMDCliUsage)
	}

}

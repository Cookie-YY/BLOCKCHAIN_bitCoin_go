package consts

import (
	"fmt"
)

type ErrorWrapper struct {
	ErrMsg  string
	ErrCode string
}

// Error : any struct which implement this method can be regarded as err
//   the error is an interface
func (ew *ErrorWrapper) Error() string {
	return ew.ErrMsg
}

func NewError(errMsg, ErrCode string) *ErrorWrapper {
	return &ErrorWrapper{errMsg, ErrCode}
}

var (
	// Wallet && Money
	ErrAddressInvalid = func(address string) *ErrorWrapper {
		return NewError(fmt.Sprintf("【Error】: address: [%v] is invalid, please check", address), "")
	}
	ErrNoEnoughMoney = func(address string, amount float64) *ErrorWrapper {
		return NewError(fmt.Sprintf("address: [%v] doesn't have enough money to transfer [%v]", address, amount), "")
	}
	ErrAddressNotFound = func(address string) *ErrorWrapper {
		return NewError(fmt.Sprintf("Error: the address: [%v] hasn't been registered", address), "")
	}
	// Block && BlockChain
	ErrBlockGetInBlockChainFailed = func() *ErrorWrapper {
		return NewError(fmt.Sprintf("Error: the block gets in block-chain failed. Maybe all txs verified failed"), "")
	}
	ErrBlockChainNotInit = func() *ErrorWrapper {
		return NewError(fmt.Sprintf("Error: the block chain is not init: [initBlockChain] first"), "")
	}
	ErrBlockChainNoGod = func() *ErrorWrapper {
		return NewError(fmt.Sprintf("Error: when init the block chain, it needs a god address. Do you forget to init the block chain?"), "")
	}
	// Transaction
	ErrTXSignRefereedOutput = func() *ErrorWrapper {
		return NewError(fmt.Sprintf("Error: when sign the transaction, the length of inputs and the length of the refereed output's pubKeyHash mush be equal"), "")
	}
)

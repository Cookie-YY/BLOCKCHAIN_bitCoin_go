package blocks

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

const reward = 12.5

type Transaction struct {
	TXId      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXId          []byte
	TXOutputIndex int64 // -1 means coinbaseTX's input
	Sig           string
}

type TXOutput struct {
	Amount     float64
	PubKeyHash string
}

// UTXO : is the wrapper of TXOutput
type UTXO struct {
	TXOutput
	TXId        []byte
	OutputIndex int64 // -1 means coinbaseTX's input
}

// getTXHash : to generate the id of transaction
func (tx *Transaction) getTXHash() []byte {
	blockInfo := []byte(fmt.Sprintf("%v", *tx))
	hash := sha256.Sum256(blockInfo)
	return hash[:]
}

// IsCoinBaseTX : coinBaseTX can bypass the input check(no corresponding output)
func (tx *Transaction) IsCoinBaseTX() bool {
	if len(tx.TXInputs) == 1 && len(tx.TXOutputs) == 1 && tx.TXInputs[0].TXOutputIndex == -1 {
		return true
	}
	return false
}

// NewCoinBaseTX : coinBaseTX can carry some message. Because it doesn't need to offer signature
func NewCoinBaseTX(address string, data string) *Transaction {
	txInput := TXInput{[]byte{}, -1, data}
	txOutput := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{txInput}, []TXOutput{txOutput}}
	tx.TXId = tx.getTXHash()
	return &tx
}

// NewTX : it records every common transactions
func NewTX(from string, to string, amount float64, bc *BlockChain) (*Transaction, error) {
	utxos, amounts := bc.FindUTXOsAndAmountsOf(from, amount)
	if amounts < amount {
		return nil, errors.New(fmt.Sprintf("address: [%v] doesn't have enough money", from))
	}
	txInputs := make([]TXInput, 0)
	txOutputs := make([]TXOutput, 0)

	for _, utxo := range utxos {
		txInput := TXInput{utxo.TXId, utxo.OutputIndex, from}
		txInputs = append(txInputs, txInput)
	}

	txOutput := TXOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)

	if amounts > amount {
		txOutputs = append(txOutputs, TXOutput{amounts - amount, from})
	}
	tx := Transaction{[]byte{}, txInputs, txOutputs}
	tx.TXId = tx.getTXHash()

	return &tx, nil
}

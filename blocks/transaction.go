package blocks

import (
	"blcokChain/utils"
	. "blcokChain/utils/consts"
	"blcokChain/wallets"
	"crypto/sha256"
	"fmt"
	"log"
)

type Transaction struct {
	TXId      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXId          []byte
	TXOutputIndex int64  // -1 means coinbaseTX's input
	Sig           []byte // r+s
	PubKey        []byte // x+y: means the pubKey of the sender
}

type TXOutput struct {
	Amount     float64
	PubKeyHash []byte // convert from receiver's address
}

func NewTXOutput(amount float64, receiverAddress string) *TXOutput {
	return &TXOutput{amount, utils.GetPubKeyHashFromAddress(receiverAddress)}
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
	txInput := TXInput{[]byte{}, -1, nil, []byte(data)}
	txOutput := NewTXOutput(MinerReward, address)
	tx := Transaction{[]byte{}, []TXInput{txInput}, []TXOutput{*txOutput}}
	tx.TXId = tx.getTXHash()
	return &tx
}

// NewTX : it records every common transactions
func NewTX(from string, to string, amount float64, bc *BlockChain) (*Transaction, error) {
	// 1. find wallet by address in wallets
	ws, err := wallets.GetOrCreateEmptyWallets()
	if err != nil {
		return nil, err
	}
	wallet, err := ws.GetWalletOf(from)
	if err != nil {
		return nil, err
	}

	// 2. find pubKey and privateKey
	pubKey, privateKey := wallet.PubKey, wallet.PrivatKey

	// 3. find available money: need pubKeyHash. Because it needs pubKeyHash in output to match
	// input's pubKey -> pubKeyHash to match pubKeyHash in output
	utxos, amounts := bc.FindUTXOsAndAmountsOf(utils.GetPubKeyHashFromPubKey(pubKey), amount)
	if amounts < amount {
		return nil, ErrNoEnoughMoney(from, amount)
	}

	// 4. create inputs and outputs
	txInputs := make([]TXInput, 0)
	txOutputs := make([]TXOutput, 0)
	refereedOutputPubKeyHashList := make([][]byte, 0)
	for _, utxo := range utxos {
		txInput := TXInput{utxo.TXId, utxo.OutputIndex, nil, pubKey}
		txInputs = append(txInputs, txInput)
		refereedOutputPubKeyHashList = append(refereedOutputPubKeyHashList, utxo.PubKeyHash)
	}
	txOutputs = append(txOutputs, *NewTXOutput(amount, to))

	// 5. return to sender some money if needed
	if amounts > amount { // needed utxo > amount
		txOutputs = append(txOutputs, *NewTXOutput(amounts-amount, from))
	}
	tx := Transaction{[]byte{}, txInputs, txOutputs}
	tx.TXId = tx.getTXHash()

	// 6. sign the transaction.
	// - from: all the outputs' pubKeyHash referred by the input
	// - amount: the amount
	// - to: the sender address
	err = tx.Sign(privateKey, pubKey, refereedOutputPubKeyHashList)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// getTXTemplate : get a copy of the current transaction to store all the infos which are needed to be signed
// store in this kind of way is not very clear. Because pubKey in input is the sender's pubKey now it will contain the refereed output's pubKeyHash
// you can use another defined struct. But it's more complex than this.
func (tx *Transaction) getTXTemplate() *Transaction {
	afterCopy := &Transaction{}
	err := utils.DeepCopy(afterCopy, tx)
	if err != nil {
		log.Fatal(err)
	}
	for index := range afterCopy.TXInputs {
		afterCopy.TXInputs[index].Sig = nil    // it will be filled after load all the key factors.
		afterCopy.TXInputs[index].PubKey = nil // it will be filled by the refereed output's pubKeyHash
	}
	return afterCopy
}

// Sign : sign the three key factors:
//   - from: refereedOutputPubKeyHashList
//   - amount: current amount in output
//   - to: current pubKey in output
//   NOTICE: the sign process should not need pubKey, but we need pubKey to restore ecdsaPrivateKey from privateKey
func (tx *Transaction) Sign(privateKey, pubKey []byte, refereedOutputPubKeyHashList [][]byte) error {
	if len(tx.TXInputs) != len(refereedOutputPubKeyHashList) { // the length of inputs and the length of the refereed output's pubKeyHash mush be equal
		return ErrTXSignRefereedOutput()
	}
	// templateTXForSignature as a container to load all the key factors
	// amount and to_address have already in the output.
	// We need to load the 「refereedOutputPubKeyHashList」 in the 「templateTXForSignature」
	templateTXForSignature := tx.getTXTemplate()
	for index := range refereedOutputPubKeyHashList {
		// 1. prepare hash data to sign. It's a hash data, using the TXID as the data
		templateTXForSignature.TXInputs[index].PubKey = refereedOutputPubKeyHashList[index] // load pubKeyHash
		hashForSignature := templateTXForSignature.getTXHash()                              // all the required infos are in this templateTX
		// 2. restore the privateKey
		ecdsaPrivateKey := utils.GetECDSAPrivateKeyFromPrivateKey(privateKey, pubKey)
		// 3. sign
		sig, err := utils.Sign(ecdsaPrivateKey, hashForSignature)
		if err != nil {
			return err
		}
		tx.TXInputs[index].Sig = sig
	}
	return nil
}

// Verify : need to verify every signature of the input in the given tx.
// refereed pubKeyHashList:
// 	- the hard-to-find refereed pubKeyHashList need to be found after scan the whole block-chain to check
// 	- actually, when signing, the refereed pubKeyHashList are also found after scan the whole block-chain.
// pubKey: from each input
// signature: from each input
func (tx *Transaction) Verify(bc *BlockChain) bool {
	templateTXForVerify := tx.getTXTemplate()
	for index, txInput := range tx.TXInputs {
		// 1. scan the block-chain to find refereed pubKeyHashList
		targetTXOutput := bc.FindOutputByTXIdAndIndex(txInput.TXId, txInput.TXOutputIndex)
		if targetTXOutput == nil {
			return false // something goes wrong, but no need to return error
		}
		templateTXForVerify.TXInputs[index].PubKey = targetTXOutput.PubKeyHash
		// 2. prepare 3 key things: hashData, signature, pubKey
		hashForVerify := templateTXForVerify.getTXHash()
		pubKey := utils.GetECDSAPubKeyFromPubKey(txInput.PubKey)
		sig := txInput.Sig
		// 3. verify
		if !utils.Verify(&pubKey, hashForVerify, sig) {
			return false
		}
	}
	return true
}

// NewTXAndAndBlock : create transaction and record it in a block
//  it means every block only record one common transaction. Also means every transaction needs a miner
func NewTXAndAndBlock(from string, to string, amount float64, miner string, bc *BlockChain) error {
	coinBaseTX := NewCoinBaseTX(miner, fmt.Sprintf(MinerMsg, miner))
	tx, err := NewTX(from, to, amount, bc)
	if err != nil {
		return err
	}
	err = bc.AddBlock([]*Transaction{coinBaseTX, tx})
	if err != nil {
		return err
	}
	return nil
}

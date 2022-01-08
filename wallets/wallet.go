package wallets

import (
	"blcokChain/utils"
	"bytes"
)

type Wallet struct {
	PrivatKey []byte `json:"privat_key_string"`
	PubKey    []byte `json:"pub_key_string"`
}

func newWallet() (*Wallet, error) {
	privatKey, pubKey, err := utils.NewKeyPair()
	if err != nil {
		return nil, err
	}
	return &Wallet{privatKey, pubKey}, nil
}

func IsValidAddress(addressList []string) (bool, string) {
	for _, address := range addressList {
		if !bytes.Equal(
			utils.GetCheckSumFromPayLoad(utils.GetPayloadFromAddress(address)),
			utils.GetAddressCheckSumFromAddress(address),
		) {
			return false, address
		}
	}
	return true, ""
}

func (w *Wallet) getAddressFromPubKey() string {
	return utils.GetAddressFromPubKey(w.PubKey)
}

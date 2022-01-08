package wallets

import (
	"blcokChain/utils"
	. "blcokChain/utils/consts"
)

type Wallets struct {
	WalletMap map[string]*Wallet // Address -> Wallet
}

func GetOrCreateEmptyWallets() (*Wallets, error) {
	jdb := utils.GetJsonDB(JsonFileNameWallets, utils.NewJsonSerializer())
	ws := &Wallets{}
	ws.WalletMap = make(map[string]*Wallet)
	if jdb.IsExist() {
		err := jdb.ReadFromDB(ws)
		if err != nil {
			return nil, err
		}
		return ws, nil
	}
	err := jdb.WriteToDB(ws)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (ws *Wallets) addWallet(db utils.DB) error {
	w, err := newWallet()
	if err != nil {
		return err
	}

	ws.WalletMap[w.getAddressFromPubKey()] = w
	err = db.WriteToDB(ws)
	if err != nil {
		return err
	}
	return nil
}

func (ws *Wallets) AddWallets(num int64) error {
	var err error
	jdb := utils.GetJsonDB(JsonFileNameWallets, utils.NewJsonSerializer())
	for i := 0; i < int(num); i++ {
		err = ws.addWallet(jdb)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ws *Wallets) GetWalletOf(address string) (*Wallet, error) {
	wallet, ok := ws.WalletMap[address]
	if ok {
		return wallet, nil
	}
	return nil, ErrAddressNotFound(address)
}

func (ws *Wallets) ListAllAddresses() []string {
	addresses := make([]string, 0)
	for address := range ws.WalletMap {
		addresses = append(addresses, address)
	}
	return addresses
}

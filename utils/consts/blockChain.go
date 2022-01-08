package consts

import "strings"

const (
	JsonFileNameWallets    = "wallets.json"
	JsonFileNameBlockChain = "block-chain.json"

	MinerReward            = 12.5
	MinerMsgOfGenesisBlock = "cookie creates it" // the msg in the coinBaseTX of the genesis block(it's the pubKey of input)
	MinerMsg               = "%v's pool name"    // the msg in the coinBaseTX of the common block(it's the pubKey of input)
	MineTip                = "** 【Mining Success】! hash: %x, nonce: %v, ceil hash: %x ** \n"
	TXVerifyFailed         = "the tx: %v, verified error\n"
)

var (
	MineCeilHash     = MineCeilHashEasy
	MineCeilHashEasy = "00010" + strings.Repeat("0", 59)
	MineCeilHashHard = "00001" + strings.Repeat("0", 59)
)

# Block-Chain-demo
> A block-chain demo go project.

- `blockChain/demo` contains some demonstration which could run independently
  - totalBTC.go: 
    - The total num of BTC: 
    `21w`
  - pow.go: 
    - The hash: `a subtle variance cause huge difference`
- To study better, I Split this project to several stages
  - Step1: use data to imitate transaction
    - block & block-chain
    - mining process
    - db to save block-chain(json file)
  - Step2: use transaction to replace data
  - Step3: use Cryptography to replace name in address
## Stage 1: 
> A simple block-chain with several blocks which contain some data.
### Notice && Techniques
- about bytes
  ~~~go
  // 1. bytes, string, hex
  []byte("hello world")        // string -> []byte
  []byte(0xa3b)                // hex -> []byte
  fmt.Sprintf("%x", hexBytes)  // []byte -> hex
  fmt.Sprintf("%s", strBytes)  // []byte -> string
  strBytes[:]                  // [size]byte -> []byte
  // 2. hash:
  sha256.Sum256()  // require input and output are []byte
  // 3. compare
  tmp := big.Int{}
  tmp.SetBytes(bigBytes)
  tmp.Cmp(tmp)
  ~~~
- go.mod
  - reason: `go test` or  `go run *.go` need an entry
  - usage: `go mod init blcokChain`
## Stage 2:
> A block-chain with transaction in the block. Still no Cryptography
### Notice && Techniques
- use UTXO to wrap output
  - the index of output is hard to get
  - maybe the best way to represent UTXO is to wrap output
- use address to imitate public key and private key
- use `handler.go` to wrap operations in block-chain.
## Stage 3:
> A block-chain with transaction in the block with Cryptography.
- Refactoring the whole project~
- add more files in demo/
  - `cryptography.go` shows Sign and Verify process using elliptic curve
  - `serializer.go` shows two serialize ways. json and gob
    - Json: easy to use and easy to read. But sometimes goes wrong
    - Gob: can handle most of the cases. But hard to read(save in bytes format)
- privateKey && pubKey && address
  - privateKey: Calculated by elliptic curve
  - pubKey: Derived from privateKey
  - pubKeyHash: Hash pubKey 2 times: sha256 + ripemd160
  - version: default 00
  - checkSum: Hash the payload and take the first 4 bytes
    - payload: concat version and pubKeyHash
    - purpose of checkSum: verify if an address is valid quickly.
  - address: base58 the payload
    - payload: concat version, puKeyHash, checkSum
    - purpose of base58: easy for human to read(base58 ignore some confusing letters compared with base64)
- sign && verify
  - sign process: sign the 3 key factors
    - from: refereedOutputPubKeyHashList
      - prepare it when scan block-chain for available utxos 
    - amount: current amount in output
    - to: current pubKey in output
  - verify process: it needs to restore the hashData
    - hashData: 
      - scan the block-chain to find refereed pubKeyHash
      - makeup data
      - hash the whole transaction
    - pubKey: save in the input
    - signature: save in the input
    
- `client/cmdClient`: wrap the Client
  - wrap some common operations into a client

## Usage: cmd：
> Based on Stage3
~~~
【usage】：./main makeTransaction [from] [to] [amount] [miner]
【usage】：./main addWallets <num>             the num param is optional.
    After 「addWallets」, you can get a json file called wallets.json
【usage】：./main initBlockChain [address]     must give a god address.It will receive the fist reward of the first block So you must run addWallets before init block chain
    After 「initBlockChain」, you can get a json file called block-chain.json
【usage】：./main listAddress                  list all the address in the wallets
【usage】：./main getBalanceOf [address]       get the balance of the given address
【usage】：./main getBalances                  get the balance of all addresses

 e.g.
 ./main addWallets 3 
 ./main initBlockChain x
 ./main makeTransaction x y a z
~~~
## TODO
### Transaction
- more than one tx in one block
  - can't consider each tx only check form block-chain.
### BlockChain
- UTXO pool
  - quickly find the utxos of the given address
- Mining
  - ceilHash
    - modify the ceilHash dynamically
    - calculate the ceilHash from difficulty
  - transaction fee
- more than one node
  - tailHash will not be one anymore
  - synchronization information between two nodes
- SPV: 
  - realize merkelRoot
  - SPV process
### Accounts && Wallets
- Mnemonic Phrase
### Project
- add WebClient and frontEnd: Go socket + react
- add GUIClient: C++

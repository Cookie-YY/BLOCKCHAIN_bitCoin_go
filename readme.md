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
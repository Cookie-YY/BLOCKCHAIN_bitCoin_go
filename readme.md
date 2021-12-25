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

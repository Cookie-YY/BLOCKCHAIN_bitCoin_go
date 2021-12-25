package main

import (
	"crypto/sha256"
	"fmt"
)

// To demonstrate that there's no regulation in encryption process.
// Which means a little change in raw data, very difference in output
//
// 1. get the different input: "hello-world-1", "hello-world-2","hello-world-2"...
// 2. encrypt and print
func main() {
	rawData := "hello-world"
	for nonce := 0; nonce < 100; nonce++ {
		waitForEncryptStr := fmt.Sprintf("%v-%v", rawData, nonce) // hello-world-{nounce}
		fmt.Printf("waitForEncryptStr: %v\n", waitForEncryptStr)
		afterEncryptBytes := sha256.Sum256([]byte(waitForEncryptStr)) // encrypt
		fmt.Printf("afterEncryptBytes: %x\n\n", afterEncryptBytes)    // 16(Hexadecimal)
	}

}

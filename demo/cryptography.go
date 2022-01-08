package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
)

// To demonstrate the cryptography in the block-chain using elliptic.
//
// 1. get the curve
// 2. get private key: using ecdsa based on the curve
// 3. get public key: based on the private key
// 4. sign: using ecdsa based on the private key to sign the data(get r and s)
//      - the data usually is the hash data
// 5. verify: using public key to verify the sign(r and s)
//      - verify you are the right person: your public key can resolve the sign(signed by the corresponding private key)
//      - verify the data has not been changed: you need know the hash of the original dada.
func main() {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	pubKey := privateKey.PublicKey
	data := "hello-world"
	hash := sha256.Sum256([]byte(data))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:]) // r, s is two parts of the signature.
	if err != nil {
		log.Fatal(err)
	}

	res := ecdsa.Verify(&pubKey, hash[:], r, s)
	fmt.Printf("verify result: %v", res)
}

package utils

import (
	"blcokChain/utils/ripemd160"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

/*
privateKey: []byte
pubKey: []byte: x+y
	x,y big.Int
address: string: []byte after base58
ecdsaPubKey: ecdsa.PublicKey
ecdsaPrivateKey: *ecdsa.PrivateKey
signature/sig: r+s
	r,s big.Int
*/

// NewKeyPair : create privateKeyString and pubKeyString
func NewKeyPair() ([]byte, []byte, error) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader) // It will go wrong when not having sufficient space
	if err != nil {
		return nil, nil, err
	}
	pubKey := privateKey.PublicKey
	pubKeyXY := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	return privateKey.D.Bytes(), pubKeyXY, nil
}

// GetAddressFromPubKey : 2 hash + version + checkSum + base58
//  Notice:
//     the step1: 2hash process are in the GetPubKeyHashFromPubKey. Here is for clearly look
//     the step3: calculate checksum from payload process are in the GetCheckSumFromPayLoad. Here is for clearly look
func GetAddressFromPubKey(pubKey []byte) string {
	// 1. two times of hash:  sha256 + ripemd160
	hash := sha256.Sum256(pubKey)
	hash160 := ripemd160.Ripemd160(hash[:])
	// 2. concat version
	version := byte(00)
	payload := append([]byte{version}, hash160...)
	// 3. calculate checksum: check if an address is valid in a short time (sha256 + first 4 bytes)
	tmp := sha256.Sum256(payload)
	hashForCheck := sha256.Sum256(tmp[:])
	checkSum := hashForCheck[:4]
	// 4. concat checksum
	waitForEncode := append(payload, checkSum...)
	// 5. base58: easy to read by human
	address := string(Base58Encode(waitForEncode))
	return address
}

func GetPubKeyHashFromAddress(address string) []byte {
	// version(1bytes) + pubKeyHash(20) + checkSum(4bytes)
	decodeBytes := Base58Decode([]byte(address)) // 25bytes
	return decodeBytes[1 : len(decodeBytes)-4]
}

func GetPayloadFromAddress(address string) []byte {
	// payload = version(1bytes) + pubKeyHash(20)
	decodeBytes := Base58Decode([]byte(address)) // 25bytes
	return decodeBytes[:len(decodeBytes)-4]
}

func GetAddressCheckSumFromAddress(address string) []byte {
	// version(1bytes) + payload(20) + checkSum(4bytes)
	decodeBytes := Base58Decode([]byte(address)) // 25bytes
	return decodeBytes[len(decodeBytes)-4:]
}

func GetPubKeyHashFromPubKey(pubKey []byte) []byte {
	// 1. two times of hash:  sha256 + ripemd160
	hash := sha256.Sum256(pubKey)
	return ripemd160.Ripemd160(hash[:])
}

func GetCheckSumFromPayLoad(payload []byte) []byte {
	tmp := sha256.Sum256(payload)
	hashForCheck := sha256.Sum256(tmp[:])
	return hashForCheck[:4]
}

func GetECDSAPubKeyFromPubKey(pubKey []byte) ecdsa.PublicKey {
	x := GetBigIntWrapperFromBytes(pubKey[:len(pubKey)/2])
	y := GetBigIntWrapperFromBytes(pubKey[len(pubKey)/2:])
	return ecdsa.PublicKey{Curve: elliptic.P256(), X: x.Value, Y: y.Value}
}

func GetECDSAPrivateKeyFromPrivateKey(privateKey []byte, pubKey []byte) *ecdsa.PrivateKey {
	d := big.Int{}
	d.SetBytes(privateKey)
	return &ecdsa.PrivateKey{PublicKey: GetECDSAPubKeyFromPubKey(pubKey), D: &d}
}

func getRSFromSignature(sig []byte) (*big.Int, *big.Int) {
	r := GetBigIntWrapperFromBytes(sig[:len(sig)/2])
	s := GetBigIntWrapperFromBytes(sig[len(sig)/2:])
	return r.Value, s.Value
}

// Sign : the basic sign process
func Sign(ecdsaPrivateKey *ecdsa.PrivateKey, dataForSign []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivateKey, dataForSign) // r, s is two parts of the signature.
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

func Verify(ecdsaPubKey *ecdsa.PublicKey, hashData []byte, sig []byte) bool {
	r, s := getRSFromSignature(sig)
	return ecdsa.Verify(ecdsaPubKey, hashData, r, s)
}

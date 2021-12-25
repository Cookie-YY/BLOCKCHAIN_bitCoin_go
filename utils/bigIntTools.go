package utils

import "math/big"

// BigIntWrapper : offers a way to compare two bytes(usually the hash value)
type BigIntWrapper struct {
	Value *big.Int
}

func GetBigIntWrapperFromStr(bigStr string, base int) *BigIntWrapper {
	targetInt := big.Int{}
	targetInt.SetString(bigStr, base)
	return &BigIntWrapper{&targetInt}
}

func GetBigIntWrapperFromBytes(bigBytes []byte) *BigIntWrapper {
	targetInt := big.Int{}
	targetInt.SetBytes(bigBytes)
	return &BigIntWrapper{&targetInt}
}

func (bw *BigIntWrapper) EqualToAnotherBigIntWrapper(anotherBw *BigIntWrapper) bool {
	return bw.Value.Cmp(anotherBw.Value) == 0
}

func (bw *BigIntWrapper) SmallerThanAnotherBigIntWrapper(anotherBw *BigIntWrapper) bool {
	return bw.Value.Cmp(anotherBw.Value) == -1
}

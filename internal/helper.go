package internal

import (
	"crypto/rand"
	"math/big"
)

func RandIntn(m int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(m))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}

func RandFloat64() float64 {
	return float64(RandIntn(1<<DefaultRandMantissa)) / (1 << DefaultRandMantissa)
}

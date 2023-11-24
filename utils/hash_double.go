package utils

import (
	"crypto/sha256"
	"math/big"
	"strconv"
)

func GetHashDoubleRule(visitorCode string, containerID int, respoolTime int) float64 {
	suffix := ""
	if respoolTime > 0 {
		suffix = strconv.Itoa(respoolTime)
	}
	return GetHashDouble(visitorCode, containerID, suffix)
}

func GetHashDouble(visitorCode string, containerID int, suffix ...string) float64 {
	var b []byte
	b = append(b, visitorCode...)
	b = append(b, WritePositiveInt(containerID)...)
	if len(suffix) > 0 {
		b = append(b, suffix[0]...)
	}

	h := sha256.New()
	h.Write(b)

	z := new(big.Int).SetBytes(h.Sum(nil))
	n1 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

	f1 := new(big.Float).SetInt(z)
	f2 := new(big.Float).SetInt(n1)
	f, _ := new(big.Float).Quo(f1, f2).Float64()
	return f
}

package utils

import (
	"crypto/sha256"
	"math"
	"strconv"
)

func ObtainHashRule(visitorCode string, containerID int, respoolTime int) float64 {
	suffix := ""
	if respoolTime > 0 {
		suffix = strconv.Itoa(respoolTime)
	}
	return ObtainHash(visitorCode, containerID, suffix)
}

func ObtainHash(visitorCode string, containerID int, suffix ...string) float64 {
	var b []byte
	b = append(b, visitorCode...)
	b = append(b, WritePositiveInt(containerID)...)
	if len(suffix) > 0 {
		b = append(b, suffix[0]...)
	}
	return CalculateHash(b)
}

func ObtainHashForMEGroup(visitorCode string, meGroupName string) float64 {
	return CalculateHash([]byte(visitorCode + meGroupName))
}

func CalculateHash(b []byte) float64 {
	h := sha256.New()
	h.Write(b)
	b = h.Sum(nil)
	parsedValue := uint64(b[7]) |
		(uint64(b[6]) << 8) |
		(uint64(b[5]) << 16) |
		(uint64(b[4]) << 24) |
		(uint64(b[3]) << 32) |
		(uint64(b[2]) << 40) |
		(uint64(b[1]) << 48) |
		(uint64(b[0]) << 56)
	return float64(parsedValue) / math.MaxUint64
}

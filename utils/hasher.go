package utils

import (
	"crypto/sha256"
	"encoding/binary"
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
	b := make([]byte, 0, len(visitorCode)+len(meGroupName))
	b = append(b, visitorCode...)
	b = append(b, meGroupName...)
	return CalculateHash(b)
}

func CalculateHash(b []byte) float64 {
	sum := sha256.Sum256(b)
	return float64(binary.BigEndian.Uint64(sum[:8])) / math.MaxUint64
}

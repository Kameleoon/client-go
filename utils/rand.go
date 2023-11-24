package utils

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

const letterBytes = "ABCDEF0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var random *rand.Rand
var randomMx sync.Mutex

// A random generator returned by `getRandom` function should be used
// in order to avoid deterministic value generation.
func getRandom() *rand.Rand {
	if random == nil {
		randomMx.Lock()
		defer randomMx.Unlock()
		if random == nil {
			random = rand.New(rand.NewSource(time.Now().UnixNano()))
		}
	}
	return random
}

func GetRandomString(n int) string {
	rnd := getRandom()
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	randomMx.Lock()
	for i, cache, remain := n-1, rnd.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rnd.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	randomMx.Unlock()

	return sb.String()
}

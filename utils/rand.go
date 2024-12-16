package utils

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
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
			random = rand.New(rand.NewSource(getRandomSeed()))
		}
	}
	return random
}

func getRandomSeed() int64 {
	var seedBuffer [8]byte
	if _, err := crand.Read(seedBuffer[:]); err != nil {
		logging.Error("Failed to generate random seed: %s", err)
		return time.Now().UnixNano()
	}
	return int64(binary.BigEndian.Uint64(seedBuffer[:]))
}

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// `letterBytes` should not be longer than 64 chars. All chars after this limit will be ignored.
func GetRandomString(n int, letterBytes string) string {
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

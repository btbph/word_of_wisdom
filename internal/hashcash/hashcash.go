package hashcash

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash"
	"math"
	"time"
)

type Clock interface {
	Now() time.Time
}

type Hashcash struct {
	zeroBits   int
	clock      Clock
	saltLength int
	hasher     hash.Hash
}

func NewHashcash(
	zeroBits, saltLength int,
	clock Clock,
	hasher hash.Hash,
) *Hashcash {
	return &Hashcash{
		zeroBits:   zeroBits,
		saltLength: saltLength,
		clock:      clock,
		hasher:     hasher,
	}
}

func (h Hashcash) Generate(resource string) string {
	salt := generateRandomString(h.saltLength)
	now := h.clock.Now().UTC().Format("060102") // YYMMDD
	base := fmt.Sprintf("1:%d:%s:%s::%s", h.zeroBits, now, resource, salt)

	counter := 0
	hashcashString := h.getHashcashString(base, counter) // TODO: introduce string builder
	for !checkZeroes(h.hasher, h.zeroBits, hashcashString) {
		hashcashString = h.getHashcashString(base, counter)
		counter++
	}

	return hashcashString
}

func (h Hashcash) getHashcashString(base string, counter int) string {
	return base + fmt.Sprintf(":%x", counter)
}

func checkZeroes(hasher hash.Hash, zeroBits int, str string) bool {
	hasher.Reset()
	hasher.Write([]byte(str))
	generatedHash := hasher.Sum(nil)
	hexDigits := int(math.Ceil(float64(zeroBits) / 4))

	i := 0
	for i < hexDigits {
		if generatedHash[i] != 0 {
			return false
		}
		i++
	}

	return true
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)[:length]
}

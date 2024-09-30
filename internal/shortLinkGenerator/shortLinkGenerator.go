package shortlinkgenerator

import (
	"crypto/rand"

	"github.com/decred/base58"
)

func GenerateShortLink() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base58.Encode(b)
}

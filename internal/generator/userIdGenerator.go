package generator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateUserID() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b) // записываем байты в слайс b
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	return hex.EncodeToString(b), nil
}

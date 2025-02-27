package util

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomAddress() string {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }
    return hex.EncodeToString(b)
}


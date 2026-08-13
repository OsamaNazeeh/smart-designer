package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func RefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}
package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return base64.StdEncoding.EncodeToString(hash[:])
}

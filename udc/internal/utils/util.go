package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashB64(data []byte) string {
	hash := sha256.Sum256(data)
	return base64.URLEncoding.EncodeToString(hash[:])
}

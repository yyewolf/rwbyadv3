package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateNewCookieId() string {
	b := make([]byte, 66)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
)

func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func CompareTokens(a, b string) bool {
	res := subtle.ConstantTimeCompare([]byte(a), []byte(b))
	return res == 1
}

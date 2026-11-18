package identity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func NewSecret() (string, []byte, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, fmt.Errorf("generate random secret: %w", err)
	}
	value := base64.RawURLEncoding.EncodeToString(raw)
	digest := sha256.Sum256([]byte(value))
	return value, digest[:], nil
}

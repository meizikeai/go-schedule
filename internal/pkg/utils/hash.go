// internal/pkg/utils/hash.go
package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func SHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func SHA256Bytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func SHA512(s string) string {
	sum := sha512.Sum512([]byte(s))
	return hex.EncodeToString(sum[:])
}

func SHA512Bytes(b []byte) string {
	sum := sha512.Sum512(b)
	return hex.EncodeToString(sum[:])
}

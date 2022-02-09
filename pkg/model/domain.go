package model

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Domain holds the config for dyndns related updates
type Domain struct {
	Credentials string
}

// ValidateCredentials compares supplied creds against hashed creds from config
func (d *Domain) ValidateCredentials(user string, password string) bool {
	return strings.ToUpper(d.Credentials) == strings.ToUpper(getSha2Hash(user+password))
}

func getSha2Hash(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

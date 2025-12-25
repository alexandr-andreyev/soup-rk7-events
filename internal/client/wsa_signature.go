package client

import (
	"crypto/sha256"
	"encoding/base64"
)

func (c *ExternalClient) GenerateHash(input []byte) string {
	data := append(input, c.config.WSASalt...)
	hash := sha256.Sum256((data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// CryptoService example, use AES-256 as an example of implementing AES-256 encryption.
// Encryption strength can be enhanced based on implementation.
type CryptoService struct {
	key []byte
}

func NewCryptoService() (*CryptoService, error) {
	key, err := generateKey()
	if err != nil {
		return nil, err
	}
	return &CryptoService{key: key}, nil
}

func NewCryptoServiceWithKey(key []byte) *CryptoService {
	return &CryptoService{key: key}
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES-256 key: %w", err)
	}

	return key, nil
}

func (c *CryptoService) GetKey() []byte {
	result := make([]byte, len(c.key))
	copy(result, c.key)
	return result
}

func (c *CryptoService) CreateCipher() (cipher.Stream, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("create AES ciper fail: %v", err)
	}

	iv := make([]byte, aes.BlockSize)
	copy(iv, c.key[:aes.BlockSize])

	return cipher.NewCTR(block, iv), nil
}

func (c *CryptoService) ProcessData(stream cipher.Stream, input, output []byte, length int) {
	stream.XORKeyStream(output[:length], input[:length])
}

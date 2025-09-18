package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base32"
	"encoding/base64"
	"log"
	"time"

	"github.com/mint8846/Traversal-Learning/odc/internal/config"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// OTPService uses TOTP based on RFC 6238, but security can be enhanced by using a custom TOTP algorithm
// that truncates timestamps by specific time intervals and generates HMAC using seed values and timestamps.
type OTPService struct {
	cfg *config.Config
}

func NewOTPService(cfg *config.Config) *OTPService {
	return &OTPService{cfg: cfg}
}

func (c *OTPService) EncryptKey(plainKey []byte) (string, error) {
	otp, err := c.generateOTP()
	if err != nil {
		log.Printf("GenerateKey: generateOTP fail %v", err)
		return "", err
	}
	log.Printf("GenerateKey: otp: %s", otp)

	// Make AES key from OTP
	key := c.makeAESKey(otp)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("GenerateKey: NewCipher fail %v", err)
		return "", err
	}

	// Add padding
	data := addPadding(plainKey)

	iv := make([]byte, aes.BlockSize)
	copy(iv, key[:aes.BlockSize])

	// Encrypt
	encryptKey := make([]byte, len(data))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(encryptKey, data)

	return base64.StdEncoding.EncodeToString(encryptKey), nil
}

func (c *OTPService) DecryptKey(cipherText string) ([]byte, error) {
	otp, err := c.generateOTP()
	if err != nil {
		log.Printf("DecryptKey: generateOTP fail %v", err)
		return nil, err
	}
	log.Printf("DecryptKey: otp: %s", otp)

	// Decode base64
	byteData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	// Make AES key from OTP
	key := c.makeAESKey(otp)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("DecryptKey: NewCipher fail %v", err)
		return nil, err
	}

	// Extract IV and ciphertext
	iv := make([]byte, aes.BlockSize)
	copy(iv, key[:aes.BlockSize])

	// Decrypt
	decryptKey := make([]byte, len(byteData))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(decryptKey, byteData)

	// Remove padding
	result := removePadding(decryptKey)
	return result, nil
}

func (c *OTPService) generateOTP() (string, error) {
	secretKey := base32.StdEncoding.EncodeToString([]byte(c.cfg.OTP.Seed))
	if len(secretKey) > 32 {
		secretKey = secretKey[:32]
	}

	otp, err := totp.GenerateCodeCustom(secretKey, time.Now(), totp.ValidateOpts{
		Period:    c.cfg.OTP.Period,
		Skew:      1,
		Digits:    otp.DigitsEight,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
		return "", err
	}

	return otp, nil
}

func (c *OTPService) makeAESKey(otp string) []byte {
	key := make([]byte, 32)
	copy(key, otp)
	return key
}

// addPadding adds PKCS7 padding
func addPadding(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	paddingText := make([]byte, padding)
	for i := range paddingText {
		paddingText[i] = byte(padding)
	}
	return append(data, paddingText...)
}

// removePadding removes PKCS7 padding
func removePadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

package service

import (
	"encoding/base32"
	"log"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// OTPService uses TOTP based on RFC 6238, but security can be enhanced by using a custom TOTP algorithm
// that truncates timestamps by specific time intervals and generates HMAC using seed values and timestamps.
type OTPService struct{}

func (o *OTPService) Generate(seed string, period uint, time time.Time) ([]byte, error) {
	secretKey := base32.StdEncoding.EncodeToString([]byte(seed))
	if len(secretKey) > 32 {
		secretKey = secretKey[:32]
	}

	value, err := totp.GenerateCodeCustom(secretKey, time, totp.ValidateOpts{
		Period:    period,
		Skew:      1,
		Digits:    otp.DigitsEight,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
		return nil, err
	}
	log.Printf("GetKey: otp(%s)", value)

	// Make AES key from OTP
	return o.makeAESKey(value), nil
}

func (o *OTPService) makeAESKey(otp string) []byte {
	key := make([]byte, 32)
	copy(key, otp)
	return key
}

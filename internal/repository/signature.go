package repository

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"time"
)

type SignatureInput = struct {
	Timestamp        int
	Signature        string
	SigningSecret    string
	Body             string
	SignatureVersion string
}

type SignatureOutput struct {
}

type ISignatureRepository interface {
	Verify(input SignatureInput) error
}

type SignatureRepository struct {
}

func NewSignatureRepository() *SignatureRepository {
	return &SignatureRepository{}
}

func sign(base string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(base))
	return hex.EncodeToString(mac.Sum(nil))
}

func (*SignatureRepository) Verify(input SignatureInput) error {
	if math.Abs(float64(time.Now().Unix()-int64(input.Timestamp))) > 5*60 {
		return errors.New("input timestamp is newer or older than 5 minutes")
	}
	base := input.SignatureVersion + ":" + strconv.Itoa(input.Timestamp) + ":" + input.Body
	signature := input.SignatureVersion + "=" + sign(base, input.SigningSecret)
	if input.Signature != signature {
		return errors.New("invalid signature detected")
	}
	return nil
}

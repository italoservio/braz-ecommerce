package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/italoservio/braz_ecommerce/packages/exception"
)

type EncryptedText struct {
	EncryptedText string
	Salt          string
}

func Encrypt(secret string, text string) (*EncryptedText, error) {
	if secret == "" || text == "" {
		return nil, errors.New(exception.CodeValidationFailed)
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, errors.New(exception.CodeInternal)
	}

	salt := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, errors.New(exception.CodeInternal)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New(exception.CodeInternal)
	}

	cipherText := aesgcm.Seal(nil, salt, []byte(text), nil)

	return &EncryptedText{
		EncryptedText: hex.EncodeToString(cipherText),
		Salt:          hex.EncodeToString(salt),
	}, nil
}

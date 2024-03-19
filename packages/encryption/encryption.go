package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
)

//go:generate mockgen --source=encryption.go --destination=./mocks/encryption_interface_mock.go --package=encryption_mocks
type EncryptionInterface interface {
	Encrypt(
		ctx context.Context,
		secret string,
		text string,
	) (*EncryptedText, error)
}

type EncryptionImpl struct {
	logger logger.LoggerInterface
}

func NewEncryptionImpl(lg logger.LoggerInterface) *EncryptionImpl {
	return &EncryptionImpl{
		logger: lg,
	}
}

type EncryptedText struct {
	EncryptedText string
	Salt          string
}

func (e *EncryptionImpl) Encrypt(ctx context.Context, secret string, text string) (*EncryptedText, error) {
	if secret == "" || text == "" {
		e.logger.WithCtx(ctx).Error("secret or text is empty")
		return nil, errors.New(exception.CodeValidationFailed)
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		e.logger.WithCtx(ctx).Error(err.Error())
		return nil, errors.New(exception.CodeInternal)
	}

	salt := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		e.logger.WithCtx(ctx).Error(err.Error())
		return nil, errors.New(exception.CodeInternal)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		e.logger.WithCtx(ctx).Error(err.Error())
		return nil, errors.New(exception.CodeInternal)
	}

	cipherText := aesgcm.Seal(nil, salt, []byte(text), nil)

	return &EncryptedText{
		EncryptedText: hex.EncodeToString(cipherText),
		Salt:          hex.EncodeToString(salt),
	}, nil
}

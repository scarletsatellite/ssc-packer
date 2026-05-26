package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

const (
	KeySize   = 32
	IterCount = 600_000
)

func DeriveKey(password []byte, salt []byte) []byte {
	return pbkdf2.Key(password, salt, IterCount, KeySize, sha256.New)
}

func EncryptGCM(plainData, password, salt, nonce []byte) ([]byte, error) {
	key := DeriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, plainData, nil), nil
}

func DecryptGCM(cipherData, password, salt, nonce []byte) ([]byte, error) {
	key := DeriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, cipherData, nil)
}

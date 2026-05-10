package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// Encrypt encrypts data using AES-256-GCM.
func Encrypt(data []byte, key []byte) ([]byte, error) {
	// Ensure key is exactly 32 bytes
	k := normalizeKey(key)

	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt decrypts data using AES-256-GCM.
func Decrypt(data []byte, key []byte) ([]byte, error) {
	k := normalizeKey(key)

	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed — data may be corrupted or wrong key")
	}
	return plain, nil
}

// DeriveKey turns a password into a 32-byte AES key via SHA-256.
func DeriveKey(password []byte) []byte {
	hash := sha256.Sum256(password)
	return hash[:]
}

// normalizeKey pads or truncates key to exactly 32 bytes.
func normalizeKey(key []byte) []byte {
	k := make([]byte, 32)
	copy(k, key)
	return k
}

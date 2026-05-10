package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

// Encrypt encrypts data using AES-256-GCM.
func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
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
	block, err := aes.NewCipher(key)
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

// DeriveHardwareKey turns a passphrase into a 32-byte AES key via Argon2id.
// It includes a device-specific secret to ensure the key is hardware-locked.
func DeriveHardwareKey(passphrase []byte, salt []byte, deviceSecret []byte) []byte {
	// Combine salt and deviceSecret for true hardware-locking
	combinedSalt := append(salt, deviceSecret...)
	
	// Argon2id parameters:
	// time=3, memory=64MB, threads=4, keyLen=32
	return argon2.IDKey(passphrase, combinedSalt, 3, 64*1024, 4, 32)
}

// DeriveSimpleKey turns a passphrase into a 32-byte AES key via Argon2id.
// Used for sharing bundles where hardware-locking is not desired.
func DeriveSimpleKey(passphrase []byte, salt []byte) []byte {
	return argon2.IDKey(passphrase, salt, 3, 64*1024, 4, 32)
}

// GenerateSalt creates a random salt of the specified size.
func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

package keychain

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/99designs/keyring"
)

const (
	serviceName = "envy-cli"
	keyName     = "device-secret"
)

func GetOrCreateDeviceSecret() ([]byte, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
	})
	if err != nil {
		return nil, err
	}

	item, err := ring.Get(keyName)
	if err == nil {
		return hex.DecodeString(string(item.Data))
	}

	// Create new secret
	secret := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		return nil, err
	}

	err = ring.Set(keyring.Item{
		Key:  keyName,
		Data: []byte(hex.EncodeToString(secret)),
	})
	if err != nil {
		return nil, err
	}

	return secret, nil
}

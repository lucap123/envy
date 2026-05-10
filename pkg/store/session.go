package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/lucap/envy/pkg/crypto"
)

type sessionFile struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	KeyEnc    string    `json:"key_enc"` // key encrypted with session token
}

const sessionDuration = 8 * time.Hour

func (s *Store) sessionPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envy", "session")
}

// tokenKey derives a 32-byte AES key from a session token.
func tokenKey(token []byte) []byte {
	h := sha256.Sum256(token)
	return h[:]
}

// SaveSession caches the encryption key for sessionDuration.
func (s *Store) SaveSession(key []byte) error {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(key, tokenKey(token))
	if err != nil {
		return err
	}

	sf := sessionFile{
		Token:     hex.EncodeToString(token),
		ExpiresAt: time.Now().Add(sessionDuration),
		KeyEnc:    hex.EncodeToString(encrypted),
	}

	data, err := json.Marshal(sf)
	if err != nil {
		return err
	}

	return os.WriteFile(s.sessionPath(), data, 0600)
}

// LoadSession returns the cached key if a valid session exists.
func (s *Store) LoadSession() ([]byte, bool) {
	data, err := os.ReadFile(s.sessionPath())
	if err != nil {
		return nil, false
	}

	var sf sessionFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return nil, false
	}

	if time.Now().After(sf.ExpiresAt) {
		os.Remove(s.sessionPath())
		return nil, false
	}

	token, err := hex.DecodeString(sf.Token)
	if err != nil {
		return nil, false
	}

	encrypted, err := hex.DecodeString(sf.KeyEnc)
	if err != nil {
		return nil, false
	}

	key, err := crypto.Decrypt(encrypted, tokenKey(token))
	if err != nil {
		return nil, false
	}

	return key, true
}

// ClearSession removes the cached session.
func (s *Store) ClearSession() {
	os.Remove(s.sessionPath())
}

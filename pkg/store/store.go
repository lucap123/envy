package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucap/envy/pkg/crypto"
)

// ProjectVars holds all profiles for a project. Exported so cmd/share.go can use it.
type ProjectVars struct {
	Profiles map[string]map[string]string `json:"profiles"`
	Active   string                       `json:"active"`
}

func (pv *ProjectVars) GetActiveVars() map[string]string {
	if vars, ok := pv.Profiles[pv.Active]; ok {
		return vars
	}
	return make(map[string]string)
}

type Store struct {
	baseDir string
	key     []byte
}

func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(home, ".envy", "sessions")
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, err
	}

	key, err := getOrCreateKey(filepath.Join(home, ".envy"))
	if err != nil {
		return nil, err
	}

	return &Store{baseDir: baseDir, key: key}, nil
}

func getOrCreateKey(dir string) ([]byte, error) {
	keyPath := filepath.Join(dir, "key")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		// Generate a random 32-byte key
		k := make([]byte, 32)
		f, err := os.Open("/dev/urandom")
		if err != nil {
			// Fallback: derive from hostname + username
			hostname, _ := os.Hostname()
			h := sha256.Sum256([]byte(hostname + os.Getenv("USER")))
			k = h[:]
		} else {
			f.Read(k)
			f.Close()
		}
		if err := os.WriteFile(keyPath, k, 0600); err != nil {
			return nil, fmt.Errorf("could not write key file: %w", err)
		}
		return k, nil
	}
	return os.ReadFile(keyPath)
}

func (s *Store) projectID() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(cwd))
	return hex.EncodeToString(hash[:16]), nil // 32 hex chars is enough
}

func (s *Store) Load() (*ProjectVars, error) {
	id, err := s.projectID()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(s.baseDir, id+".enc")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &ProjectVars{
			Profiles: map[string]map[string]string{"default": {}},
			Active:   "default",
		}, nil
	}

	encrypted, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	decrypted, err := crypto.Decrypt(encrypted, s.key)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt store: %w", err)
	}

	var pv ProjectVars
	if err := json.Unmarshal(decrypted, &pv); err != nil {
		return nil, fmt.Errorf("corrupted store: %w", err)
	}
	if pv.Profiles == nil {
		pv.Profiles = map[string]map[string]string{"default": {}}
	}
	if pv.Active == "" {
		pv.Active = "default"
	}

	return &pv, nil
}

func (s *Store) Save(pv *ProjectVars) error {
	id, err := s.projectID()
	if err != nil {
		return err
	}

	data, err := json.Marshal(pv)
	if err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(data, s.key)
	if err != nil {
		return err
	}

	path := filepath.Join(s.baseDir, id+".enc")
	return os.WriteFile(path, encrypted, 0600)
}

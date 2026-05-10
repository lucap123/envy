package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucap123/envy/pkg/crypto"
)

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
	salt    []byte
}

type Config struct {
	Salt   string `json:"salt"`
	Canary string `json:"canary"`
}

func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	envyDir := filepath.Join(home, ".envy")
	baseDir := filepath.Join(envyDir, "sessions")
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, err
	}

	configPath := filepath.Join(envyDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Store{baseDir: baseDir}, nil
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	salt, err := hex.DecodeString(config.Salt)
	if err != nil {
		return nil, err
	}

	return &Store{baseDir: baseDir, salt: salt}, nil
}

func (s *Store) IsInitialized() bool {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".envy", "config.json")
	_, err := os.Stat(configPath)
	return err == nil
}

func (s *Store) Initialize(passphrase string, deviceSecret []byte) error {
	salt, err := crypto.GenerateSalt(16)
	if err != nil {
		return err
	}

	key := crypto.DeriveHardwareKey([]byte(passphrase), salt, deviceSecret)
	canary, err := crypto.Encrypt([]byte("envy-canary"), key)
	if err != nil {
		return err
	}

	config := Config{
		Salt:   hex.EncodeToString(salt),
		Canary: hex.EncodeToString(canary),
	}

	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	home, _ := os.UserHomeDir()
	envyDir := filepath.Join(home, ".envy")
	configPath := filepath.Join(envyDir, "config.json")
	
	if err := os.WriteFile(configPath, configData, 0600); err != nil {
		return err
	}

	s.salt = salt
	return nil
}

func (s *Store) VerifyPassphrase(passphrase string, deviceSecret []byte) ([]byte, error) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".envy", "config.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	salt, _ := hex.DecodeString(config.Salt)
	canary, _ := hex.DecodeString(config.Canary)

	key := crypto.DeriveHardwareKey([]byte(passphrase), salt, deviceSecret)
	decrypted, err := crypto.Decrypt(canary, key)
	if err != nil || string(decrypted) != "envy-canary" {
		return nil, errors.New("invalid passphrase")
	}

	return key, nil
}

func (s *Store) projectID() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(cwd))
	return hex.EncodeToString(hash[:16]), nil
}

func (s *Store) Load(key []byte) (*ProjectVars, error) {
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

	decrypted, err := crypto.Decrypt(encrypted, key)
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

func (s *Store) Save(pv *ProjectVars, key []byte) error {
	id, err := s.projectID()
	if err != nil {
		return err
	}

	data, err := json.Marshal(pv)
	if err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(data, key)
	if err != nil {
		return err
	}

	path := filepath.Join(s.baseDir, id+".enc")
	return os.WriteFile(path, encrypted, 0600)
}

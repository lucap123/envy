package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/lucap/envy/pkg/crypto"
	"github.com/lucap/envy/pkg/store"
	"golang.org/x/term"
)

func Share(s *store.Store, args []string) error {
	outFile := "team.env.enc"
	for i, a := range args {
		if a == "--output" && i+1 < len(args) {
			outFile = args[i+1]
		}
	}

	pv, err := s.Load()
	if err != nil {
		return err
	}

	if len(pv.GetActiveVars()) == 0 {
		return fmt.Errorf("no variables in active profile to share")
	}

	fmt.Print("Enter encryption password: ")
	pw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("could not read password: %w", err)
	}
	if len(pw) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	fmt.Print("Confirm password: ")
	pw2, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("could not read password: %w", err)
	}
	if string(pw) != string(pw2) {
		return fmt.Errorf("passwords do not match")
	}

	data, err := json.Marshal(pv)
	if err != nil {
		return err
	}

	key := crypto.DeriveKey(pw)
	encrypted, err := crypto.Encrypt(data, key)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outFile, encrypted, 0600); err != nil {
		return fmt.Errorf("could not write bundle: %w", err)
	}

	fmt.Printf("Exported encrypted bundle to %s\n", outFile)
	fmt.Println("Share this file + the password via separate channels.")
	fmt.Printf("Recipient runs: envy import %s\n", outFile)
	return nil
}

func Import(s *store.Store, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envy import <file>")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	fmt.Print("Enter decryption password: ")
	pw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("could not read password: %w", err)
	}

	key := crypto.DeriveKey(pw)
	decrypted, err := crypto.Decrypt(data, key)
	if err != nil {
		return fmt.Errorf("decryption failed — wrong password?")
	}

	var imported store.ProjectVars
	if err := json.Unmarshal(decrypted, &imported); err != nil {
		return fmt.Errorf("invalid bundle format")
	}

	// Merge into current store — don't overwrite, add missing profiles
	current, err := s.Load()
	if err != nil {
		return err
	}

	imported_count := 0
	for profile, vars := range imported.Profiles {
		if _, exists := current.Profiles[profile]; !exists {
			current.Profiles[profile] = vars
			imported_count += len(vars)
		}
	}

	if err := s.Save(current); err != nil {
		return err
	}

	fmt.Printf("Imported %d variables across %d profile(s)\n", imported_count, len(imported.Profiles))
	return nil
}

package cmd

import (
	"fmt"
	"syscall"

	"github.com/awnumar/memguard"
	"github.com/lucap123/envy/pkg/keychain"
	"github.com/lucap123/envy/pkg/store"
	"golang.org/x/term"
)

func Execute(args []string) error {
	memguard.CatchInterrupt()
	defer memguard.Purge()

	s, err := store.NewStore()
	if err != nil {
		return fmt.Errorf("failed to init store: %w", err)
	}

	ds, err := keychain.GetOrCreateDeviceSecret()
	if err != nil {
		return fmt.Errorf("failed to get device secret: %w", err)
	}

	cmdName := args[0]

	// Commands that never need a key
	switch cmdName {
	case "help", "--help", "-h":
		PrintUsage()
		return nil
	case "hook":
		return Hook(args[1:])
	case "logout":
		s.ClearSession()
		fmt.Println("Session cleared.")
		return nil
	}

	// init: may need to set up master passphrase first
	if cmdName == "init" {
		if !s.IsInitialized() {
			fmt.Println("Welcome to envy! Let's set your master passphrase.")
			fmt.Println("This protects all your stored secrets.")
			fmt.Println()
			pass, err := promptPassphrase("Set master passphrase: ", true)
			if err != nil {
				return err
			}
			defer pass.Destroy()
			if err := s.Initialize(string(pass.Bytes()), ds); err != nil {
				return fmt.Errorf("failed to initialize: %w", err)
			}
			fmt.Println("✓ Master passphrase set. You're ready to use envy.")
			fmt.Println()
		}
		return Init(args[1:])
	}

	// All other commands need the encryption key
	key, err := getKey(s, ds)
	if err != nil {
		return err
	}
	defer key.Destroy()

	switch cmdName {
	case "set":
		return Set(s, key.Bytes(), args[1:])
	case "get":
		return Get(s, key.Bytes(), args[1:])
	case "list":
		return List(s, key.Bytes(), args[1:])
	case "run":
		return Run(s, key.Bytes(), args[1:])
	case "export":
		return Export(s, key.Bytes(), args[1:])
	case "profile":
		return Profile(s, key.Bytes(), args[1:])
	case "share":
		return Share(s, key.Bytes(), args[1:])
	case "import":
		return Import(s, key.Bytes(), args[1:])
	default:
		return fmt.Errorf("unknown command: %s\n\nRun 'envy help' for usage", cmdName)
	}
}

// getKey returns the encryption key from session cache or passphrase prompt.
func getKey(s *store.Store, ds []byte) (*memguard.LockedBuffer, error) {
	if !s.IsInitialized() {
		return nil, fmt.Errorf("envy is not initialized. Run 'envy init' first")
	}

	// Try session cache first (valid for 8 hours)
	if cached, ok := s.LoadSession(); ok {
		return memguard.NewBufferFromBytes(cached), nil
	}

	// Prompt for passphrase
	pass, err := promptPassphrase("Master passphrase: ", false)
	if err != nil {
		return nil, err
	}
	defer pass.Destroy()

	keyBytes, err := s.VerifyPassphrase(string(pass.Bytes()), ds)
	if err != nil {
		return nil, fmt.Errorf("wrong passphrase")
	}

	// Cache session so next commands skip the prompt (8 hours)
	_ = s.SaveSession(keyBytes)

	return memguard.NewBufferFromBytes(keyBytes), nil
}

func promptPassphrase(prompt string, confirm bool) (*memguard.LockedBuffer, error) {
	fmt.Print(prompt)
	raw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return nil, err
	}

	pass := memguard.NewBufferFromBytes(raw)

	if pass.Size() < 8 {
		pass.Destroy()
		return nil, fmt.Errorf("passphrase must be at least 8 characters")
	}

	if confirm {
		fmt.Print("Confirm passphrase: ")
		raw2, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			pass.Destroy()
			return nil, err
		}
		c := memguard.NewBufferFromBytes(raw2)
		defer c.Destroy()

		if !pass.EqualTo(c.Bytes()) {
			pass.Destroy()
			return nil, fmt.Errorf("passphrases do not match")
		}
	}

	return pass, nil
}

func PrintUsage() {
	fmt.Println(`envy — zero-config environment variable manager
  Hardware-locked · AES-256-GCM · Argon2id · No cloud

Usage:
  envy init                        setup envy & detect project type
  envy set <key> <value>           store a variable (encrypted)
  envy get <key>                   retrieve a variable
  envy list [--reveal]             list variables (masked by default)
  envy run <cmd> [args...]         run command with vars injected
  envy export [--file <path>]      write a .env file
  envy hook install                install git pre-commit secret scanner
  envy hook uninstall              remove git pre-commit hook
  envy profile list                list all profiles
  envy profile add <name>          create a new profile
  envy profile use <name>          switch active profile
  envy profile delete <name>       delete a profile
  envy share --output <file>       export encrypted bundle for team sharing
  envy import <file>               import an encrypted bundle
  envy logout                      clear session (force passphrase re-entry)

Session: passphrase cached for 8 hours after first unlock
Storage: ~/.envy/  (AES-256-GCM · Argon2id · hardware-locked)`)
}

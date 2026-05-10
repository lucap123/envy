package cmd

import (
	"fmt"

	"github.com/lucap/envy/pkg/store"
)

func Execute(args []string) error {
	s, err := store.NewStore()
	if err != nil {
		return fmt.Errorf("failed to init store: %w", err)
	}

	switch args[0] {
	case "init":
		return Init(args[1:])
	case "set":
		return Set(s, args[1:])
	case "get":
		return Get(s, args[1:])
	case "list":
		return List(s, args[1:])
	case "run":
		return Run(s, args[1:])
	case "export":
		return Export(s, args[1:])
	case "hook":
		return Hook(args[1:])
	case "profile":
		return Profile(s, args[1:])
	case "share":
		return Share(s, args[1:])
	case "import":
		return Import(s, args[1:])
	case "help", "--help", "-h":
		PrintUsage()
		return nil
	default:
		return fmt.Errorf("unknown command: %s\n\nRun 'envy help' for usage", args[0])
	}
}

func PrintUsage() {
	fmt.Println(`envy — zero-config environment variable manager

Usage:
  envy init                        detect project type, create .env.example
  envy set <key> <value>           store a variable in the active profile
  envy get <key>                   retrieve a variable (prints plain value)
  envy list [--reveal]             list all variables (masked by default)
  envy run <cmd> [args...]         run a command with all vars injected
  envy export [--file <path>]      write a .env file from active profile
  envy hook install                install git pre-commit secret scanner
  envy profile list                list all profiles
  envy profile add <name>          create a new profile
  envy profile use <name>          switch active profile
  envy share --output <file>       export encrypted bundle for team sharing
  envy import <file>               import an encrypted bundle

Storage: ~/.envy/  (AES-256-GCM encrypted, chmod 600, no cloud)`)
}

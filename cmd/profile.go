package cmd

import (
	"fmt"

	"github.com/lucap/envy/pkg/store"
)

func Profile(s *store.Store, key []byte, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envy profile [list|add|use|delete] <name>")
	}

	pv, err := s.Load(key)
	if err != nil {
		return err
	}

	switch args[0] {
	case "list":
		fmt.Printf("Profiles (%d):\n", len(pv.Profiles))
		for name, vars := range pv.Profiles {
			marker := "  "
			if name == pv.Active {
				marker = "→ "
			}
			fmt.Printf("%s%-20s (%d vars)\n", marker, name, len(vars))
		}

	case "add":
		if len(args) < 2 {
			return fmt.Errorf("usage: envy profile add <name>")
		}
		name := args[1]
		if _, ok := pv.Profiles[name]; ok {
			return fmt.Errorf("profile '%s' already exists", name)
		}
		pv.Profiles[name] = make(map[string]string)
		if err := s.Save(pv, key); err != nil {
			return err
		}
		fmt.Printf("Created profile '%s'\n", name)
		fmt.Printf("Tip: switch to it with 'envy profile use %s'\n", name)

	case "use":
		if len(args) < 2 {
			return fmt.Errorf("usage: envy profile use <name>")
		}
		name := args[1]
		if _, ok := pv.Profiles[name]; !ok {
			return fmt.Errorf("profile '%s' does not exist — create it with 'envy profile add %s'", name, name)
		}
		pv.Active = name
		if err := s.Save(pv, key); err != nil {
			return err
		}
		fmt.Printf("Switched to profile '%s'\n", name)

	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("usage: envy profile delete <name>")
		}
		name := args[1]
		if name == "default" {
			return fmt.Errorf("cannot delete the default profile")
		}
		if name == pv.Active {
			return fmt.Errorf("cannot delete the active profile — switch away first")
		}
		if _, ok := pv.Profiles[name]; !ok {
			return fmt.Errorf("profile '%s' does not exist", name)
		}
		delete(pv.Profiles, name)
		if err := s.Save(pv, key); err != nil {
			return err
		}
		fmt.Printf("Deleted profile '%s'\n", name)

	default:
		return fmt.Errorf("unknown subcommand: %s", args[0])
	}

	return nil
}

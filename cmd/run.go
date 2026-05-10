package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lucap123/envy/pkg/store"
)

func Run(s *store.Store, key []byte, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envy run <command> [args...]")
	}

	verbose := false
	cmdArgs := args
	if args[0] == "--verbose" || args[0] == "-v" {
		verbose = true
		cmdArgs = args[1:]
		if len(cmdArgs) < 1 {
			return fmt.Errorf("usage: envy run [--verbose] <command> [args...]")
		}
	}

	pv, err := s.Load(key)
	if err != nil {
		return err
	}

	vars := pv.GetActiveVars()

	if verbose {
		fmt.Printf("[envy] profile: %s (%d vars injected)\n", pv.Active, len(vars))
		fmt.Printf("[envy] running: %s\n", joinCmd(cmdArgs))
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()
	for k, v := range vars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}

func joinCmd(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += " "
		}
		if containsSpace(a) {
			result += fmt.Sprintf("%q", a)
		} else {
			result += a
		}
	}
	return result
}

func containsSpace(s string) bool {
	for _, c := range s {
		if c == ' ' || c == '\t' {
			return true
		}
	}
	return false
}

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lucap123/envy/pkg/secrets"
)

const hookScript = `#!/bin/sh
# envy pre-commit hook — blocks commits containing secrets
envy hook run
`

func Hook(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envy hook [install|uninstall|run]")
	}

	switch args[0] {
	case "install":
		return hookInstall()
	case "uninstall":
		return hookUninstall()
	case "run":
		return hookRun()
	default:
		return fmt.Errorf("unknown hook subcommand: %s", args[0])
	}
}

func hookInstall() error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("no .git directory found — are you in a git repository?")
	}

	hookPath := ".git/hooks/pre-commit"

	// Don't overwrite an existing hook without warning
	if _, err := os.Stat(hookPath); err == nil {
		content, _ := os.ReadFile(hookPath)
		if !strings.Contains(string(content), "envy") {
			return fmt.Errorf("a pre-commit hook already exists at %s\nManually add 'envy hook run' to it, or remove it first", hookPath)
		}
		fmt.Println("envy hook already installed.")
		return nil
	}

	if err := os.WriteFile(hookPath, []byte(hookScript), 0755); err != nil {
		return fmt.Errorf("could not write hook: %w", err)
	}

	fmt.Println("✓ Installed git pre-commit hook")
	fmt.Println("  Every commit will now be scanned for secrets.")
	fmt.Println("  To remove: envy hook uninstall")
	return nil
}

func hookUninstall() error {
	hookPath := ".git/hooks/pre-commit"
	content, err := os.ReadFile(hookPath)
	if err != nil {
		return fmt.Errorf("no pre-commit hook found")
	}

	if !strings.Contains(string(content), "envy") {
		return fmt.Errorf("the existing hook was not installed by envy — remove it manually")
	}

	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("could not remove hook: %w", err)
	}

	fmt.Println("Removed envy pre-commit hook.")
	return nil
}

func hookRun() error {
	out, err := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM").Output()
	if err != nil {
		return fmt.Errorf("could not get staged files: %w", err)
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil // nothing staged
	}

	files := strings.Split(raw, "\n")
	fmt.Printf("envy: scanning %d staged file(s)...\n", len(files))

	var allFindings []secrets.Finding
	for _, f := range files {
		findings, err := secrets.ScanFile(f)
		if err != nil {
			continue // file may have been deleted
		}
		allFindings = append(allFindings, findings...)
	}

	if len(allFindings) == 0 {
		fmt.Println("envy: no secrets found ✓")
		return nil
	}

	fmt.Printf("\nenvy: BLOCKED — %d potential secret(s) found:\n\n", len(allFindings))
	for _, f := range allFindings {
		fmt.Println(f.String())
	}
	fmt.Println("\nFix the issues above, then commit again.")
	fmt.Println("To skip scanning (not recommended): git commit --no-verify")
	os.Exit(1)
	return nil
}

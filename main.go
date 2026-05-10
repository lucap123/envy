package main

import (
	"fmt"
	"os"

	"github.com/lucap123/envy/cmd"
)

// Version is set at build time via -ldflags
var Version = "dev"

func main() {
	if len(os.Args) < 2 {
		cmd.PrintUsage()
		os.Exit(0)
	}

	if os.Args[1] == "version" || os.Args[1] == "--version" {
		fmt.Printf("envy %s\n", Version)
		os.Exit(0)
	}

	if err := cmd.Execute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

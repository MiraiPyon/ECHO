package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ciferia/echo/internal/cli"
	"github.com/ciferia/echo/internal/tui"
	"golang.org/x/term"
)

func main() {
	binary := os.Args[0]
	args := os.Args

	for {
		if len(args) < 2 {
			if term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd())) {
				interactiveArgs, err := tui.Launch(binary)
				switch {
				case err == nil:
					args = interactiveArgs
				case errors.Is(err, tui.ErrCancelled):
					return
				case strings.Contains(err.Error(), "cursor addressable"):
					args = cli.PromptInteractiveArgs(binary)
				default:
					fmt.Fprintf(os.Stderr, "Could not launch the interface: %v\n", err)
					args = cli.PromptInteractiveArgs(binary)
				}
			} else {
				args = cli.PromptInteractiveArgs(binary)
			}
		}

		err := cli.Execute(args)
		if err == nil {
			return
		}
		if errors.Is(err, tui.ErrBackToMainMenu) {
			args = []string{binary}
			continue
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

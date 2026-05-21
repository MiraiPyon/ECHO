package cli

import (
	"fmt"

	"github.com/ciferia/echo/internal/commands/concat"
	"github.com/ciferia/echo/internal/commands/extract"
	"github.com/ciferia/echo/internal/commands/play"
	"github.com/ciferia/echo/internal/commands/trim"
	"github.com/ciferia/echo/internal/commands/volume"
)

// Execute parses the command line and dispatches to the right handler.
func Execute(args []string) error {
	if len(args) < 2 {
		PrintHelp()
		return nil
	}

	command := args[1]
	flags := args[2:]

	switch command {
	case "trim":
		return trim.Run(flags)
	case "play":
		return play.Run(flags)
	case "concat":
		return concat.Run(flags)
	case "extract":
		return extract.Run(flags)
	case "volume":
		return volume.Run(flags)
	case "help", "--help", "-h":
		PrintHelp()
		return nil
	default:
		return fmt.Errorf("invalid command %q. Run 'echo --help' to see the available commands", command)
	}
}

// PrintHelp prints the basic usage guide.
func PrintHelp() {
	fmt.Println("ECHO - Terminal audio tool")
	fmt.Println("Usage: echo <command> [input files] [flags]")
	fmt.Println("\nCommands:")
	fmt.Println("  play     - Preview an audio file")
	fmt.Println("  trim     - Trim an audio file")
	fmt.Println("  concat   - Concatenate multiple audio files")
	fmt.Println("  extract  - Extract audio from a video file")
	fmt.Println("  volume   - Adjust audio volume")
	fmt.Println("\nTip: running `echo` without a command opens the interactive launcher.")
}

package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// PromptInteractiveArgs opens a simple terminal prompt and returns argv-style arguments.
func PromptInteractiveArgs(binary string) []string {
	return PromptInteractiveArgsWithIO(binary, os.Stdin, os.Stdout)
}

// PromptInteractiveArgsWithIO is the testable version of the interactive prompt.
func PromptInteractiveArgsWithIO(binary string, stdin io.Reader, stdout io.Writer) []string {
	reader := bufio.NewReader(stdin)

	fmt.Fprintln(stdout, "ECHO - Terminal audio tool")
	fmt.Fprintln(stdout, "Choose a command:")
	fmt.Fprintln(stdout, "  1) play")
	fmt.Fprintln(stdout, "  2) trim")
	fmt.Fprintln(stdout, "  3) concat")
	fmt.Fprintln(stdout, "  4) extract")
	fmt.Fprintln(stdout, "  5) volume")
	fmt.Fprint(stdout, "\nEnter a command number or name: ")

	choice, err := reader.ReadString('\n')
	if err != nil && strings.TrimSpace(choice) == "" {
		return []string{binary}
	}

	command := NormalizeCommand(choice)
	if command == "" {
		return []string{binary}
	}

	if command == "help" {
		return []string{binary, command}
	}

	if command == "play" {
		fmt.Fprint(stdout, "Enter the audio file to play: ")
		file, err := reader.ReadString('\n')
		if err != nil && strings.TrimSpace(file) == "" {
			return []string{binary, command}
		}

		file = strings.TrimSpace(file)
		if file == "" {
			return []string{binary, command}
		}

		return []string{binary, command, file}
	}

	fmt.Fprint(stdout, "Enter command arguments (example: input.mp3 --start 00:30 --end 01:30 --out out.mp3): ")
	rest, err := reader.ReadString('\n')
	if err != nil && strings.TrimSpace(rest) == "" {
		return []string{binary, command}
	}

	fields := strings.Fields(strings.TrimSpace(rest))
	return append([]string{binary, command}, fields...)
}

// NormalizeCommand maps menu input to a canonical command name.
func NormalizeCommand(input string) string {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "1", "play":
		return "play"
	case "2", "trim":
		return "trim"
	case "3", "concat":
		return "concat"
	case "4", "extract":
		return "extract"
	case "5", "volume":
		return "volume"
	case "help", "--help", "-h":
		return "help"
	default:
		return ""
	}
}

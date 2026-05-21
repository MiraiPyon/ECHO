package play

import (
	"fmt"
	"os"
	"strings"

	"github.com/ciferia/echo/internal/utils"
	"golang.org/x/term"
)

func Run(flags []string) error {
	if len(flags) == 0 {
		return fmt.Errorf("missing audio file. Example: echo play input.mp3")
	}
	if len(flags) > 1 {
		return fmt.Errorf("play accepts exactly one audio file. Example: echo play input.mp3")
	}

	filePath := strings.TrimSpace(flags[0])
	if filePath == "" {
		return fmt.Errorf("missing audio file. Example: echo play input.mp3")
	}

	if err := utils.ValidateFileExists(filePath); err != nil {
		return err
	}

	if term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd())) {
		return launchPlayerUI(filePath)
	}

	return playSimple(filePath)
}

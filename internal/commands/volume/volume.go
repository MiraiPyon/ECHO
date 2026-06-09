package volume

import (
	"fmt"

	"github.com/ciferia/echo/internal/audio"
	"github.com/ciferia/echo/internal/utils"
)

// Run executes the volume command.
// Usage: echo volume input.mp3 --level 2.0 --out louder.mp3
func Run(flags []string) error {
	parsed, err := utils.ParseArgs(flags)
	if err != nil {
		return err
	}

	if len(parsed.Positional) == 0 {
		return fmt.Errorf("missing input audio file. Usage: echo volume input.mp3 --level 2.0 --out louder.mp3")
	}
	if len(parsed.Positional) > 1 {
		return fmt.Errorf("volume command accepts exactly one input audio file. Usage: echo volume input.mp3 --level 2.0 --out louder.mp3")
	}

	inputFile := parsed.Positional[0]
	if err := utils.ValidateFileExists(inputFile); err != nil {
		return err
	}

	out := parsed.Flags["out"]
	if out == "" {
		out = parsed.Flags["o"]
	}
	if out == "" {
		return fmt.Errorf("missing output file flag (--out or -o)")
	}

	level := parsed.Flags["level"]
	if level == "" {
		level = parsed.Flags["l"]
	}
	if level == "" {
		return fmt.Errorf("missing volume level flag (--level or -l). Example: 2.0 or 0.5")
	}

	// Build ffmpeg arguments
	ffmpegArgs := []string{
		"-y",
		"-i", inputFile,
		"-filter:a", fmt.Sprintf("volume=%s", level),
		out,
	}

	fmt.Printf("Adjusting volume of %q to %s and saving to %q...\n", inputFile, level, out)
	if err := audio.RunFFmpeg(ffmpegArgs); err != nil {
		return err
	}

	fmt.Printf("Successfully adjusted volume and saved to %q\n", out)
	return nil
}

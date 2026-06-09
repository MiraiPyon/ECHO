package extract

import (
	"fmt"

	"github.com/ciferia/echo/internal/audio"
	"github.com/ciferia/echo/internal/utils"
)

// Run executes the extract command.
// Usage: echo extract video.mp4 --out audio.mp3
func Run(flags []string) error {
	parsed, err := utils.ParseArgs(flags)
	if err != nil {
		return err
	}

	if len(parsed.Positional) == 0 {
		return fmt.Errorf("missing input video file. Usage: echo extract video.mp4 --out audio.mp3")
	}
	if len(parsed.Positional) > 1 {
		return fmt.Errorf("extract command accepts exactly one input video file. Usage: echo extract video.mp4 --out audio.mp3")
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

	// Build ffmpeg arguments
	ffmpegArgs := []string{
		"-y",
		"-i", inputFile,
		"-vn",
		out,
	}

	fmt.Printf("Extracting audio from %q to %q...\n", inputFile, out)
	if err := audio.RunFFmpeg(ffmpegArgs); err != nil {
		return err
	}

	fmt.Printf("Successfully extracted audio and saved to %q\n", out)
	return nil
}

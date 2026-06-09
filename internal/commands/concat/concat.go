package concat

import (
	"fmt"
	"os"
	"strings"

	"github.com/ciferia/echo/internal/audio"
	"github.com/ciferia/echo/internal/utils"
)

// Run executes the concat command.
// Usage: echo concat file1.mp3 file2.mp3 --out out.mp3
func Run(flags []string) error {
	parsed, err := utils.ParseArgs(flags)
	if err != nil {
		return err
	}

	if len(parsed.Positional) < 2 {
		return fmt.Errorf("concat command requires at least two input files. Usage: echo concat file1.mp3 file2.mp3 --out out.mp3")
	}

	out := parsed.Flags["out"]
	if out == "" {
		out = parsed.Flags["o"]
	}
	if out == "" {
		return fmt.Errorf("missing output file flag (--out or -o)")
	}

	// Validate input files exist
	for _, file := range parsed.Positional {
		if err := utils.ValidateFileExists(file); err != nil {
			return err
		}
	}

	// Create temporary file for ffmpeg's concat demuxer
	tmpFile, err := os.CreateTemp("", "echo-concat-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	for _, file := range parsed.Positional {
		// Prepare absolute or relative paths with forward slashes for FFmpeg
		escapedPath := strings.ReplaceAll(file, "\\", "/")
		escapedPath = strings.ReplaceAll(escapedPath, "'", "'\\''")
		if _, err := fmt.Fprintf(tmpFile, "file '%s'\n", escapedPath); err != nil {
			return fmt.Errorf("failed to write to temporary file: %w", err)
		}
	}

	// Close the file so FFmpeg can read it on Windows
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Build ffmpeg arguments
	ffmpegArgs := []string{
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", tmpFile.Name(),
		"-c", "copy",
		out,
	}

	fmt.Printf("Concatenating %d files to %q...\n", len(parsed.Positional), out)
	if err := audio.RunFFmpeg(ffmpegArgs); err != nil {
		return err
	}

	fmt.Printf("Successfully concatenated files and saved to %q\n", out)
	return nil
}

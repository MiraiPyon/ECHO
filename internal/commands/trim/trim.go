package trim

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ciferia/echo/internal/audio"
	"github.com/ciferia/echo/internal/utils"
)

// Run executes the trim command.
// Usage: echo trim input.mp3 --start 00:30 --end 01:30 --out out.mp3
func Run(flags []string) error {
	parsed, err := utils.ParseArgs(flags)
	if err != nil {
		return err
	}

	if len(parsed.Positional) == 0 {
		return fmt.Errorf("missing input audio file. Usage: echo trim input.mp3 --start 00:30 --end 01:30 --out out.mp3")
	}
	if len(parsed.Positional) > 1 {
		return fmt.Errorf("trim command accepts exactly one input audio file. Usage: echo trim input.mp3 --start 00:30 --end 01:30 --out out.mp3")
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

	start := parsed.Flags["start"]
	if start == "" {
		start = parsed.Flags["s"]
	}

	end := parsed.Flags["end"]
	if end == "" {
		end = parsed.Flags["e"]
	}

	if start == "" && end == "" {
		return fmt.Errorf("at least one of --start (-s) or --end (-e) must be specified to trim")
	}

	if start != "" {
		if err := utils.ValidateTimeFormat(start); err != nil {
			return fmt.Errorf("invalid start time: %w", err)
		}
	}
	if end != "" {
		if err := utils.ValidateTimeFormat(end); err != nil {
			return fmt.Errorf("invalid end time: %w", err)
		}
	}

	// Build ffmpeg arguments
	ffmpegArgs := []string{"-y"}
	
	// Add input file
	ffmpegArgs = append(ffmpegArgs, "-i", inputFile)

	if start != "" {
		ffmpegArgs = append(ffmpegArgs, "-ss", start)
	}
	if end != "" {
		ffmpegArgs = append(ffmpegArgs, "-to", end)
	}

	// Use stream copy if extensions match to prevent unnecessary transcoding
	if strings.ToLower(filepath.Ext(inputFile)) == strings.ToLower(filepath.Ext(out)) {
		ffmpegArgs = append(ffmpegArgs, "-c", "copy")
	}

	ffmpegArgs = append(ffmpegArgs, out)

	fmt.Printf("Trimming %q to %q...\n", inputFile, out)
	if err := audio.RunFFmpeg(ffmpegArgs); err != nil {
		return err
	}

	fmt.Printf("Successfully trimmed audio and saved to %q\n", out)
	return nil
}

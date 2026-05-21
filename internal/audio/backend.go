package audio

import "fmt"

// RunFFmpeg wraps execution of the ffmpeg binary.
func RunFFmpeg(args []string) error {
	// TODO: Use os/exec to call `exec.Command("ffmpeg", args...)`.
	// TODO: Capture stdout and stderr.
	// TODO: Inspect the exit code and return a useful error on failure.
	return fmt.Errorf("RunFFmpeg is not implemented yet")
}

// RunFFprobe wraps execution of ffprobe when metadata is needed.
func RunFFprobe(args []string) error {
	return fmt.Errorf("RunFFprobe is not implemented yet")
}

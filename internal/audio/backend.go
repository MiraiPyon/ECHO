package audio

import (
	"bytes"
	"fmt"
	"os/exec"
)

// RunFFmpeg wraps execution of the ffmpeg binary.
var RunFFmpeg = func(args []string) error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg is not installed or not found in system PATH. Please install FFmpeg to run this command")
	}

	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg execution failed: %w (stderr: %s)", err, stderr.String())
	}
	return nil
}

// RunFFprobe wraps execution of ffprobe when metadata is needed and returns its stdout.
var RunFFprobe = func(args []string) ([]byte, error) {
	_, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, fmt.Errorf("ffprobe is not installed or not found in system PATH. Please install FFmpeg to run this command")
	}

	cmd := exec.Command("ffprobe", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffprobe execution failed: %w (stderr: %s)", err, stderr.String())
	}
	return stdout.Bytes(), nil
}

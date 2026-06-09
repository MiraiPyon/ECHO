package utils

import (
	"fmt"
	"os"
	"regexp"
)

var timeRegex = regexp.MustCompile(`^(?:(?:\d+:)?\d+:)?\d+(?:\.\d+)?$`)

// ValidateTimeFormat checks whether a time string such as "01:30" or "12:05:10" is supported.
func ValidateTimeFormat(timeStr string) error {
	if !timeRegex.MatchString(timeStr) {
		return fmt.Errorf("invalid time format %q. Expected formats: SS, MM:SS, HH:MM:SS (e.g., 30, 01:30, 12:05:10)", timeStr)
	}
	return nil
}

// ValidateFileExists ensures a path points to an existing file.
func ValidateFileExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("could not open file %q: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%q is a directory, not a file", path)
	}
	return nil
}

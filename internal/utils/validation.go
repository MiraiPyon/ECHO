package utils

import (
	"fmt"
	"os"
)

// ValidateTimeFormat checks whether a time string such as "01:30" or "12:05:10" is supported.
func ValidateTimeFormat(timeStr string) error {
	// TODO: Implement validation with a regexp or manual parsing.
	return fmt.Errorf("ValidateTimeFormat is not implemented yet")
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

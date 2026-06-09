package utils

import (
	"fmt"
	"strings"
)

// ParsedArgs holds the parsed positional arguments and named flags.
type ParsedArgs struct {
	Positional []string
	Flags      map[string]string
}

// ParseArgs parses a slice of string arguments.
// It supports flags in the format: --flag value, -flag value, --flag=value, -flag=value.
func ParseArgs(args []string) (*ParsedArgs, error) {
	parsed := &ParsedArgs{
		Positional: make([]string, 0),
		Flags:      make(map[string]string),
	}

	ignoreFlags := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if ignoreFlags {
			parsed.Positional = append(parsed.Positional, arg)
			continue
		}

		if arg == "--" {
			ignoreFlags = true
			continue
		}

		if strings.HasPrefix(arg, "-") && arg != "-" {
			parts := strings.SplitN(arg, "=", 2)
			flagName := parts[0]
			// Trim leading dashes (supports both -flag and --flag)
			flagName = strings.TrimLeft(flagName, "-")
			if flagName == "" {
				return nil, fmt.Errorf("empty flag name: %s", arg)
			}

			var flagVal string
			if len(parts) == 2 {
				flagVal = parts[1]
			} else {
				// The value is the next argument, unless it doesn't exist or is another flag
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					flagVal = args[i+1]
					i++
				} else {
					return nil, fmt.Errorf("flag --%s requires a value", flagName)
				}
			}
			parsed.Flags[flagName] = flagVal
		} else {
			parsed.Positional = append(parsed.Positional, arg)
		}
	}
	return parsed, nil
}

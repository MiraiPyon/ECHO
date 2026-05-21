# ECHO

ECHO is a terminal-based audio tool written in Go. It includes an interactive launcher and a custom player UI for previewing audio files without leaving the terminal.

## What Works Now

- `play` opens a terminal player for a single audio file.
- The player supports play/pause, seek, mute, volume changes, speed selection, and returning to the main launcher.
- The launcher can be used with the keyboard or mouse.

## Command Status

- `play` - implemented
- `trim` - planned
- `concat` - planned
- `extract` - planned
- `volume` - planned

## Requirements

- Go 1.25 or newer
- A terminal that supports interactive input
- Mouse support for the richer TUI experience

## Quick Start

```bash
go run ./cmd/echo/main.go
```

If you prefer a compiled binary:

```bash
make build
./bin/echo
```

Note: `echo` is also a shell builtin on many systems, so running the binary path directly is usually the safest option.

## Development

```bash
make test
make clean
```

`make test` runs the test suite under `tests/`.

Local audio fixtures for manual checks live in `testdata/` and are ignored by git.

## Project Layout

- `cmd/echo` - program entry point
- `internal/cli` - command parsing and help output
- `internal/tui` - interactive launcher
- `internal/commands/play` - the current audio player
- `internal/commands/{trim,concat,extract,volume}` - command shells for future work
- `internal/audio` - backend wrappers for ffmpeg and ffprobe
- `internal/utils` - shared validation helpers

## Status

This repository is ready for iterative development and GitHub hosting. The main UI and player are usable today, while the remaining command shells are in place for future implementation.

# Architecture Overview

This document describes a practical way to build ECHO so the project stays simple, maintainable, and easy to extend.

## Product Idea

ECHO is a terminal-first audio utility. The user should be able to run one command, pass a file and a few flags, and get a processed audio file back.

## Core Design Goals

- Keep commands short and predictable.
- Make audio operations available from the terminal.
- Use a clear flow from argument parsing to audio processing.
- Fail with readable error messages when input is invalid.

## Suggested High-Level Flow

1. The user runs a command such as `echo trim input.mp3 --start 00:30 --end 01:30 --out output.mp3`.
2. The CLI parser reads the command name and flags.
3. The app validates the input values.
4. The matching handler prepares the job.
5. ECHO calls the audio backend.
6. The backend writes the output file.
7. The CLI prints success or an error message.

## Suggested Internal Layers

### 1. CLI Layer

Responsible for:

- Reading command-line arguments
- Showing help text
- Routing to the correct command
- Validating required flags early

### 2. Command Layer

One package or module per command:

- `play`
- `trim`
- `concat`
- `extract`
- `volume`

Each command should contain only the logic specific to that action.

### 3. Audio Backend Layer

Responsible for the actual processing work.

For a Go CLI like ECHO, the most practical first implementation is to wrap FFmpeg and FFprobe.

This layer should handle:

- Running backend commands
- Building the correct arguments
- Capturing stdout/stderr
- Returning clear errors to the CLI layer

### 4. Utilities Layer

Shared helpers for:

- Time parsing
- Path validation
- File existence checks
- Temporary file handling

## Proposed Folder Structure

```text
ECHO/
  cmd/
    echo/
      main.go
  internal/
    app/
    cli/
    commands/
      play/
      trim/
      concat/
      extract/
      volume/
    audio/
    ffmpeg/
    utils/
  docs/
    README.md
    architecture.md
    commands.md
    roadmap.md
```

## What Each Command Does Internally

### `play`

- Accepts one audio file.
- Opens or streams it through the playback backend.
- Prints basic playback status.

### `trim`

- Reads the input file.
- Parses `--start` and `--end`.
- Extracts the selected time range.
- Writes a new audio file to `--out`.

### `concat`

- Collects multiple audio files in order.
- Verifies that the files are compatible or can be re-encoded.
- Joins them into one output file.

### `extract`

- Accepts a video file.
- Detects or selects the audio stream.
- Writes the audio stream to the output file.

### `volume`

- Reads the input file.
- Applies the requested gain or attenuation factor.
- Writes the adjusted audio to `--out`.

## Validation Rules To Expect

- `trim` requires both `--start` and `--end`.
- `trim` should reject invalid time ranges where start is greater than or equal to end.
- `concat` should require at least two input files.
- `extract` should require a video input file.
- `volume` should require a positive numeric `--level`.

## Example User Journey

### Trim a chorus

Input:

```bash
echo trim "song.mp3" --start 00:30 --end 01:30 --out "chorus.mp3"
```

Expected result:

- The app parses the times.
- The backend trims the selected range.
- A new file named `chorus.mp3` is created.

### Merge several clips

Input:

```bash
echo concat "a.mp3" "b.mp3" "c.mp3" --out "combined.mp3"
```

Expected result:

- The app reads the files in order.
- The backend joins them into one track.
- A single output file is written.

### Extract audio from a video

Input:

```bash
echo extract "video.mp4" --out "audio.mp3"
```

Expected result:

- The app finds the audio stream.
- The backend writes the extracted track to `audio.mp3`.

## Error Handling Philosophy

ECHO should be strict about invalid input and friendly in its messages.

Examples:

- Missing file path
- Invalid time format
- Output path not writable
- Unsupported file format
- Failed backend execution

When something goes wrong, the user should get a short message that says what failed and what to fix.


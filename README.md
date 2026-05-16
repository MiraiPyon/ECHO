# ECHO

**ECHO** (Efficient, Command-line, Harmony, Orchestrator) is a **Go**-based CLI tool that runs in the terminal and helps users work with audio quickly and conveniently.

## Goals

- Handle common audio tasks directly from the terminal.
- Avoid opening heavy editors for simple operations.
- Keep the command syntax short, memorable, and practical for both casual and technical users.

## Core Features

### 1. Play Audio

ECHO can play audio files directly from the terminal.

- Useful for quickly checking a file.
- Great when you want to preview or replay audio without opening a separate media player.

### 2. Trim Audio (`trim`)

The `trim` command lets users cut out a specific section from an audio file, such as extracting the chorus from 00:30 to 01:30.

**Usage:**

```bash
echo trim "input.mp3" --start 00:30 --end 01:30 --out "output.mp3"
```

### 3. Concatenate Audio (`concat`)

The `concat` command joins multiple audio files in order into one longer file.

**Usage:**

```bash
echo concat "file1.mp3" "file2.mp3" "file3.mp3" --out "output.mp3"
```

### 4. Extract Audio from Video (`extract`)

The `extract` command separates the audio track from a video file. This is especially useful when users want to save a good background song from a video.

**Usage:**

```bash
echo extract "video.mp4" --out "output.mp3"
```

### 5. Adjust Volume (`volume`)

The `volume` command changes the playback level of an audio file. For example, a quiet voice recording can be boosted to twice the original volume.

**Usage:**

```bash
echo volume "file.mp3" --level 2.0 --out "output.mp3"
```

## Parameter Conventions

- `input`: the source file to process.
- `--out`: the output file path.
- `--start`: the start time for trimming.
- `--end`: the end time for trimming.
- `--level`: the volume multiplier.

## Quick Examples

```bash
# Trim a segment from 00:30 to 01:30
echo trim "input.mp3" --start 00:30 --end 01:30 --out "chorus.mp3"

# Merge 3 audio files
echo concat "intro.mp3" "main.mp3" "outro.mp3" --out "full_track.mp3"

# Extract audio from a video
echo extract "movie.mp4" --out "audio.mp3"

# Increase volume by 2x
echo volume "recording.mp3" --level 2.0 --out "louder.mp3"
```

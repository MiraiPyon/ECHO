# Command Reference

This file documents the intended command set for ECHO.

## General Syntax

```bash
echo <command> [input files] [flags]
```

## 1. Play Audio

Planned syntax:

```bash
echo play "input.mp3"
```

Purpose:

- Play an audio file from the terminal.

Typical use:

- Quick preview
- Simple playback without opening another app

## 2. Trim Audio

```bash
echo trim "input.mp3" --start 00:30 --end 01:30 --out "output.mp3"
```

Purpose:

- Extract a specific time range from an audio file.

Important flags:

- `--start`: start time
- `--end`: end time
- `--out`: output file

Suggested validation:

- Start time must be before end time.
- Time format should be consistent, such as `mm:ss` or `hh:mm:ss`.

## 3. Concatenate Audio

```bash
echo concat "file1.mp3" "file2.mp3" "file3.mp3" --out "output.mp3"
```

Purpose:

- Join multiple audio files into a single file in the same order they are passed.

Important flags:

- `--out`: output file

Suggested validation:

- At least two input files are required.

## 4. Extract Audio From Video

```bash
echo extract "video.mp4" --out "output.mp3"
```

Purpose:

- Save the audio track from a video file as a separate audio file.

Important flags:

- `--out`: output file

Suggested validation:

- Input should be a video file.

## 5. Adjust Volume

```bash
echo volume "file.mp3" --level 2.0 --out "output.mp3"
```

Purpose:

- Increase or decrease the audio volume.

Important flags:

- `--level`: volume multiplier
- `--out`: output file

Suggested validation:

- `--level` must be a positive number.
- `1.0` means no change.
- `2.0` means twice as loud.
- `0.5` means half as loud.

## Common Output Behavior

Every command should ideally:

- Print a short success message
- Print the output file location
- Return a non-zero exit code on failure

## Help Behavior

ECHO should support helpful terminal output such as:

```bash
echo --help
echo trim --help
echo concat --help
```

The help text should show:

- What the command does
- Required inputs
- Available flags
- A short example


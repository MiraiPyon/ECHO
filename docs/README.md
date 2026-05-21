# ECHO Documentation

This folder explains how ECHO is expected to work before the codebase grows further.

ECHO is a Go-based CLI tool for common audio tasks in the terminal.

## What ECHO Does

- Play audio files
- Trim audio clips
- Concatenate audio files
- Extract audio from video files
- Adjust audio volume

## Recommended Command Style

The project is designed around a simple command pattern:

```bash
echo <command> [flags]
```

Planned commands:

- `play`
- `trim`
- `concat`
- `extract`
- `volume`

## Docs In This Folder

- [Architecture](architecture.md)
- [Commands](commands.md)
- [Roadmap](roadmap.md)

## Assumed Backend

For the first version, the simplest implementation path is to use FFmpeg and FFprobe under the hood.


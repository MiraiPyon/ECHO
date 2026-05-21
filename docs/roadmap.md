# Roadmap

This roadmap helps describe how ECHO can grow from an MVP into a more complete audio CLI tool.

## MVP

The first version should focus on the core user value:

- Command parsing
- Help output
- `play`
- `trim`
- `concat`
- `extract`
- `volume`
- Basic input validation
- Clear success and error messages

## Recommended First Backend

The easiest and most reliable approach is:

- Use Go for the CLI and orchestration
- Use FFmpeg for audio processing
- Use FFprobe when metadata or stream inspection is needed

## Nice-to-Have Features

- Batch processing
- Presets for common tasks
- Audio normalization
- Metadata reading and editing
- Progress indicators for long operations
- Dry-run mode
- Config file support

## Possible Future Commands

- `normalize`
- `info`
- `split`
- `merge`
- `tag`

## Design Principles For Future Growth

- Keep command names short.
- Keep flags consistent across commands.
- Prefer readable error messages over raw backend output.
- Add features only when they fit the main goal of fast terminal-based audio work.


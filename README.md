# Audio File Converter

This is a command line tool for converting audio files.

Currently, it supports the following conversions:

- Source File Types
    - .m4a
    - .wav
    - .flac
    - .aac
    - .ogg
    - .opus
    - .m4r

- Target File Types
    - .mp3

## Manual Installation

```
cd cli
make all
```

This creates binaries for the following platforms:

- darwin/arm64
- darwin/amd64
- linux/amd64
- linux/arm64

## Usage

```
./bin/converter-cli --source-dir /path/to/source/dir --target-dir /path/to/target/dir
```

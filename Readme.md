# ASCII-Art-Justify

ASCII-Art-Justify is a Go CLI project that renders text as ASCII art and supports horizontal alignment based on terminal width.

Supported banners:
- `standard`
- `shadow`
- `thinkertoy`

Supported alignment types:
- `left`
- `center`
- `right`
- `justify`

## Usage

Required justify format:

```bash
go run . [OPTION] [STRING] [BANNER]
```

Main option for this subject:

```bash
--align=<type>
```

Example:

```bash
go run . --align=right something standard
```

Invalid option formats print:

```text
Usage: go run . [OPTION] [STRING] [BANNER]

Example: go run . --align=right something standard
```

The program also supports a single `[STRING]` argument (default banner: `standard`).

## Terminal Adaptation

Alignment is applied using terminal width.
- Width is read from `COLUMNS` (with fallback width if not set).
- `center` and `right` add left padding.
- `justify` expands spaces between words so rendered output fills the line when possible.

Only text that fits terminal size is expected by tests.

## Optional Compatibility

This repository also keeps compatibility with previously implemented optional flags (when correctly formatted):
- `--color=<color>`
- `--out=<file.txt>`
- `--output=<file.txt>`
- `--font=<banner>`
- `--align=<type>`

## Examples

```bash
go run . --align=center "hello" standard
go run . --align=left "Hello There" standard
go run . --align=right "hello" shadow
go run . --align=justify "how are you" shadow
```

## Project Structure

- `main.go`: entry point
- `pipeline/pipeline.go`: run orchestration
- `pipeline/args.go`: CLI parsing and usage validation (`--align` included)
- `pipeline/alignment.go`: terminal-width alignment logic
- `pipeline/loadBanner.go`: banner loading
- `pipeline/tokenize.go`: tokenization
- `pipeline/renderLines.go`: ASCII art rendering
- `pipeline/colorFormating.go`: optional color support
- `pipeline/validateInput.go`: input validation
- `pipeline/writeOutput.go`: output writer
- `tests/args_test.go`: argument/option behavior tests
- `tests/alingment_test.go`: alignment behavior tests
- `tests/`: additional unit tests

## Testing

Run:

```bash
go test ./...
```

## Allowed Packages

Only Go standard library packages are used.

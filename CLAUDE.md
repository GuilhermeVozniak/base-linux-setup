# CLAUDE.md

## Build & Test

```bash
make build                          # Build to build/base-linux-setup
go test ./... -v                    # Run unit tests
go vet ./...                        # Static analysis
golangci-lint run                   # Lint (gocritic, gofmt, errcheck, etc.)
python3 -m json.tool scripts/*.json # Validate JSON presets
```

Always verify changes with: `go test ./... && go vet ./... && go build -o /dev/null .`

## Architecture

- `main.go` / `assets.go` — Entry point, `go:embed` for `scripts/*.json`
- `cmd/` — Cobra subcommands: `detect`, `list-presets`
- `internal/detector/` — Environment detection (neofetch + fallbacks)
- `internal/presets/` — JSON preset parsing, matching, validation, dependency sorting
- `internal/executor/` — Task execution (command, script, file, service types)
- `internal/ui/` — Interactive prompts (promptui)
- `scripts/*.json` — Preset config files, embedded at build time

## Key Design Rules

- **All presets are JSON** in `scripts/`. To add a preset, create JSON and rebuild. No Go changes needed.
- **`match` field** controls auto-detection. Most specific match (most fields) wins. `default.json` has no match and is the fallback.
- **`go:embed`** makes the binary standalone. No external files at runtime.
- **`--config` flag** loads an external JSON file, bypassing auto-detection.

## Task Types

- **command**: Runs `commands[]` via `sh -c`
- **script**: Writes `script` to temp file, executes it
- **file**: Creates file at `commands[0]` with content from `script`, permissions from `commands[1]`
- **service**: Runs `systemctl commands[1] commands[0]`

## Conventions

- Conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`
- Tasks using `sudo` internally must set `"elevated": true`
- Table-driven tests with `t.Run()` subtests
- Error messages start with lowercase verb (e.g., "failed to detect environment")
- `GetPreset()` and `GetAllPresets()` return `(T, error)` — always handle the error

## Dependencies

Go 1.21+, [cobra](https://github.com/spf13/cobra), [promptui](https://github.com/manifoldco/promptui), [color](https://github.com/fatih/color). Runtime: `neofetch`.

## Custom Commands

- `/test` — Run full verification pipeline (tests, vet, build, JSON validation)
- `/lint` — Run golangci-lint and auto-fix issues
- `/add-preset` — Guided creation of a new JSON preset
- `/release <version>` — Tag and push a release (triggers GitHub Actions)

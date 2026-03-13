# CLAUDE.md

## Project Overview

**base-linux-setup** is a Go CLI tool that auto-detects your Linux environment (via neofetch) and runs setup presets defined as JSON configuration files. All presets are embedded into the binary at build time — no external files needed at runtime.

Primary target: **Kali Linux on Raspberry Pi**, but any OS is supported by providing the appropriate JSON config.

## Architecture

```
base-linux-setup/
├── main.go                  # Entry point, wires embedded FS to presets package
├── assets.go                # go:embed for scripts/*.json → PresetFiles embed.FS
├── cmd/
│   ├── detect.go            # `detect` subcommand (neofetch-based)
│   └── list.go              # `list-presets` subcommand
├── internal/
│   ├── detector/detector.go # Environment detection (neofetch + fallbacks)
│   ├── presets/presets.go   # Preset/Task structs, matching, loading, validation
│   ├── executor/executor.go # Task execution (command, script, file, service types)
│   └── ui/ui.go             # Interactive prompts (promptui)
└── scripts/                 # JSON preset config files (embedded at build time)
    ├── kali-raspberry-pi.json
    ├── debian-base.json
    ├── ubuntu.json
    ├── arch.json
    └── default.json
```

## Key Design Decisions

- **All presets are JSON files** in `scripts/`. No presets are hardcoded in Go. To add or modify a preset, edit JSON — no Go changes needed.
- **`match` field** in each JSON preset controls auto-detection. The preset with the highest specificity (most matching fields) wins. `default.json` has no `match` field and serves as fallback.
- **`--config` flag** loads an external JSON preset file, bypassing auto-detection entirely.
- **go:embed** (`scripts/*.json`) embeds all presets into the binary. The binary is fully standalone.
- **No template variables** in presets — JSON is self-contained. Users edit values directly.

## Build & Test

```bash
make build                    # Build to build/base-linux-setup
make clean                    # Clean build artifacts
go build -o /dev/null .       # Quick compile check
go vet ./...                  # Static analysis
python3 -m json.tool scripts/*.json  # Validate all JSON presets
./build/base-linux-setup list-presets # Verify presets load correctly
./build/base-linux-setup --config scripts/kali-raspberry-pi.json  # Test external config
```

No test suite exists yet (`go test ./...` runs but there are no `_test.go` files).

## Adding a New Preset

1. Create `scripts/my-distro.json` with a `match` field:
   ```json
   {
     "name": "My Distro Setup",
     "environment": "My Distro",
     "description": "Setup tasks for My Distro",
     "match": { "distribution": "mydistro" },
     "tasks": [...]
   }
   ```
2. Rebuild (`make build`). That's it — no Go code changes.

## JSON Preset Schema

```json
{
  "name": "string (required)",
  "environment": "string (required)",
  "description": "string",
  "match": {
    "distribution": "string (case-insensitive contains match)",
    "os": "string (case-insensitive contains match)",
    "architecture": "string (case-insensitive contains match)",
    "is_raspberry_pi": "bool (exact match, null = don't care)"
  },
  "tasks": [
    {
      "name": "string (required)",
      "description": "string",
      "type": "command | script | file | service",
      "commands": ["string array"],
      "script": "string (script content for script/file types)",
      "elevated": "bool (requires sudo)",
      "optional": "bool (can be skipped)",
      "condition": "string (shell command, exit 0 = run task)",
      "depends_on": ["string array of task names"]
    }
  ]
}
```

## Task Types

- **command**: Runs each entry in `commands[]` via `sh -c`
- **script**: Writes `script` to a temp file, makes executable, runs it
- **file**: Creates file at `commands[0]` with content from `script`, permissions from `commands[1]`
- **service**: Runs `systemctl commands[1] commands[0]` (e.g., `["docker", "enable"]`)

## Dependencies

- Go 1.21+
- [cobra](https://github.com/spf13/cobra) — CLI framework
- [promptui](https://github.com/manifoldco/promptui) — interactive prompts
- [color](https://github.com/fatih/color) — terminal colors
- `neofetch` — required at runtime for environment detection

## Conventions

- Commit messages use conventional commits (`feat:`, `fix:`, `docs:`, etc.)
- Commands execute via `sh -c` (supports pipes, redirects, quoted args)
- All script tasks that use `sudo` internally should have `"elevated": true`

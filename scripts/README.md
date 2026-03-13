# Scripts Directory

This directory contains preset configuration files in JSON format. These files define the setup tasks for different operating systems and environments. All JSON files here are embedded into the binary at build time via `go:embed`.

## JSON Preset Format

Each preset file follows this JSON structure:

```json
{
  "name": "Preset Name",
  "environment": "Environment Description",
  "description": "Detailed description of what this preset does",
  "match": {
    "distribution": "distro-name",
    "os": "os-name",
    "architecture": "arch",
    "is_raspberry_pi": true
  },
  "tasks": [
    {
      "name": "Task Name",
      "description": "Task description",
      "type": "command|script|file|service",
      "commands": ["command1", "command2"],
      "script": "#!/bin/bash\necho 'script content'",
      "elevated": true|false,
      "optional": true|false,
      "condition": "shell-command-returns-0-to-run",
      "depends_on": ["other-task-name"]
    }
  ]
}
```

## Match Criteria

The `match` field controls which preset is auto-selected based on the detected environment. All fields are optional and use case-insensitive substring matching:

| Field | Type | Description |
|-------|------|-------------|
| `distribution` | string | Matches against detected distribution name (e.g., "kali", "ubuntu") |
| `os` | string | Matches against detected OS name |
| `architecture` | string | Matches against detected architecture (e.g., "amd64", "arm64") |
| `is_raspberry_pi` | bool | Exact match: `true` requires Pi, `false` requires non-Pi |

The preset with the **most matching fields** (highest specificity) wins. A preset without a `match` field serves as the default fallback.

## Task Types

### 1. Command Tasks

Execute shell commands sequentially via `sh -c`.

```json
{
  "name": "Install Packages",
  "description": "Install essential packages",
  "type": "command",
  "commands": ["sudo apt-get update", "sudo apt-get install -y git curl"],
  "elevated": true,
  "optional": false
}
```

### 2. Script Tasks

Execute bash scripts (written to a temp file and run).

```json
{
  "name": "Setup Environment",
  "description": "Configure development environment",
  "type": "script",
  "script": "#!/bin/bash\nset -e\necho 'Setting up environment...'\nexport PATH=$PATH:/usr/local/bin",
  "elevated": false,
  "optional": false
}
```

### 3. File Tasks

Create files with specific content.

```json
{
  "name": "Create Config File",
  "description": "Create application configuration",
  "type": "file",
  "commands": ["/path/to/file", "644"],
  "script": "config content goes here\nline 2\nline 3",
  "elevated": false,
  "optional": true
}
```

### 4. Service Tasks

Manage system services with systemctl.

```json
{
  "name": "Enable Service",
  "description": "Enable and start a service",
  "type": "service",
  "commands": ["service-name", "enable"],
  "elevated": true,
  "optional": false
}
```

Service actions: `start`, `stop`, `enable`, `disable`, `restart`, `reload`, `status`

## Field Descriptions

- **name**: Display name for the task (required)
- **description**: Detailed description shown to the user
- **type**: Task type (`command`, `script`, `file`, `service`)
- **commands**: Array of commands or parameters (usage varies by task type)
- **script**: Script content or file content (for script and file tasks)
- **elevated**: Whether the task requires sudo privileges
- **optional**: Whether the task can be skipped by the user
- **condition**: Shell command — task runs only if this exits with status 0
- **depends_on**: Array of task names that must run before this task

## Available Presets

### kali-raspberry-pi.json

Complete setup for Kali Linux on Raspberry Pi including:

- System updates
- Golang installation with architecture detection
- Development packages
- raspi-config installation
- I2C interface configuration
- Static IP address configuration
- mDNS/Avahi networking setup

### debian-base.json

Basic setup for Debian-based systems: system updates and essential development tools.

### ubuntu.json

Setup for Ubuntu systems: system updates and optional snap package installation.

### arch.json

Setup for Arch Linux: system updates and base development tools via pacman.

### default.json

Generic fallback preset with multi-package-manager commands for basic tool installation.

## Adding New Presets

1. Create a new JSON file in this directory
2. Add a `match` field so it auto-selects for the right environment
3. Rebuild with `make build` — no Go code changes needed
4. Verify with `./build/base-linux-setup list-presets`

You can also use any preset directly with `--config`:
```bash
./build/base-linux-setup --config scripts/my-preset.json
```

## Testing Presets

```bash
# 1. Validate JSON syntax
python3 -m json.tool scripts/my-preset.json

# 2. Build and verify
make build
./build/base-linux-setup list-presets

# 3. Test with external config (without rebuilding)
./build/base-linux-setup --config scripts/my-preset.json
```

## Notes

- Scripts should include proper error handling (`set -e`)
- Use absolute paths where necessary
- Test on target systems before deploying
- Consider making destructive operations optional
- Tasks that use `sudo` internally should set `"elevated": true`
- Document any system requirements or dependencies

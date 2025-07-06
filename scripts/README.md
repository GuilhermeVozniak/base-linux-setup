# Scripts Directory

This directory contains preset configuration files in JSON format. These files define the setup tasks for different operating systems and environments.

## JSON Preset Format

Each preset file follows this JSON structure:

```json
{
  "name": "Preset Name",
  "environment": "Environment Description",
  "description": "Detailed description of what this preset does",
  "tasks": [
    {
      "name": "Task Name",
      "description": "Task description",
      "type": "command|script|file|service",
      "commands": ["command1", "command2"],
      "script": "#!/bin/bash\necho 'script content'",
      "elevated": true|false,
      "optional": true|false
    }
  ]
}
```

## Task Types

### 1. Command Tasks
Execute shell commands sequentially.

```json
{
  "name": "Install Packages",
  "description": "Install essential packages",
  "type": "command",
  "commands": [
    "sudo apt-get update",
    "sudo apt-get install -y git curl"
  ],
  "elevated": true,
  "optional": false
}
```

### 2. Script Tasks
Execute bash scripts.

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

- **name**: Display name for the task
- **description**: Detailed description shown to the user
- **type**: Task type (`command`, `script`, `file`, `service`)
- **commands**: Array of commands or parameters (usage varies by task type)
- **script**: Script content or file content (for script and file tasks)
- **elevated**: Whether the task requires sudo privileges
- **optional**: Whether the task can be skipped by the user

## Available Presets

### kali-raspberry-pi.json
Complete setup for Kali Linux on Raspberry Pi including:
- System updates
- Golang installation with architecture detection
- Development packages
- I2C interface configuration
- Docker installation and setup
- Development aliases

## Adding New Presets

1. Create a new JSON file in this directory
2. Follow the JSON format above
3. Update `internal/presets/presets.go` to load your preset:

```go
// In the appropriate detection function
if isYourEnvironment(env) {
    if preset, err := loadPresetFromJSON("your-preset.json"); err == nil {
        return preset
    }
    // Fallback preset if needed
}
```

## Testing Presets

You can test your JSON presets by:

1. Building the application: `make build`
2. Listing presets: `./build/base-linux-setup list-presets`
3. Running in dry-run mode (future feature)

## Notes

- Scripts should include proper error handling (`set -e`)
- Use absolute paths where necessary
- Test on target systems before deploying
- Consider making destructive operations optional
- Document any system requirements or dependencies 
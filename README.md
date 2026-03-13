# Base Linux Setup

A powerful CLI tool that automatically detects your Linux environment and provides customizable setup presets for different distributions and hardware configurations.

## Features

- 🔍 **Automatic Environment Detection**: Uses `neofetch` to detect OS, distribution, architecture, and hardware
- 🍓 **Raspberry Pi Support**: Special detection and optimized presets for Raspberry Pi devices
- 🎯 **Targeted Presets**: Specific configurations for:
  - Kali Linux on Raspberry Pi
  - Ubuntu
  - Debian-based systems
  - Arch Linux
  - Generic Linux fallback
- 🛠️ **Customizable Tasks**: Add, remove, or modify setup tasks interactively
- 🔧 **Multiple Task Types**: Support for commands, scripts, file operations, and service management
- 🎨 **Beautiful CLI Interface**: Colorful output with interactive prompts
- 💾 **Backup System**: Automatic backup of important configuration files
- 🧪 **Dry Run Mode**: Test configurations without making changes

## Installation

### Prerequisites

- Go 1.21 or higher
- `neofetch` (for environment detection)

```bash
# Install neofetch on Debian/Ubuntu/Kali
sudo apt-get install neofetch

# Install neofetch on Arch Linux
sudo pacman -S neofetch

# Install neofetch on macOS
brew install neofetch
```

### Build from Source

```bash
git clone https://github.com/GuilhermeVozniak/base-linux-setup.git
cd base-linux-setup
go build -o base-linux-setup .
```

### Using Make

```bash
make build          # Build the application
make install        # Install to /usr/local/bin
make clean          # Clean build artifacts
make release        # Create cross-platform release builds
```

### Pre-built Binaries

After building with `make build` or `make release`, you can run the executable directly:

```bash
# Run from build directory
./build/base-linux-setup

# Or use a specific platform binary (after make release)
./build/release/base-linux-setup-linux-amd64    # Linux x86_64
./build/release/base-linux-setup-linux-arm64    # Linux ARM64 (Raspberry Pi 4, etc.)
./build/release/base-linux-setup-linux-arm      # Linux ARM (Raspberry Pi 3, etc.)
./build/release/base-linux-setup-darwin-amd64   # macOS Intel
./build/release/base-linux-setup-darwin-arm64   # macOS Apple Silicon
```

### Download and Run (No Build Required)

If you have a pre-built binary, you can run it directly:

```bash
# Make executable and run
chmod +x base-linux-setup-linux-arm64
./base-linux-setup-linux-arm64

# Or move to your PATH
sudo mv base-linux-setup-linux-arm64 /usr/local/bin/base-linux-setup
base-linux-setup
```

## Usage

### Basic Usage

Run the setup with automatic environment detection:

```bash
# If built with make build
./build/base-linux-setup

# If using a release binary
./build/release/base-linux-setup-linux-arm64

# If installed to PATH
base-linux-setup
```

### Available Commands

```bash
# Detect current environment
./build/base-linux-setup detect

# List all available presets
./build/base-linux-setup list-presets

# Show version information
./build/base-linux-setup --version

# Show help
./build/base-linux-setup --help

# Generate shell completion
./build/base-linux-setup completion bash > base-linux-setup.bash
source base-linux-setup.bash
```

### Quick Start with Example Script

For a guided experience, use the included example script:

```bash
# Make executable and run
chmod +x example.sh
./example.sh

# The script will:
# 1. Build the application if needed
# 2. Show available commands
# 3. Let you choose what to run
# 4. Provide helpful guidance
```

### Interactive Setup Process

1. **Environment Detection**: The tool automatically detects your system
2. **Preset Selection**: Shows the recommended preset for your environment
3. **Task Customization**: Option to modify, add, or remove tasks
4. **Confirmation**: Review final task list before execution
5. **Execution**: Run tasks with progress tracking and error handling

## Kali Linux Raspberry Pi Preset

The comprehensive preset includes **7 tasks** for complete Raspberry Pi setup:

### 1. Update and Upgrade System

- Updates package lists (`apt-get update`)
- Upgrades all installed packages (`apt-get upgrade`)
- Performs distribution upgrade (`apt-get dist-upgrade`)

### 2. Install Golang

- Detects system architecture (ARM64, ARM7, etc.)
- Downloads and installs latest Go version (1.21.5)
- Configures PATH and GOPATH in `.bashrc`
- Creates Go workspace structure

### 3. Install Required System Packages

- Build tools (`build-essential`)
- Version control (`git`)
- Network tools (`curl`, `wget`)
- Text editors (`vim`, `nano`)
- Programming languages (`python3`, `nodejs`)
- System utilities (`htop`, `tree`)
- I2C development tools (`i2c-tools`, `libi2c-dev`)

### 4. Enable I2C Interface

- Enables I2C in `/boot/config.txt`
- Loads I2C kernel modules (`i2c-bcm2708`, `i2c-dev`)
- Installs I2C development tools and Python libraries
- Configures user permissions for I2C group

### 5. Configure Fixed IP Address

- Sets static IP address to `192.168.1.100/24`
- Configures gateway as `192.168.1.1`
- Uses Google DNS servers (`8.8.8.8`, `8.8.4.4`)
- Backs up original `dhcpcd.conf` configuration
- Includes optional Wi-Fi configuration

### 6. Install and Configure mDNS

- Installs Avahi daemon for mDNS/Zeroconf networking
- Configures hostname as `kali-pi.local`
- Enables network discovery and service publishing
- Makes Pi accessible via `kali-pi.local`
- Supports SSH access using `ssh user@kali-pi.local`

### 7. Install raspi-config

- Adds Raspbian repository for official Pi tools
- Installs `raspi-config` configuration tool
- Includes additional tools (`rpi-update`, `raspberrypi-bootloader`)
- Creates compatibility symbolic links
- Enables `sudo raspi-config` for system configuration

> **Note**: The raspi-config installation follows the approach documented in community guides ([GitHub](https://github.com/EmilGus/install_raspi-config)) for adding Raspberry Pi tools to Kali Linux.

## Customization

### Adding Custom Tasks

During the setup process, you can add custom tasks by:

1. Selecting "Customize tasks"
2. Choosing "Add custom tasks"
3. Providing task details:
   - Task name and description
   - Commands to execute
   - Whether elevated privileges are needed

### Task Types

- **Command**: Execute shell commands
- **Script**: Run bash scripts
- **File**: Create or modify files
- **Service**: Manage system services

### Example Custom Task

```bash
Task name: Install VS Code
Commands: wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg; sudo install -o root -g root -m 644 packages.microsoft.gpg /etc/apt/trusted.gpg.d/; sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/trusted.gpg.d/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'; sudo apt update; sudo apt install code
Elevated privileges: Yes
```

## Configuration

### Environment Variables

- `GOPATH`: Set automatically during Go installation
- `PATH`: Updated to include Go binaries

### Backup Location

Configuration backups are stored in:

```
~/.config/base-linux-setup/backups/
```

## Safety Features

- **Automatic Backups**: Important files are backed up before modification
- **Error Handling**: Graceful error handling with continue/cancel options
- **Dry Run Mode**: Test configurations without making changes
- **Permission Checks**: Validates required permissions before execution

## Development

### Project Structure

```
base-linux-setup/
├── cmd/                    # CLI commands
│   ├── detect.go          # Environment detection command
│   └── list.go            # List presets command
├── internal/              # Internal packages
│   ├── detector/          # Environment detection logic
│   ├── presets/           # Preset loading, matching, and validation
│   ├── ui/               # User interface components
│   └── executor/          # Task execution engine
├── scripts/               # JSON preset configurations (embedded at build time)
│   ├── kali-raspberry-pi.json  # Kali Linux Raspberry Pi preset
│   ├── debian-base.json        # Debian preset
│   ├── ubuntu.json             # Ubuntu preset
│   ├── arch.json               # Arch Linux preset
│   ├── default.json            # Generic fallback preset
│   └── README.md               # JSON format documentation
├── main.go               # Main entry point
├── assets.go             # go:embed for all preset JSON files
├── go.mod               # Go module dependencies
├── Makefile             # Build automation
├── LICENSE              # MIT License
├── example.sh           # Usage examples
└── README.md            # This file
```

### Adding New Presets

All presets are JSON files in the `scripts/` directory. No Go code changes are needed.

1. Create a new JSON file in the `scripts/` directory
2. Add a `match` field to control auto-detection (see below)
3. Rebuild with `make build`
4. Test with `./build/base-linux-setup list-presets`

**JSON Format Example:**

```json
{
  "name": "My Custom Preset",
  "environment": "Custom Linux",
  "description": "Custom setup tasks",
  "match": {
    "distribution": "mydistro"
  },
  "tasks": [
    {
      "name": "Install Tools",
      "type": "command",
      "commands": ["sudo apt-get install -y htop"],
      "elevated": true,
      "optional": false
    }
  ]
}
```

**Match fields** (all optional, case-insensitive substring match):
- `distribution` — matches against detected distribution name
- `os` — matches against detected OS name
- `architecture` — matches against detected architecture
- `is_raspberry_pi` — `true`/`false` for exact Raspberry Pi match

The preset with the most matching fields wins. A preset without a `match` field serves as the default fallback.

You can also use any preset file directly with `--config`:
```bash
./build/base-linux-setup --config scripts/kali-raspberry-pi.json
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Troubleshooting

### Common Issues

**neofetch not found**

```bash
sudo apt-get install neofetch
```

**Permission denied**

```bash
# Make sure you have sudo privileges
sudo -v
```

**Go installation fails**

```bash
# Check internet connectivity
ping -c 1 golang.org
```

**I2C interface not working**

```bash
# Reboot after enabling I2C
sudo reboot
```

### Platform-Specific Notes

**Raspberry Pi Users:**

- Use `base-linux-setup-linux-arm64` for Raspberry Pi 4/5 (64-bit)
- Use `base-linux-setup-linux-arm` for Raspberry Pi 3 and older (32-bit)

**Linux Users:**

- Use `base-linux-setup-linux-amd64` for standard Linux x86_64 systems
- Use `base-linux-setup-linux-arm64` for ARM64 systems

**macOS Users:**

- Use `base-linux-setup-darwin-arm64` for Apple Silicon Macs (M1/M2/M3)
- Use `base-linux-setup-darwin-amd64` for Intel Macs

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- Create an issue on GitHub for bugs or feature requests
- Check the troubleshooting section for common problems
- Review the examples for usage patterns

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [promptui](https://github.com/manifoldco/promptui) for interactive prompts
- Styled with [color](https://github.com/fatih/color) for terminal colors

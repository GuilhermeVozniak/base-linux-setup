# Usage Guide

This guide covers how to use Base Linux Setup effectively, from basic commands to advanced customization.

## Basic Commands

### Help and Information
```bash
# Show help
base-linux-setup --help

# Show version information
base-linux-setup --version

# Generate shell completion
base-linux-setup completion bash > /etc/bash_completion.d/base-linux-setup
base-linux-setup completion zsh > ~/.zsh_completions/_base-linux-setup
```

### Environment Detection
```bash
# Detect current environment
base-linux-setup detect

# Example output:
# Environment Information:
#   OS: Linux
#   Distribution: kali
#   Version: 2023.4
#   Architecture: aarch64
#   Hardware: Raspberry Pi
#   Kernel: 6.1.0-kali7-arm64
#   🍓 Raspberry Pi detected!
```

### List Available Presets
```bash
# Show all presets
base-linux-setup list-presets

# Example output:
# Available Presets:
# 
# ▶ Kali Linux - Raspberry Pi
#   Environment: Kali Linux (Raspberry Pi)
#   Description: Complete setup for Kali Linux on Raspberry Pi
#   Tasks: 4
#     1. Update and Upgrade System
#     2. Install Golang
#     3. Install Required System Packages
#     4. Enable I2C Interface
```

### Interactive Setup
```bash
# Run the main setup process
base-linux-setup

# This will:
# 1. Detect your environment
# 2. Show the recommended preset
# 3. Allow customization
# 4. Execute selected tasks
```

## Interactive Setup Process

### Step 1: Environment Detection
The tool automatically detects your system and displays:
- Operating System
- Distribution and version
- Architecture (x86_64, ARM64, ARM)
- Hardware type (Raspberry Pi detection)
- Kernel version

### Step 2: Preset Selection
Based on your environment, the tool will:
- Show the best matching preset
- Display the preset description
- List all included tasks

### Step 3: Customization Options
You'll be presented with choices:

**Use preset as-is**
- Proceeds with all default tasks
- Recommended for most users

**Customize tasks (add/remove/modify)**
- Allows detailed task modification
- Lets you add custom tasks
- Good for advanced users

**Cancel setup**
- Exits without making changes

### Step 4: Task Customization (if selected)

For each task, you can:
- **Include this task** - Add to execution list
- **Skip this task** - Remove from execution
- **Mark as optional** - Include but allow failure

#### Adding Custom Tasks
If you choose to add custom tasks, you'll be prompted for:

**Task Information:**
- **Task name**: Display name for the task
- **Description**: Detailed explanation (optional)
- **Commands**: Shell commands to execute (separate with `;`)
- **Elevated privileges**: Whether the task needs `sudo`

**Example Custom Task:**
```
Task name: Install VS Code
Description: Install Visual Studio Code editor
Commands: wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg; sudo install packages.microsoft.gpg /etc/apt/trusted.gpg.d/; sudo sh -c 'echo "deb [signed-by=/etc/apt/trusted.gpg.d/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'; sudo apt update; sudo apt install code
Elevated privileges: Yes
```

### Step 5: Confirmation
Review the final task list:
- ✓ indicates required tasks
- ? indicates optional tasks
- Total task count is shown

### Step 6: Execution
Tasks are executed in order with:
- Progress tracking (Task X/Y)
- Command output in real-time
- Execution time for each task
- Success/failure indicators

#### Error Handling
If a task fails, you can choose to:
- **Continue with remaining tasks**
- **Cancel setup**

## Working with Presets

### Kali Linux - Raspberry Pi Preset

This preset is automatically selected for Kali Linux running on Raspberry Pi hardware.

**Included Tasks:**
1. **Update and Upgrade System**
   ```bash
   sudo apt-get update
   sudo apt-get upgrade -y
   sudo apt-get dist-upgrade -y
   ```

2. **Install Golang**
   - Detects ARM architecture automatically
   - Downloads appropriate Go version
   - Configures PATH and GOPATH
   - Creates Go workspace

3. **Install Required System Packages**
   ```bash
   sudo apt-get install -y build-essential git curl wget vim nano
   sudo apt-get install -y python3 python3-pip nodejs npm
   sudo apt-get install -y htop tree i2c-tools libi2c-dev python3-smbus
   ```

4. **Enable I2C Interface**
   - Enables I2C in `/boot/config.txt`
   - Loads kernel modules
   - Configures user permissions
   - Requires reboot to take effect

**Optional Tasks:**
- Docker installation and configuration
- Development aliases creation
- Service management

### Ubuntu Preset

Optimized for Ubuntu desktop and server systems.

**Included Tasks:**
- System updates via `apt`
- Essential development tools
- Snap package installations (optional)

### Debian Preset

Basic setup for Debian-based systems.

**Included Tasks:**
- Package manager updates
- Essential development packages

### Arch Linux Preset

Tailored for Arch Linux systems.

**Included Tasks:**
- System updates via `pacman`
- Base development tools from official repositories

### Generic Linux Preset

Fallback preset for unsupported distributions.

**Included Tasks:**
- Multi-package-manager update commands
- Basic tool installation with fallbacks

## Advanced Usage

### Understanding Task Types

#### Command Tasks
Execute shell commands sequentially:
```json
{
  "name": "Install Packages",
  "type": "command",
  "commands": [
    "sudo apt-get update",
    "sudo apt-get install -y git"
  ],
  "elevated": true
}
```

#### Script Tasks
Run complex bash scripts:
```json
{
  "name": "Setup Environment",
  "type": "script",
  "script": "#!/bin/bash\nset -e\necho 'Configuring environment...'\nexport PATH=$PATH:/usr/local/bin",
  "elevated": false
}
```

#### File Tasks
Create configuration files:
```json
{
  "name": "Create Config",
  "type": "file",
  "commands": ["/path/to/file", "644"],
  "script": "configuration content here",
  "elevated": false
}
```

#### Service Tasks
Manage systemd services:
```json
{
  "name": "Enable Service",
  "type": "service",
  "commands": ["docker", "enable"],
  "elevated": true
}
```

### Backup and Safety

#### Automatic Backups
Before making system changes, the tool automatically backs up:
- `/etc/fstab`
- `/boot/config.txt` (Raspberry Pi)
- `/etc/modules`
- `~/.bashrc`
- `~/.profile`

**Backup Location**: `~/.config/base-linux-setup/backups/`

**Backup Format**: `filename.timestamp`

#### Safety Features
- **Permission validation** before execution
- **Network connectivity** checks
- **Disk space** verification
- **Error handling** with user choices
- **Graceful cancellation** at any point

### Configuration Files

#### User Preferences
The tool stores minimal configuration in:
`~/.config/base-linux-setup/`

#### Custom Presets
You can create custom presets by:
1. Creating JSON files in the `scripts/` directory
2. Following the format in `scripts/README.md`
3. Modifying detection logic if needed

## Tips and Best Practices

### Before Running Setup
1. **Backup important data** independently
2. **Ensure stable internet** connection
3. **Have sudo privileges** available
4. **Run on a fresh system** when possible

### During Setup
1. **Read task descriptions** carefully
2. **Skip unknown or risky tasks** if unsure
3. **Monitor command output** for errors
4. **Don't interrupt critical operations**

### After Setup
1. **Reboot if requested** (especially for I2C changes)
2. **Source shell configurations** or restart terminal
3. **Verify installations** worked correctly
4. **Test new functionality** before relying on it

### Customization Strategies

**Conservative Approach:**
- Use presets as-is initially
- Add custom tasks gradually
- Test on non-critical systems first

**Advanced Approach:**
- Modify existing presets heavily
- Create custom JSON presets
- Contribute improvements back to project

## Example Workflows

### Basic Kali Pi Setup
```bash
# 1. Install the tool
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-arm64
chmod +x base-linux-setup-linux-arm64

# 2. Run setup with defaults
./base-linux-setup-linux-arm64
# Choose "Use preset as-is"

# 3. Reboot after I2C setup
sudo reboot
```

### Custom Development Environment
```bash
# 1. Run setup
./base-linux-setup

# 2. Choose "Customize tasks"

# 3. Add custom tasks like:
# - Install specific editors
# - Configure development tools
# - Set up project directories

# 4. Execute customized setup
```

### Testing New Preset
```bash
# 1. Check current detection
./base-linux-setup detect

# 2. List available presets
./base-linux-setup list-presets

# 3. Run setup to test preset
./base-linux-setup
```

## Troubleshooting Usage

### Common Issues

**"No preset found for your environment"**
- The tool will use the generic fallback preset
- Consider contributing a preset for your environment
- Manually add tasks during customization

**"Task failed to execute"**
- Check the error message carefully
- Verify internet connectivity
- Ensure proper permissions
- Choose to continue or cancel

**"Environment detection incomplete"**
- Install or update `neofetch`
- Check if `neofetch --stdout` works manually
- Some fields may show "Unknown" on unsupported systems

### Getting Help

1. **Check command output** for specific error messages
2. **Review the [[Troubleshooting]]** guide
3. **Search [existing issues](https://github.com/GuilhermeVozniak/base-linux-setup/issues)**
4. **Create a [bug report](https://github.com/GuilhermeVozniak/base-linux-setup/issues/new?template=bug_report.md)** with:
   - Output of `./base-linux-setup detect`
   - Full error messages
   - System information

---

**Next Steps**: Learn about **[[Configuration]]** options or explore **[[Preset Development]]** to create custom setups! 
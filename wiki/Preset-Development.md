# Preset Development Guide

This comprehensive guide covers creating, testing, and maintaining presets for Base Linux Setup. Whether you're adding support for a new OS or customizing existing setups, this guide has you covered.

## Overview

Presets define the setup tasks for different operating systems and environments. All presets are **JSON configuration files** in the `scripts/` directory — no Go code changes are needed to add or modify presets. JSON files are embedded into the binary at build time via `go:embed`.

## JSON Preset Structure

Each preset follows this structure:

```json
{
  "name": "Display Name",
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
      "script": "script content for script/file tasks",
      "elevated": true,
      "optional": false,
      "condition": "shell-command-returns-0-to-run",
      "depends_on": ["other-task-name"]
    }
  ]
}
```

### Field Reference

#### Preset Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Display name for the preset |
| `environment` | string | ✅ | Environment description |
| `description` | string | ✅ | Detailed preset description |
| `match` | object | ❌ | Auto-detection criteria (omit for default fallback) |
| `tasks` | array | ✅ | Array of task objects |

#### Match Criteria Fields

| Field | Type | Description |
|-------|------|-------------|
| `distribution` | string | Case-insensitive substring match against detected distro |
| `os` | string | Match against detected OS name |
| `architecture` | string | Match against detected architecture (e.g., `amd64`, `arm64`) |
| `is_raspberry_pi` | bool | Exact match: `true` requires Pi, `false` requires non-Pi |

The preset with the **most matching fields** (highest specificity) wins. A preset without a `match` field serves as the default fallback.

#### Task Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Task display name |
| `description` | string | ❌ | Task description |
| `type` | string | ✅ | Task type: `command`, `script`, `file`, `service` |
| `commands` | array | Varies | Commands or parameters (usage varies by type) |
| `script` | string | ❌ | Script content (for script/file tasks) |
| `elevated` | boolean | ✅ | Whether task requires sudo |
| `optional` | boolean | ✅ | Whether task can be skipped |
| `condition` | string | ❌ | Shell command — task runs only if this exits with status 0 |
| `depends_on` | array | ❌ | Array of task names that must run before this task |

### Task Types in Detail

#### 1. Command Tasks

Execute shell commands sequentially.

```json
{
  "name": "Install Development Tools",
  "description": "Install essential development packages",
  "type": "command",
  "commands": [
    "sudo apt-get update",
    "sudo apt-get install -y build-essential git curl wget",
    "sudo apt-get install -y vim nano htop tree"
  ],
  "elevated": true,
  "optional": false
}
```

**Best Practices:**
- Break complex operations into multiple commands
- Use package manager's batch installation when possible
- Include error handling where appropriate

#### 2. Script Tasks

Execute bash scripts with full control.

```json
{
  "name": "Install Golang",
  "description": "Install Go with architecture detection",
  "type": "script",
  "script": "#!/bin/bash\nset -e\n\n# Remove old installation\nsudo rm -rf /usr/local/go\n\n# Detect architecture\nARCH=$(uname -m)\ncase $ARCH in\n    \"x86_64\") GOARCH=\"amd64\" ;;\n    \"aarch64\"|\"arm64\") GOARCH=\"arm64\" ;;\n    \"armv7l\"|\"armv6l\") GOARCH=\"armv6l\" ;;\n    *) echo \"Unsupported architecture: $ARCH\"; exit 1 ;;\nesac\n\n# Download and install\nGO_VERSION=\"1.21.5\"\nwget https://golang.org/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz\nsudo tar -C /usr/local -xzf go${GO_VERSION}.linux-${GOARCH}.tar.gz\nrm go${GO_VERSION}.linux-${GOARCH}.tar.gz\n\n# Configure environment\necho 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc\necho 'export GOPATH=$HOME/go' >> ~/.bashrc\necho 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc\n\n# Create workspace\nmkdir -p $HOME/go/{bin,pkg,src}\n\necho \"Go installed successfully!\"",
  "elevated": false,
  "optional": false
}
```

**Best Practices:**
- Always include `#!/bin/bash` and `set -e`
- Use proper error handling
- Include informative echo statements
- Test scripts independently before integration

#### 3. File Tasks

Create configuration files with custom content.

```json
{
  "name": "Create Development Aliases",
  "description": "Create useful command aliases",
  "type": "file",
  "commands": ["/home/user/.bash_aliases", "644"],
  "script": "# Development aliases\nalias ll='ls -alF'\nalias la='ls -A'\nalias l='ls -CF'\nalias ..='cd ..'\nalias ...='cd ../..'\n\n# Git aliases\nalias gs='git status'\nalias ga='git add'\nalias gc='git commit'\nalias gp='git push'\nalias gl='git log --oneline'\n\n# System aliases\nalias update='sudo apt update && sudo apt upgrade'\nalias install='sudo apt install'",
  "elevated": false,
  "optional": true
}
```

**Commands Array for File Tasks:**
- `commands[0]`: File path
- `commands[1]`: File permissions (octal, e.g., "644", "755")

**Best Practices:**
- Use appropriate file permissions
- Consider user-specific vs system-wide files
- Include helpful comments in configuration files

#### 4. Service Tasks

Manage systemd services.

```json
{
  "name": "Enable Docker Service",
  "description": "Enable and start Docker service",
  "type": "service",
  "commands": ["docker", "enable"],
  "elevated": true,
  "optional": false
}
```

**Commands Array for Service Tasks:**
- `commands[0]`: Service name
- `commands[1]`: Action (`start`, `stop`, `enable`, `disable`, `restart`, `reload`, `status`)

**Available Actions:**
- `start` - Start the service
- `stop` - Stop the service
- `enable` - Enable service for auto-start
- `disable` - Disable auto-start
- `restart` - Restart the service
- `reload` - Reload service configuration
- `status` - Check service status

## Creating a New Preset

### Step 1: Research the Target Environment

Before creating a preset, gather information about:

**System Information:**
- OS name and common detection strings
- Package manager (apt, yum, pacman, etc.)
- Service manager (systemd, openrc, etc.)
- Architecture support
- Special requirements or configurations

**Common Setup Tasks:**
- System update commands
- Essential package lists
- Development tool installations
- Service configurations
- User customizations

### Step 2: Create the JSON File

1. **Create file** in `scripts/` directory:
   ```bash
   touch scripts/my-environment.json
   ```

2. **Start with basic structure**:
   ```json
   {
     "name": "My Environment Setup",
     "environment": "My Linux Distribution",
     "description": "Complete setup for My Linux Distribution",
     "tasks": []
   }
   ```

3. **Add tasks incrementally** and test each one

### Step 3: Add Match Criteria

Add a `match` field to your JSON so the preset auto-selects for the right environment:

```json
{
  "name": "My Environment Setup",
  "environment": "My Linux Distribution",
  "description": "Complete setup for My Linux Distribution",
  "match": {
    "distribution": "mydist"
  },
  "tasks": [...]
}
```

No Go code changes are needed. The `match` field uses case-insensitive substring matching. The system automatically selects the preset with the most matching criteria (highest specificity).

### Step 4: Rebuild and Test

Rebuild with `make build` — the new JSON file is automatically embedded into the binary.

### Step 5: Testing

Test your preset thoroughly:

```bash
# 1. Validate JSON syntax
python3 -m json.tool scripts/my-environment.json

# 2. Build and test
make build
./build/base-linux-setup list-presets

# 3. Test detection
./build/base-linux-setup detect

# 4. Test on target system
./build/base-linux-setup
```

## Example: Creating a CentOS Preset

Let's walk through creating a complete preset for CentOS:

### 1. Research CentOS

- **Package Manager**: `yum` (older) or `dnf` (newer)
- **Service Manager**: `systemd`
- **Common Detection**: Distribution contains "centos" or "rhel"
- **Architecture**: Supports x86_64, ARM64

### 2. Create JSON File

```json
{
  "name": "CentOS Setup",
  "environment": "CentOS Linux",
  "description": "Complete development setup for CentOS systems",
  "tasks": [
    {
      "name": "Update System",
      "description": "Update package lists and upgrade system",
      "type": "command",
      "commands": [
        "sudo yum update -y"
      ],
      "elevated": true,
      "optional": false
    },
    {
      "name": "Install EPEL Repository",
      "description": "Install Extra Packages for Enterprise Linux",
      "type": "command",
      "commands": [
        "sudo yum install -y epel-release"
      ],
      "elevated": true,
      "optional": false
    },
    {
      "name": "Install Development Tools",
      "description": "Install essential development packages",
      "type": "command",
      "commands": [
        "sudo yum groupinstall -y \"Development Tools\"",
        "sudo yum install -y git curl wget vim nano",
        "sudo yum install -y python3 python3-pip nodejs npm"
      ],
      "elevated": true,
      "optional": false
    },
    {
      "name": "Configure Firewall",
      "description": "Configure firewalld for development",
      "type": "script",
      "script": "#!/bin/bash\nset -e\n\n# Enable firewalld\nsudo systemctl enable firewalld\nsudo systemctl start firewalld\n\n# Open common development ports\nsudo firewall-cmd --permanent --add-port=3000/tcp\nsudo firewall-cmd --permanent --add-port=8000/tcp\nsudo firewall-cmd --permanent --add-port=8080/tcp\n\n# Reload firewall\nsudo firewall-cmd --reload\n\necho \"Firewall configured for development\"",
      "elevated": false,
      "optional": true
    }
  ]
}
```

### 3. Add Match Criteria

The `match` field in the JSON handles auto-detection — no Go code needed:

```json
{
  "match": {
    "distribution": "centos"
  }
}
```

This will auto-select on any system where the detected distribution contains "centos" (case-insensitive).

For more specific matching, add more criteria:
```json
{
  "match": {
    "distribution": "centos",
    "architecture": "arm64"
  }
}
```

### 4. Rebuild

```bash
make build
./build/base-linux-setup list-presets
```

## Advanced Preset Features

### Conditional Tasks

Use the `condition` field to run tasks only when a shell command succeeds:

```json
{
  "name": "Install CentOS 8 Packages",
  "description": "Packages specific to CentOS 8",
  "type": "command",
  "commands": ["sudo dnf install -y container-tools"],
  "elevated": true,
  "optional": false,
  "condition": "grep -q '8' /etc/redhat-release"
}
```

### Task Dependencies

Use `depends_on` to control execution order:

```json
{
  "name": "Configure Service",
  "type": "command",
  "commands": ["sudo systemctl enable myservice"],
  "elevated": true,
  "optional": false,
  "depends_on": ["Install Service Package"]
}

### Architecture-Specific Tasks

Handle different architectures within tasks:

```json
{
  "name": "Install Architecture-Specific Tools",
  "type": "script",
  "script": "#!/bin/bash\nset -e\n\nARCH=$(uname -m)\ncase $ARCH in\n    \"x86_64\")\n        sudo yum install -y x86-specific-package\n        ;;\n    \"aarch64\")\n        sudo yum install -y arm64-specific-package\n        ;;\n    *)\n        echo \"Unsupported architecture: $ARCH\"\n        exit 1\n        ;;\nesac",
  "elevated": false,
  "optional": false
}
```

### Environment Variables

Use environment variables in scripts:

```json
{
  "name": "Setup User Environment",
  "type": "script",
  "script": "#!/bin/bash\nset -e\n\n# Get current user\nCURRENT_USER=${USER:-$(whoami)}\nHOME_DIR=${HOME:-/home/$CURRENT_USER}\n\n# Create directories\nmkdir -p $HOME_DIR/{projects,bin,scripts}\n\n# Set permissions\nchown $CURRENT_USER:$CURRENT_USER $HOME_DIR/{projects,bin,scripts}\n\necho \"Environment setup for user: $CURRENT_USER\"",
  "elevated": false,
  "optional": false
}
```

## Testing and Validation

### Automated Testing

Create test scripts for your presets:

```bash
#!/bin/bash
# test-preset.sh

set -e

echo "Testing preset on $(lsb_release -d)"

# Build application
make build

# Test detection
echo "=== Environment Detection ==="
./build/base-linux-setup detect

# Test preset listing
echo "=== Preset Listing ==="
./build/base-linux-setup list-presets | grep -A 10 "My Environment"

# Test JSON validation
echo "=== JSON Validation ==="
python3 -m json.tool scripts/my-environment.json > /dev/null
echo "JSON is valid"

echo "All tests passed!"
```

### Manual Testing Checklist

Before submitting a preset:

- [ ] **JSON syntax** is valid
- [ ] **All required fields** are present
- [ ] **Task commands** are correct for the target OS
- [ ] **File paths** are appropriate for the environment
- [ ] **Service names** exist on the target system
- [ ] **Package names** are correct
- [ ] **Architecture detection** works correctly
- [ ] **Elevated privileges** are set appropriately
- [ ] **Optional tasks** are marked correctly
- [ ] **Error handling** is implemented in scripts
- [ ] **Match criteria** are set correctly for the target environment

### Testing on Target Systems

**Virtual Machine Testing:**
```bash
# Create VM with target OS
# Install base-linux-setup
# Run through complete setup process
# Verify all tasks execute successfully
# Check installed software works correctly
```

**Container Testing:**
```dockerfile
FROM centos:8
RUN yum update -y
COPY base-linux-setup /usr/local/bin/
RUN base-linux-setup list-presets
# Test specific commands
```

## Best Practices

### General Guidelines

1. **Start Simple**: Begin with basic tasks, add complexity gradually
2. **Test Early**: Validate each task as you add it
3. **Document Everything**: Include clear descriptions for all tasks
4. **Handle Errors**: Use `set -e` in scripts and check command results
5. **Be Consistent**: Follow established patterns in existing presets

### Security Considerations

1. **Minimize Elevated Tasks**: Only use `sudo` when absolutely necessary
2. **Validate Inputs**: Check for expected conditions before execution
3. **Use Trusted Sources**: Download software from official repositories
4. **Verify Downloads**: Use checksums when possible
5. **Avoid Hardcoded Credentials**: Never include passwords or keys

### Performance Optimization

1. **Batch Operations**: Group similar commands together
2. **Use Package Groups**: Install multiple packages in single commands
3. **Cache Downloads**: Avoid redundant downloads
4. **Parallel Safe**: Ensure tasks can run in sequence safely

### Maintenance

1. **Version Compatibility**: Test with different OS versions
2. **Package Updates**: Keep package lists current
3. **Deprecation Handling**: Update commands as tools evolve
4. **Documentation**: Keep descriptions accurate and helpful

## Contributing Presets

### Submission Process

1. **Create the preset** following this guide
2. **Test thoroughly** on target systems
3. **Create a pull request** with:
   - The JSON preset file (with `match` field)
   - Documentation updates
   - Test results

### Pull Request Template

When submitting a preset, include:

```markdown
## New Preset: [OS Name]

### Environment Information
- **OS**: [OS Name and Version]
- **Package Manager**: [apt/yum/pacman/etc.]
- **Architecture**: [Supported architectures]
- **Hardware**: [Special hardware support]

### Tasks Included
- [ ] System updates
- [ ] Essential packages
- [ ] Development tools
- [ ] [Other tasks...]

### Testing
- [ ] JSON syntax validated
- [ ] Tested on target OS
- [ ] All tasks execute successfully
- [ ] Detection logic works
- [ ] Documentation updated

### Additional Notes
[Any special considerations or requirements]
```

## Getting Help

### Resources

- **[[Contributing]]** - General contribution guidelines
- **[Example Presets](https://github.com/GuilhermeVozniak/base-linux-setup/tree/main/scripts)** - Reference implementations
- **[Issue Templates](https://github.com/GuilhermeVozniak/base-linux-setup/issues/new/choose)** - Bug reports and feature requests

### Community Support

- **[Discussions](https://github.com/GuilhermeVozniak/base-linux-setup/discussions)** - Ask questions and share ideas
- **[Issues](https://github.com/GuilhermeVozniak/base-linux-setup/issues)** - Report bugs or request features
- **[Pull Requests](https://github.com/GuilhermeVozniak/base-linux-setup/pulls)** - Contribute code and presets

---

**Next Steps**: Start creating your preset or explore the **[[API Reference]]** for advanced development! 
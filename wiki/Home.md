# Base Linux Setup Wiki

Welcome to the **Base Linux Setup** Wiki! This comprehensive guide will help you get started with the CLI tool that automatically detects your Linux environment and provides customizable setup presets.

## 🎯 What is Base Linux Setup?

Base Linux Setup is a powerful CLI tool that:
- **Automatically detects** your Linux distribution, architecture, and hardware
- **Provides targeted presets** for different environments (Kali Linux, Ubuntu, Raspberry Pi, etc.)
- **Allows customization** of setup tasks through an interactive interface
- **Supports multiple task types** including commands, scripts, file operations, and service management
- **Uses JSON configuration** for easy preset management

## 🚀 Quick Start

### 1. Download and Install
```bash
# For Raspberry Pi 4+ (ARM64)
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-arm64
chmod +x base-linux-setup-linux-arm64
./base-linux-setup-linux-arm64

# For standard Linux x86_64
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-amd64
chmod +x base-linux-setup-linux-amd64
./base-linux-setup-linux-amd64
```

### 2. Run Setup
```bash
# Interactive setup
./base-linux-setup

# Detect environment only
./base-linux-setup detect

# List available presets
./base-linux-setup list-presets
```

## 📚 Documentation Index

### Getting Started
- **[[Installation]]** - Download, install, and configure Base Linux Setup
- **[[Usage Guide]]** - Basic usage, commands, and interactive features
- **[[Configuration]]** - Environment detection and preset selection

### Development
- **[[Development Guide]]** - Setting up development environment
- **[[Preset Development]]** - Creating custom presets and JSON configurations
- **[[API Reference]]** - Internal APIs and package documentation

### Support
- **[[Troubleshooting]]** - Common issues and solutions
- **[[FAQ]]** - Frequently asked questions
- **[[Contributing]]** - How to contribute to the project

## 🎯 Supported Environments

| Environment | Status | Tasks |
|-------------|--------|-------|
| **Kali Linux (Raspberry Pi)** | ✅ Full Support | 4 core tasks |
| **Ubuntu** | ✅ Full Support | 2 core tasks |
| **Debian** | ✅ Full Support | 2 core tasks |
| **Arch Linux** | ✅ Full Support | 2 core tasks |
| **Generic Linux** | ✅ Fallback | 2 basic tasks |

## 🔧 Key Features

### Environment Detection
- Uses `neofetch` for comprehensive system information
- Detects OS, distribution, architecture, and hardware
- Special Raspberry Pi detection and optimization

### Task Types
1. **Command Tasks** - Execute shell commands
2. **Script Tasks** - Run bash scripts with full control
3. **File Tasks** - Create configuration files with custom content
4. **Service Tasks** - Manage systemd services

### Interactive Customization
- Add, remove, or modify tasks before execution
- Create custom tasks on the fly
- Optional task support for flexible configurations

### Safety Features
- Automatic backup of important configuration files
- Error handling with continue/cancel options
- Dry-run mode for testing (planned feature)

## 🏗️ Architecture

```
base-linux-setup/
├── cmd/                 # CLI commands (detect, list-presets)
├── internal/            # Core packages
│   ├── detector/        # Environment detection using neofetch
│   ├── presets/         # Preset management and JSON loading
│   ├── ui/             # Interactive user interface
│   └── executor/        # Task execution engine
├── scripts/            # JSON preset configurations
└── .github/            # CI/CD workflows and templates
```

## 🤝 Community

### Getting Help
- 🐛 **Bug Reports**: [Create an issue](https://github.com/GuilhermeVozniak/base-linux-setup/issues/new?template=bug_report.md)
- ✨ **Feature Requests**: [Request a feature](https://github.com/GuilhermeVozniak/base-linux-setup/issues/new?template=feature_request.md)
- 🛠️ **New Presets**: [Request OS support](https://github.com/GuilhermeVozniak/base-linux-setup/issues/new?template=preset_request.md)
- 💬 **Discussions**: [Join the conversation](https://github.com/GuilhermeVozniak/base-linux-setup/discussions)

### Contributing
We welcome contributions! See our **[[Contributing]]** guide for:
- Development setup
- Code guidelines
- Testing procedures
- Pull request process

## 📋 Latest Release

Check out the [latest release](https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest) for:
- Cross-platform binaries
- Changelog
- Installation instructions
- SHA256 checksums

## 🎉 What's Next?

- Explore the **[[Installation]]** guide to get started
- Check out **[[Preset Development]]** to create custom configurations
- Join our community and help improve Base Linux Setup!

---

*This wiki is collaboratively maintained. Feel free to contribute improvements and additions!* 
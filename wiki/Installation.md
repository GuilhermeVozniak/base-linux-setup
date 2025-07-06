# Installation Guide

This guide covers multiple ways to install and set up Base Linux Setup on your system.

## Prerequisites

### Required
- **Linux-based operating system** (Kali, Ubuntu, Debian, Arch, etc.)
- **Internet connection** for downloading and package installation
- **Terminal access** with basic command knowledge

### Optional (for development)
- **Go 1.21+** for building from source
- **Git** for cloning the repository
- **Make** for using build automation

## Installation Methods

### Method 1: Download Pre-built Binary (Recommended)

The easiest way to get started is downloading a pre-built binary from our releases.

#### Step 1: Choose Your Platform

| Platform | Architecture | Binary Name |
|----------|-------------|-------------|
| **Raspberry Pi 4+** | ARM64 | `base-linux-setup-linux-arm64` |
| **Raspberry Pi 3** | ARM | `base-linux-setup-linux-arm` |
| **Standard Linux** | x86_64 | `base-linux-setup-linux-amd64` |
| **macOS Intel** | x86_64 | `base-linux-setup-darwin-amd64` |
| **macOS Apple Silicon** | ARM64 | `base-linux-setup-darwin-arm64` |

#### Step 2: Download

**Option A: Latest Release (Recommended)**
```bash
# For Raspberry Pi 4+ (ARM64)
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-arm64

# For standard Linux (x86_64)
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-amd64

# For Raspberry Pi 3 and older (ARM)
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-arm
```

**Option B: Specific Version**
```bash
# Replace v1.0.0 with desired version
VERSION="v1.0.0"
PLATFORM="linux-arm64"  # Change to your platform
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/download/$VERSION/base-linux-setup-$PLATFORM
```

#### Step 3: Make Executable and Run
```bash
# Make executable
chmod +x base-linux-setup-*

# Run directly
./base-linux-setup-linux-arm64

# Or move to PATH for global access
sudo mv base-linux-setup-linux-arm64 /usr/local/bin/base-linux-setup
base-linux-setup --version
```

#### Step 4: Verify Installation (Optional)
```bash
# Download checksum file
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-arm64.sha256

# Verify integrity
sha256sum -c base-linux-setup-linux-arm64.sha256
```

### Method 2: Build from Source

For developers or users who want the latest features:

#### Step 1: Install Prerequisites
```bash
# Ubuntu/Debian/Kali
sudo apt-get install git golang-go make

# Arch Linux
sudo pacman -S git go make

# CentOS/RHEL/Fedora
sudo yum install git golang make  # or dnf instead of yum
```

#### Step 2: Clone Repository
```bash
git clone https://github.com/GuilhermeVozniak/base-linux-setup.git
cd base-linux-setup
```

#### Step 3: Build
```bash
# Simple build
make build

# Or build manually
go build -o base-linux-setup .

# Build for multiple platforms
make release
```

#### Step 4: Install (Optional)
```bash
# Install to /usr/local/bin
make install

# Or copy manually
sudo cp build/base-linux-setup /usr/local/bin/
```

### Method 3: One-Line Installer Script

**⚠️ Coming Soon**: We're working on a one-line installer script:
```bash
# Planned feature
curl -sSL https://raw.githubusercontent.com/GuilhermeVozniak/base-linux-setup/main/install.sh | bash
```

## Platform-Specific Notes

### Raspberry Pi
- **Pi 4/5 (64-bit)**: Use `base-linux-setup-linux-arm64`
- **Pi 3 and older (32-bit)**: Use `base-linux-setup-linux-arm`
- **Check your architecture**: Run `uname -m` to confirm

### Kali Linux
```bash
# Install neofetch (required for detection)
sudo apt-get update
sudo apt-get install neofetch

# Then install base-linux-setup
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-arm64
chmod +x base-linux-setup-linux-arm64
./base-linux-setup-linux-arm64
```

### Ubuntu/Debian
```bash
# Install neofetch
sudo apt update
sudo apt install neofetch

# Continue with standard installation
```

### Arch Linux
```bash
# Install neofetch
sudo pacman -S neofetch

# Continue with standard installation
```

### macOS (Development/Testing)
```bash
# Install neofetch via Homebrew
brew install neofetch

# Download macOS binary
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-darwin-arm64
chmod +x base-linux-setup-darwin-arm64
./base-linux-setup-darwin-arm64
```

## Post-Installation Setup

### 1. Verify Installation
```bash
# Check version
base-linux-setup --version

# Test environment detection
base-linux-setup detect

# List available presets
base-linux-setup list-presets
```

### 2. Install neofetch (if not already installed)
Base Linux Setup requires `neofetch` for environment detection:

```bash
# Ubuntu/Debian/Kali
sudo apt-get install neofetch

# Arch Linux
sudo pacman -S neofetch

# CentOS/RHEL (EPEL required)
sudo yum install epel-release
sudo yum install neofetch

# Fedora
sudo dnf install neofetch

# macOS
brew install neofetch
```

### 3. Run Your First Setup
```bash
# Start interactive setup
base-linux-setup

# Or explore options first
base-linux-setup --help
```

## Troubleshooting Installation

### Common Issues

**"Command not found"**
```bash
# Make sure the binary is executable
chmod +x base-linux-setup-*

# Check if it's in your PATH
echo $PATH

# Add to PATH temporarily
export PATH=$PATH:$(pwd)

# Add to PATH permanently
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

**"Permission denied"**
```bash
# Make the binary executable
chmod +x base-linux-setup-*

# Check permissions
ls -la base-linux-setup-*
```

**"neofetch not found"**
```bash
# Install neofetch for your distribution
sudo apt-get install neofetch    # Debian/Ubuntu/Kali
sudo pacman -S neofetch          # Arch
sudo dnf install neofetch        # Fedora
brew install neofetch            # macOS
```

**"Architecture mismatch"**
```bash
# Check your architecture
uname -m

# Download the correct binary:
# x86_64 → base-linux-setup-linux-amd64
# aarch64/arm64 → base-linux-setup-linux-arm64
# armv7l → base-linux-setup-linux-arm
```

### Getting Help

If you encounter issues:

1. **Check our [[Troubleshooting]]** guide
2. **Search [existing issues](https://github.com/GuilhermeVozniak/base-linux-setup/issues)**
3. **Create a [bug report](https://github.com/GuilhermeVozniak/base-linux-setup/issues/new?template=bug_report.md)**

## Updating

### Update Pre-built Binary
```bash
# Download new version
wget https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest/download/base-linux-setup-linux-arm64

# Replace old version
chmod +x base-linux-setup-linux-arm64
sudo mv base-linux-setup-linux-arm64 /usr/local/bin/base-linux-setup

# Verify update
base-linux-setup --version
```

### Update from Source
```bash
cd base-linux-setup
git pull origin main
make build
make install
```

## Uninstallation

### Remove Binary
```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/base-linux-setup

# If using locally
rm base-linux-setup-*
```

### Remove Configuration (Optional)
```bash
# Remove backup directory
rm -rf ~/.config/base-linux-setup/
```

---

**Next Steps**: Check out the **[[Usage Guide]]** to learn how to use Base Linux Setup! 
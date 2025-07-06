package presets

import (
	"base-linux-setup/internal/detector"
	"strings"
)

// Task represents a single setup task
type Task struct {
	Name        string
	Description string
	Type        string   // "command", "script", "file", "service"
	Commands    []string
	Script      string
	Elevated    bool // requires sudo
	Optional    bool
}

// Preset represents a collection of tasks for a specific environment
type Preset struct {
	Name        string
	Environment string
	Description string
	Tasks       []Task
}

// GetPreset returns the appropriate preset for the given environment
func GetPreset(env *detector.Environment) *Preset {
	// Check for Kali Linux on Raspberry Pi
	if isKaliRaspberryPi(env) {
		return getKaliRaspberryPiPreset()
	}

	// Check for other Debian-based systems
	if isDebianBased(env) {
		return getDebianBasePreset()
	}

	// Check for Ubuntu
	if isUbuntu(env) {
		return getUbuntuPreset()
	}

	// Check for Arch Linux
	if isArch(env) {
		return getArchPreset()
	}

	return nil
}

// GetDefaultPreset returns a basic preset for unknown environments
func GetDefaultPreset() *Preset {
	return &Preset{
		Name:        "Basic Linux Setup",
		Environment: "Generic Linux",
		Description: "Basic setup tasks for generic Linux systems",
		Tasks: []Task{
			{
				Name:        "Update Package List",
				Description: "Update the package manager cache",
				Type:        "command",
				Commands:    []string{"sudo apt-get update || sudo yum update || sudo pacman -Sy"},
				Elevated:    true,
			},
			{
				Name:        "Install Basic Tools",
				Description: "Install essential development tools",
				Type:        "command",
				Commands:    []string{"sudo apt-get install -y curl wget git || sudo yum install -y curl wget git || sudo pacman -S curl wget git"},
				Elevated:    true,
			},
		},
	}
}

// GetAllPresets returns all available presets
func GetAllPresets() []*Preset {
	return []*Preset{
		getKaliRaspberryPiPreset(),
		getDebianBasePreset(),
		getUbuntuPreset(),
		getArchPreset(),
		GetDefaultPreset(),
	}
}

// getKaliRaspberryPiPreset returns the preset for Kali Linux on Raspberry Pi
func getKaliRaspberryPiPreset() *Preset {
	return &Preset{
		Name:        "Kali Linux - Raspberry Pi",
		Environment: "Kali Linux (Raspberry Pi)",
		Description: "Complete setup for Kali Linux on Raspberry Pi with development tools",
		Tasks: []Task{
			{
				Name:        "Update and Upgrade System",
				Description: "Update package lists and upgrade all installed packages",
				Type:        "command",
				Commands: []string{
					"sudo apt-get update",
					"sudo apt-get upgrade -y",
					"sudo apt-get dist-upgrade -y",
				},
				Elevated: true,
			},
			{
				Name:        "Install Golang",
				Description: "Install Go programming language",
				Type:        "script",
				Script: `
#!/bin/bash
set -e

# Remove old Go installation
sudo rm -rf /usr/local/go

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    "x86_64") GOARCH="amd64" ;;
    "aarch64"|"arm64") GOARCH="arm64" ;;
    "armv7l"|"armv6l") GOARCH="armv6l" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Download and install Go
GO_VERSION="1.21.5"
wget https://golang.org/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-${GOARCH}.tar.gz
rm go${GO_VERSION}.linux-${GOARCH}.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

# Create GOPATH directory
mkdir -p $HOME/go/{bin,pkg,src}

echo "Go installed successfully!"
echo "Please run 'source ~/.bashrc' or restart your terminal"
`,
				Elevated: false,
			},
			{
				Name:        "Install Required System Packages",
				Description: "Install essential development and system packages",
				Type:        "command",
				Commands: []string{
					"sudo apt-get install -y build-essential",
					"sudo apt-get install -y git curl wget",
					"sudo apt-get install -y vim nano",
					"sudo apt-get install -y python3 python3-pip",
					"sudo apt-get install -y nodejs npm",
					"sudo apt-get install -y htop tree",
					"sudo apt-get install -y i2c-tools",
					"sudo apt-get install -y libi2c-dev",
					"sudo apt-get install -y python3-smbus",
				},
				Elevated: true,
			},
			{
				Name:        "Enable I2C Interface",
				Description: "Enable I2C interface for hardware communication",
				Type:        "script",
				Script: `
#!/bin/bash
set -e

# Enable I2C in config.txt
if ! grep -q "dtparam=i2c_arm=on" /boot/config.txt; then
    echo "dtparam=i2c_arm=on" | sudo tee -a /boot/config.txt
fi

# Load I2C kernel modules
if ! grep -q "i2c-bcm2708" /etc/modules; then
    echo "i2c-bcm2708" | sudo tee -a /etc/modules
fi

if ! grep -q "i2c-dev" /etc/modules; then
    echo "i2c-dev" | sudo tee -a /etc/modules
fi

# Load modules now
sudo modprobe i2c-bcm2708
sudo modprobe i2c-dev

# Add user to i2c group
sudo usermod -a -G i2c $USER

echo "I2C interface enabled!"
echo "Please reboot your system for changes to take effect"
`,
				Elevated: false,
			},
			{
				Name:        "Install Docker",
				Description: "Install Docker container platform",
				Type:        "script",
				Script: `
#!/bin/bash
set -e

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
rm get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER

# Enable Docker service
sudo systemctl enable docker
sudo systemctl start docker

echo "Docker installed successfully!"
echo "Please log out and log back in for group changes to take effect"
`,
				Elevated: false,
				Optional: true,
			},
		},
	}
}

// Helper functions for environment detection
func isKaliRaspberryPi(env *detector.Environment) bool {
	return (strings.Contains(strings.ToLower(env.Distribution), "kali") ||
		strings.Contains(strings.ToLower(env.OS), "kali")) &&
		env.IsRaspberryPi
}

func isDebianBased(env *detector.Environment) bool {
	dist := strings.ToLower(env.Distribution)
	return strings.Contains(dist, "debian") ||
		strings.Contains(dist, "ubuntu") ||
		strings.Contains(dist, "kali")
}

func isUbuntu(env *detector.Environment) bool {
	return strings.Contains(strings.ToLower(env.Distribution), "ubuntu")
}

func isArch(env *detector.Environment) bool {
	return strings.Contains(strings.ToLower(env.Distribution), "arch")
}

// Additional preset functions
func getDebianBasePreset() *Preset {
	return &Preset{
		Name:        "Debian Base",
		Environment: "Debian Linux",
		Description: "Basic setup for Debian-based systems",
		Tasks: []Task{
			{
				Name:        "Update System",
				Description: "Update and upgrade system packages",
				Type:        "command",
				Commands: []string{
					"sudo apt-get update",
					"sudo apt-get upgrade -y",
				},
				Elevated: true,
			},
			{
				Name:        "Install Essential Packages",
				Description: "Install essential development tools",
				Type:        "command",
				Commands: []string{
					"sudo apt-get install -y build-essential git curl wget vim",
				},
				Elevated: true,
			},
		},
	}
}

func getUbuntuPreset() *Preset {
	return &Preset{
		Name:        "Ubuntu Setup",
		Environment: "Ubuntu Linux",
		Description: "Setup for Ubuntu systems",
		Tasks: []Task{
			{
				Name:        "Update System",
				Description: "Update package lists and upgrade system",
				Type:        "command",
				Commands: []string{
					"sudo apt update",
					"sudo apt upgrade -y",
				},
				Elevated: true,
			},
			{
				Name:        "Install Snap Packages",
				Description: "Install useful snap packages",
				Type:        "command",
				Commands: []string{
					"sudo snap install code --classic",
					"sudo snap install discord",
				},
				Elevated: true,
				Optional: true,
			},
		},
	}
}

func getArchPreset() *Preset {
	return &Preset{
		Name:        "Arch Linux Setup",
		Environment: "Arch Linux",
		Description: "Setup for Arch Linux systems",
		Tasks: []Task{
			{
				Name:        "Update System",
				Description: "Update system packages",
				Type:        "command",
				Commands: []string{
					"sudo pacman -Syu --noconfirm",
				},
				Elevated: true,
			},
			{
				Name:        "Install Base Development Tools",
				Description: "Install essential development packages",
				Type:        "command",
				Commands: []string{
					"sudo pacman -S --noconfirm base-devel git curl wget vim",
				},
				Elevated: true,
			},
		},
	}
} 
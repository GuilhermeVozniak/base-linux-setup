package presets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"base-linux-setup/internal/detector"
)

// EmbeddedJSONGetter is a function type for getting embedded JSON data
type EmbeddedJSONGetter func(filename string) ([]byte, error)

// embeddedJSONGetter is the function to get embedded JSON data
var embeddedJSONGetter EmbeddedJSONGetter

// SetEmbeddedJSONGetter sets the function to get embedded JSON data
func SetEmbeddedJSONGetter(getter EmbeddedJSONGetter) {
	embeddedJSONGetter = getter
}

// Task represents a single setup task
type Task struct {
	Name        string
	Description string
	Type        string // "command", "script", "file", "service"
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
	// Try to load from embedded JSON first
	if preset, err := loadPresetFromEmbeddedJSON("kali-raspberry-pi.json"); err == nil {
		return preset
	}
	
	// Fallback to external JSON file (for development)
	if preset, err := loadPresetFromJSON("kali-raspberry-pi.json"); err == nil {
		return preset
	}
	
	// Final fallback to hardcoded preset
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
				Type:        "command",
				Commands: []string{
					"sudo apt-get install -y golang-go",
				},
				Elevated: true,
			},
			{
				Name:        "Install Required System Packages",
				Description: "Install essential development and system packages",
				Type:        "command",
				Commands: []string{
					"sudo apt-get install -y build-essential git curl wget vim python3 python3-pip nodejs npm htop tree i2c-tools libi2c-dev python3-smbus",
				},
				Elevated: true,
			},
			{
				Name:        "Enable I2C Interface",
				Description: "Enable I2C interface for hardware communication",
				Type:        "command",
				Commands: []string{
					"sudo raspi-config nonint do_i2c 0",
				},
				Elevated: true,
			},
		},
	}
}

// loadPresetFromEmbeddedJSON loads a preset from embedded JSON data
func loadPresetFromEmbeddedJSON(filename string) (*Preset, error) {
	if embeddedJSONGetter == nil {
		return nil, fmt.Errorf("embedded JSON getter not set")
	}
	
	// Get embedded JSON data
	data, err := embeddedJSONGetter(filename)
	if err != nil {
		return nil, err
	}
	
	// Parse JSON
	var preset Preset
	if err := json.Unmarshal(data, &preset); err != nil {
		return nil, fmt.Errorf("failed to parse embedded preset JSON: %v", err)
	}
	
	return &preset, nil
}

// loadPresetFromJSON loads a preset from a JSON file (fallback for development)
func loadPresetFromJSON(filename string) (*Preset, error) {
	// Get the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}
	
	// Look for scripts directory relative to executable
	scriptDir := filepath.Join(filepath.Dir(execPath), "scripts")
	
	// If not found, try relative to source code (for development)
	if _, err := os.Stat(scriptDir); os.IsNotExist(err) {
		// Get current file's directory for development
		_, currentFile, _, _ := runtime.Caller(0)
		projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
		scriptDir = filepath.Join(projectRoot, "scripts")
	}
	
	filePath := filepath.Join(scriptDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("preset file not found: %s", filePath)
	}
	
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read preset file: %v", err)
	}
	
	// Parse JSON
	var preset Preset
	if err := json.Unmarshal(data, &preset); err != nil {
		return nil, fmt.Errorf("failed to parse preset JSON: %v", err)
	}
	
	return &preset, nil
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

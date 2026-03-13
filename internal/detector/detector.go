package detector

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Environment represents the detected system environment
type Environment struct {
	OS            string
	Distribution  string
	Version       string
	Architecture  string
	Hardware      string
	Kernel        string
	IsRaspberryPi bool
	RawOutput     string
}

// DetectEnvironment detects the current environment using neofetch
func DetectEnvironment() (*Environment, error) {
	// Check if neofetch is available
	if _, err := exec.LookPath("neofetch"); err != nil {
		return nil, fmt.Errorf("neofetch is not installed. Please install it first: sudo apt-get install neofetch")
	}

	// Run neofetch with minimal output
	cmd := exec.Command("neofetch", "--stdout")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run neofetch: %v", err)
	}

	rawOutput := string(output)
	env := &Environment{
		RawOutput: rawOutput,
	}

	// Parse neofetch output
	lines := strings.Split(rawOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract information using regex patterns
		// Note: neofetch uses "OS:" and "Distro:" for distribution info,
		// and "Host:" for hardware. It does not output "Architecture:" directly,
		// so architecture is detected via the uname fallback below.
		if extractField(line, `^OS:\s*(.+)`, &env.OS) ||
			extractField(line, `^Distro:\s*(.+)`, &env.Distribution) ||
			extractField(line, `^Kernel:\s*(.+)`, &env.Kernel) ||
			extractField(line, `^Host:\s*(.+)`, &env.Hardware) {
			continue
		}

		// Check for Raspberry Pi indicators
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "raspberry") ||
			strings.Contains(lowerLine, "raspi") ||
			strings.Contains(lowerLine, "raspberry pi") {
			env.IsRaspberryPi = true
		}

		// Extract version info
		if strings.Contains(line, "OS:") || strings.Contains(line, "Distro:") {
			if versionMatch := regexp.MustCompile(`(\d+\.\d+|\d+)`).FindString(line); versionMatch != "" {
				env.Version = versionMatch
			}
		}
	}

	// Fallback detection methods
	if env.OS == "" {
		env.OS = detectOSFallback()
	}
	if env.Distribution == "" {
		env.Distribution = detectDistributionFallback()
	}
	if env.Architecture == "" {
		env.Architecture = detectArchitectureFallback()
	}

	// Detect hardware via /proc/cpuinfo if not already set from neofetch Host field
	if env.Hardware == "" {
		env.Hardware = detectHardware()
	} else if strings.Contains(strings.ToLower(env.Hardware), "raspberry") {
		env.IsRaspberryPi = true
	}

	return env, nil
}

// extractField extracts a field from a line using regex
func extractField(line, pattern string, target *string) bool {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 && target != nil {
		*target = strings.TrimSpace(matches[1])
		return true
	}
	return false
}

// detectOSFallback provides fallback OS detection
func detectOSFallback() string {
	if cmd := exec.Command("uname", "-s"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			return strings.TrimSpace(string(output))
		}
	}
	return "Unknown"
}

// detectDistributionFallback provides fallback distribution detection
func detectDistributionFallback() string {
	// Check /etc/os-release
	if cmd := exec.Command("cat", "/etc/os-release"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "ID=") {
					return strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
				}
			}
		}
	}

	// Check /etc/debian_version for Debian-based systems
	if cmd := exec.Command("cat", "/etc/debian_version"); cmd != nil {
		if _, err := cmd.Output(); err == nil {
			return "debian"
		}
	}

	return "Unknown"
}

// detectArchitectureFallback provides fallback architecture detection
func detectArchitectureFallback() string {
	if cmd := exec.Command("uname", "-m"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			return strings.TrimSpace(string(output))
		}
	}
	return "Unknown"
}

// detectHardware detects hardware type
func detectHardware() string {
	// Check for Raspberry Pi
	if cmd := exec.Command("cat", "/proc/cpuinfo"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			cpuinfo := strings.ToLower(string(output))
			if strings.Contains(cpuinfo, "raspberry") || strings.Contains(cpuinfo, "bcm") {
				return "Raspberry Pi"
			}
		}
	}

	// Check for other hardware indicators
	if cmd := exec.Command("dmidecode", "-t", "system"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			dmidecode := strings.ToLower(string(output))
			if strings.Contains(dmidecode, "raspberry") {
				return "Raspberry Pi"
			}
		}
	}

	return "Generic"
}

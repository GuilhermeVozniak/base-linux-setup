package detector

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Pre-compiled regexes for parsing neofetch output
var (
	reOS     = regexp.MustCompile(`^OS:\s*(.+)`)
	reDistro = regexp.MustCompile(`^Distro:\s*(.+)`)
	reKernel = regexp.MustCompile(`^Kernel:\s*(.+)`)
	reHost   = regexp.MustCompile(`^Host:\s*(.+)`)

	reVersion = regexp.MustCompile(`(\d+\.\d+|\d+)`)
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

		// Extract information using pre-compiled regex patterns
		// Note: neofetch uses "OS:" and "Distro:" for distribution info,
		// and "Host:" for hardware. It does not output "Architecture:" directly,
		// so architecture is detected via the uname fallback below.
		if extractFieldRe(line, reOS, &env.OS) ||
			extractFieldRe(line, reDistro, &env.Distribution) ||
			extractFieldRe(line, reKernel, &env.Kernel) ||
			extractFieldRe(line, reHost, &env.Hardware) {
			continue
		}

		// Check for Raspberry Pi indicators
		if containsRaspberryPiIndicator(line) {
			env.IsRaspberryPi = true
		}

		// Extract version info
		if strings.Contains(line, "OS:") || strings.Contains(line, "Distro:") {
			if versionMatch := reVersion.FindString(line); versionMatch != "" {
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

// containsRaspberryPiIndicator checks if a line contains Raspberry Pi indicators
func containsRaspberryPiIndicator(line string) bool {
	lower := strings.ToLower(line)
	return strings.Contains(lower, "raspberry") || strings.Contains(lower, "raspi")
}

// extractFieldRe extracts a field from a line using a pre-compiled regex
func extractFieldRe(line string, re *regexp.Regexp, target *string) bool {
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 && target != nil {
		*target = strings.TrimSpace(matches[1])
		return true
	}
	return false
}

// detectOSFallback provides fallback OS detection
func detectOSFallback() string {
	if output, err := exec.Command("uname", "-s").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	return "Unknown"
}

// detectDistributionFallback provides fallback distribution detection
func detectDistributionFallback() string {
	// Check /etc/os-release (read file directly instead of spawning cat)
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "ID=") {
				return strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
			}
		}
	}

	// Check /etc/debian_version for Debian-based systems
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}

	return "Unknown"
}

// detectArchitectureFallback provides fallback architecture detection
func detectArchitectureFallback() string {
	if output, err := exec.Command("uname", "-m").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	return "Unknown"
}

// detectHardware detects hardware type
func detectHardware() string {
	// Check for Raspberry Pi via /proc/cpuinfo
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		cpuinfo := strings.ToLower(string(data))
		if strings.Contains(cpuinfo, "raspberry") || strings.Contains(cpuinfo, "bcm") {
			return "Raspberry Pi"
		}
	}

	// Check for other hardware indicators
	if output, err := exec.Command("dmidecode", "-t", "system").Output(); err == nil {
		dmidecode := strings.ToLower(string(output))
		if strings.Contains(dmidecode, "raspberry") {
			return "Raspberry Pi"
		}
	}

	return "Generic"
}

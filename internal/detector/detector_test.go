package detector

import (
	"regexp"
	"testing"
)

func TestExtractFieldRe(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		re      *regexp.Regexp
		want    string
		matched bool
	}{
		{
			name:    "OS field",
			line:    "OS: Kali GNU/Linux Rolling x86_64",
			re:      reOS,
			want:    "Kali GNU/Linux Rolling x86_64",
			matched: true,
		},
		{
			name:    "Distro field",
			line:    "Distro: Kali Linux 2023.4",
			re:      reDistro,
			want:    "Kali Linux 2023.4",
			matched: true,
		},
		{
			name:    "Kernel field",
			line:    "Kernel: 6.1.0-kali9-arm64",
			re:      reKernel,
			want:    "6.1.0-kali9-arm64",
			matched: true,
		},
		{
			name:    "Host field",
			line:    "Host: Raspberry Pi 4 Model B Rev 1.4",
			re:      reHost,
			want:    "Raspberry Pi 4 Model B Rev 1.4",
			matched: true,
		},
		{
			name:    "no match",
			line:    "CPU: Cortex-A72 (4) @ 1.800GHz",
			re:      reOS,
			want:    "",
			matched: false,
		},
		{
			name:    "empty line",
			line:    "",
			re:      reOS,
			want:    "",
			matched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			matched := extractFieldRe(tt.line, tt.re, &got)
			if matched != tt.matched {
				t.Errorf("extractFieldRe() matched = %v, want %v", matched, tt.matched)
			}
			if got != tt.want {
				t.Errorf("extractFieldRe() value = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVersionRegex(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{"OS: Kali GNU/Linux Rolling 2023.4", "2023.4"},
		{"Distro: Ubuntu 22.04 LTS", "22.04"},
		{"OS: Arch Linux", ""},
		{"OS: Debian 12", "12"},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := reVersion.FindString(tt.line)
			if got != tt.want {
				t.Errorf("reVersion.FindString(%q) = %q, want %q", tt.line, got, tt.want)
			}
		})
	}
}

func TestRaspberryPiDetection(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"Host: Raspberry Pi 4 Model B", true},
		{"Host: raspi", true},
		{"Host: Dell Latitude", false},
		{"Memory: 3779MiB", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			lower := tt.line
			got := containsRaspberryPiIndicator(lower)
			if got != tt.want {
				t.Errorf("containsRaspberryPiIndicator(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

package main

import (
	_ "embed"
	"fmt"
)

//go:embed scripts/kali-raspberry-pi.json
var KaliRaspberryPiJSON []byte

// GetEmbeddedJSON returns the embedded JSON data for a given filename
func GetEmbeddedJSON(filename string) ([]byte, error) {
	switch filename {
	case "kali-raspberry-pi.json":
		return KaliRaspberryPiJSON, nil
	default:
		return nil, fmt.Errorf("embedded preset file not found: %s", filename)
	}
} 
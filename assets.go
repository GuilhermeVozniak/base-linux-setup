package main

import "embed"

//go:embed scripts/*.json
var PresetFiles embed.FS

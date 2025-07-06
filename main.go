package main

import (
	"fmt"
	"os"

	"base-linux-setup/cmd"
	"base-linux-setup/internal/detector"
	"base-linux-setup/internal/executor"
	"base-linux-setup/internal/presets"
	"base-linux-setup/internal/ui"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Version information - set by build flags
var (
	version   = "dev"
	buildTime = "unknown"
	commit    = "unknown"
)

func main() {
	// Set the embedded JSON getter for presets
	presets.SetEmbeddedJSONGetter(GetEmbeddedJSON)
	
	rootCmd := &cobra.Command{
		Use:     "base-linux-setup",
		Short:   "A CLI tool to setup your local environment based on detected OS",
		Long:    `Base Linux Setup detects your environment and provides customizable presets for system configuration.`,
		Run:     runSetup,
		Version: fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, commit),
	}

	rootCmd.AddCommand(cmd.NewDetectCommand())
	rootCmd.AddCommand(cmd.NewListPresetsCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runSetup(cmd *cobra.Command, args []string) {
	// Print banner
	printBanner()

	// Detect environment
	env, err := detector.DetectEnvironment()
	if err != nil {
		color.Red("Error detecting environment: %v", err)
		os.Exit(1)
	}

	// Display detected environment
	color.Cyan("Detected Environment:")
	color.White("  OS: %s", env.OS)
	color.White("  Distribution: %s", env.Distribution)
	color.White("  Architecture: %s", env.Architecture)
	color.White("  Hardware: %s", env.Hardware)
	fmt.Println()

	// Get preset for environment
	preset := presets.GetPreset(env)
	if preset == nil {
		color.Yellow("No preset found for your environment. Creating a basic preset...")
		preset = presets.GetDefaultPreset()
	}

	// Display preset
	color.Green("Available Preset: %s", preset.Name)
	color.White("Description: %s", preset.Description)
	fmt.Println()

	// Show tasks
	color.Cyan("Preset Tasks:")
	for i, task := range preset.Tasks {
		color.White("  %d. %s", i+1, task.Name)
		if task.Description != "" {
			color.HiBlack("     %s", task.Description)
		}
	}
	fmt.Println()

	// Ask user for customization
	customizedPreset, err := ui.CustomizePreset(preset)
	if err != nil {
		color.Red("Error customizing preset: %v", err)
		os.Exit(1)
	}

	// Confirm execution
	if !ui.ConfirmExecution(customizedPreset) {
		color.Yellow("Setup cancelled by user.")
		os.Exit(0)
	}

	// Execute tasks
	color.Green("Starting setup...")
	fmt.Println()

	executor := executor.NewExecutor()
	for i, task := range customizedPreset.Tasks {
		color.Cyan("Executing task %d/%d: %s", i+1, len(customizedPreset.Tasks), task.Name)

		if err := executor.ExecuteTask(task); err != nil {
			color.Red("Error executing task '%s': %v", task.Name, err)

			if !ui.ContinueOnError() {
				color.Yellow("Setup cancelled.")
				os.Exit(1)
			}
		} else {
			color.Green("✓ Task completed: %s", task.Name)
		}
		fmt.Println()
	}

	color.Green("Setup completed successfully!")
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                    Base Linux Setup                          ║
║              Environment Detection & Setup Tool             ║
╚══════════════════════════════════════════════════════════════╝
`
	color.HiCyan(banner)
}

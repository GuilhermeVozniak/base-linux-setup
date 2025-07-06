package cmd

import (
	"fmt"

	"base-linux-setup/internal/presets"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewListPresetsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list-presets",
		Short: "List all available presets",
		Long:  `List all available presets for different environments.`,
		Run: func(cmd *cobra.Command, args []string) {
			presetList := presets.GetAllPresets()

			color.Cyan("Available Presets:")
			fmt.Println()

			for _, preset := range presetList {
				color.Green("▶ %s", preset.Name)
				color.White("  Environment: %s", preset.Environment)
				color.White("  Description: %s", preset.Description)
				color.HiBlack("  Tasks: %d", len(preset.Tasks))
				
				for i, task := range preset.Tasks {
					color.HiBlack("    %d. %s", i+1, task.Name)
				}
				fmt.Println()
			}
		},
	}
} 
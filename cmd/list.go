package cmd

import (
	"fmt"
	"strings"

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
			presetList, err := presets.GetAllPresets()
			if err != nil {
				color.Red("Error loading presets: %v", err)
				return
			}

			if len(presetList) == 0 {
				color.Yellow("No presets available.")
				return
			}

			color.Cyan("Available Presets:")
			fmt.Println()

			for _, preset := range presetList {
				color.Green("▶ %s", preset.Name)
				color.White("  Environment: %s", preset.Environment)
				color.White("  Description: %s", preset.Description)

				if preset.Match != nil {
					matchInfo := []string{}
					if preset.Match.Distribution != "" {
						matchInfo = append(matchInfo, fmt.Sprintf("distribution=%s", preset.Match.Distribution))
					}
					if preset.Match.OS != "" {
						matchInfo = append(matchInfo, fmt.Sprintf("os=%s", preset.Match.OS))
					}
					if preset.Match.Architecture != "" {
						matchInfo = append(matchInfo, fmt.Sprintf("arch=%s", preset.Match.Architecture))
					}
					if preset.Match.IsRaspberryPi != nil {
						matchInfo = append(matchInfo, fmt.Sprintf("raspberry_pi=%v", *preset.Match.IsRaspberryPi))
					}
					color.White("  Match: [%s]", strings.Join(matchInfo, ", "))
				} else {
					color.White("  Match: [default fallback]")
				}

				color.HiBlack("  Tasks: %d", len(preset.Tasks))

				for i, task := range preset.Tasks {
					color.HiBlack("    %d. %s", i+1, task.Name)
				}
				fmt.Println()
			}
		},
	}
}

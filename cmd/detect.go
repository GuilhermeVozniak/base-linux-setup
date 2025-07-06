package cmd

import (
	"fmt"

	"base-linux-setup/internal/detector"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewDetectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect the current environment",
		Long:  `Detect the current operating system, distribution, architecture, and hardware.`,
		Run: func(cmd *cobra.Command, args []string) {
			env, err := detector.DetectEnvironment()
			if err != nil {
				color.Red("Error detecting environment: %v", err)
				return
			}

			color.Cyan("Environment Information:")
			color.White("  OS: %s", env.OS)
			color.White("  Distribution: %s", env.Distribution)
			color.White("  Version: %s", env.Version)
			color.White("  Architecture: %s", env.Architecture)
			color.White("  Hardware: %s", env.Hardware)
			color.White("  Kernel: %s", env.Kernel)

			if env.IsRaspberryPi {
				color.Green("  🍓 Raspberry Pi detected!")
			}

			fmt.Println()
			color.HiBlack("Raw neofetch output:")
			fmt.Println(env.RawOutput)
		},
	}
} 
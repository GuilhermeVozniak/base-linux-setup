package ui

import (
	"fmt"
	"strconv"
	"strings"

	"base-linux-setup/internal/presets"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// CustomizePreset allows the user to customize a preset
func CustomizePreset(preset *presets.Preset) (*presets.Preset, error) {
	// Ask if user wants to use the preset as-is or customize it
	prompt := promptui.Select{
		Label: "How would you like to proceed?",
		Items: []string{
			"Use preset as-is",
			"Customize tasks (add/remove/modify)",
			"Cancel setup",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	switch result {
	case "Use preset as-is":
		return preset, nil
	case "Cancel setup":
		return nil, fmt.Errorf("setup cancelled by user")
	case "Customize tasks (add/remove/modify)":
		return customizeTasks(preset)
	}

	return preset, nil
}

// customizeTasks allows detailed task customization
func customizeTasks(preset *presets.Preset) (*presets.Preset, error) {
	customizedPreset := &presets.Preset{
		Name:        preset.Name + " (Customized)",
		Environment: preset.Environment,
		Description: preset.Description,
		Tasks:       make([]presets.Task, 0),
	}

	color.Cyan("Customizing tasks for: %s", preset.Name)
	fmt.Println()

	// Go through each task and ask user
	for i, task := range preset.Tasks {
		color.White("Task %d: %s", i+1, task.Name)
		if task.Description != "" {
			color.HiBlack("  Description: %s", task.Description)
		}

		options := []string{"Include this task", "Skip this task"}
		if task.Optional {
			options = append(options, "Mark as optional")
		}

		prompt := promptui.Select{
			Label: "What would you like to do with this task?",
			Items: options,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		switch result {
		case "Include this task":
			customizedPreset.Tasks = append(customizedPreset.Tasks, task)
		case "Skip this task":
			// Skip this task
			continue
		case "Mark as optional":
			task.Optional = true
			customizedPreset.Tasks = append(customizedPreset.Tasks, task)
		}

		fmt.Println()
	}

	// Ask if user wants to add custom tasks
	addCustomPrompt := promptui.Select{
		Label: "Would you like to add custom tasks?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := addCustomPrompt.Run()
	if err == nil && result == "Yes" {
		if err := addCustomTasks(customizedPreset); err != nil {
			return nil, err
		}
	}

	return customizedPreset, nil
}

// addCustomTasks allows adding custom tasks
func addCustomTasks(preset *presets.Preset) error {
	for {
		color.Cyan("Adding custom task...")

		// Get task name
		namePrompt := promptui.Prompt{
			Label: "Task name",
		}
		name, err := namePrompt.Run()
		if err != nil {
			return err
		}

		// Get task description
		descPrompt := promptui.Prompt{
			Label: "Task description (optional)",
		}
		description, _ := descPrompt.Run()

		// Get task commands
		commandPrompt := promptui.Prompt{
			Label: "Commands (separate multiple commands with ';')",
		}
		commandStr, err := commandPrompt.Run()
		if err != nil {
			return err
		}

		commands := strings.Split(commandStr, ";")
		for i, cmd := range commands {
			commands[i] = strings.TrimSpace(cmd)
		}

		// Ask if elevated privileges are needed
		elevatedPrompt := promptui.Select{
			Label: "Does this task require elevated privileges (sudo)?",
			Items: []string{"Yes", "No"},
		}
		_, elevatedResult, _ := elevatedPrompt.Run()
		elevated := elevatedResult == "Yes"

		// Create and add the task
		task := presets.Task{
			Name:        name,
			Description: description,
			Type:        "command",
			Commands:    commands,
			Elevated:    elevated,
		}

		preset.Tasks = append(preset.Tasks, task)

		// Ask if user wants to add more tasks
		morePrompt := promptui.Select{
			Label: "Add another custom task?",
			Items: []string{"Yes", "No"},
		}
		_, result, err := morePrompt.Run()
		if err != nil || result != "Yes" {
			break
		}
	}

	return nil
}

// ConfirmExecution asks user to confirm execution of the preset
func ConfirmExecution(preset *presets.Preset) bool {
	color.Cyan("Final Task List:")
	for i, task := range preset.Tasks {
		status := "✓"
		if task.Optional {
			status = "?"
		}
		color.White("  %s %d. %s", status, i+1, task.Name)
	}
	fmt.Println()

	prompt := promptui.Select{
		Label: fmt.Sprintf("Execute %d tasks?", len(preset.Tasks)),
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false
	}

	return result == "Yes"
}

// ContinueOnError asks user whether to continue when an error occurs
func ContinueOnError() bool {
	color.Yellow("An error occurred during task execution.")

	prompt := promptui.Select{
		Label: "What would you like to do?",
		Items: []string{
			"Continue with remaining tasks",
			"Cancel setup",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false
	}

	return result == "Continue with remaining tasks"
}

// SelectFromList allows user to select items from a list
func SelectFromList(label string, items []string) ([]string, error) {
	selected := make([]string, 0)

	for {
		// Show current selection
		if len(selected) > 0 {
			color.Green("Currently selected: %s", strings.Join(selected, ", "))
		}

		// Create options
		options := make([]string, 0)
		for _, item := range items {
			if !contains(selected, item) {
				options = append(options, item)
			}
		}
		options = append(options, "Done selecting")

		if len(options) == 1 {
			break
		}

		prompt := promptui.Select{
			Label: label,
			Items: options,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if result == "Done selecting" {
			break
		}

		selected = append(selected, result)
	}

	return selected, nil
}

// GetNumberInput gets a number input from user
func GetNumberInput(label string, min, max int) (int, error) {
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s (%d-%d)", label, min, max),
		Validate: func(input string) error {
			num, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid number")
			}
			if num < min || num > max {
				return fmt.Errorf("number must be between %d and %d", min, max)
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(result)
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

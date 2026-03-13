package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"base-linux-setup/internal/presets"

	"github.com/fatih/color"
)

// Executor handles the execution of tasks
type Executor struct {
	dryRun bool
}

// NewExecutor creates a new executor
func NewExecutor() *Executor {
	return &Executor{
		dryRun: false,
	}
}

// NewDryRunExecutor creates a new executor in dry-run mode
func NewDryRunExecutor() *Executor {
	return &Executor{
		dryRun: true,
	}
}

// ExecuteTask executes a single task
func (e *Executor) ExecuteTask(task presets.Task) error {
	// Check task condition before execution
	shouldRun, err := presets.CheckTaskCondition(&task)
	if err != nil {
		color.Yellow("Warning: Failed to check condition for task '%s': %v", task.Name, err)
		color.Yellow("Proceeding with task execution...")
	} else if !shouldRun {
		color.HiBlack("⊗ Skipping task '%s' - condition not met", task.Name)
		return nil
	}

	if e.dryRun {
		return e.dryRunTask(task)
	}

	switch task.Type {
	case "command":
		return e.executeCommands(task)
	case "script":
		return e.executeScript(task)
	case "file":
		return e.createFile(task)
	case "service":
		return e.manageService(task)
	default:
		return fmt.Errorf("unknown task type: %s", task.Type)
	}
}

// ExecutePresetWithDependencies executes a preset with proper dependency ordering
func (e *Executor) ExecutePresetWithDependencies(preset *presets.Preset) error {
	// Sort tasks by dependencies
	sortedTasks, err := presets.SortTasksByDependencies(preset.Tasks)
	if err != nil {
		return fmt.Errorf("failed to resolve task dependencies: %v", err)
	}

	// Execute tasks in dependency order
	for i, task := range sortedTasks {
		color.Cyan("Executing task %d/%d: %s", i+1, len(sortedTasks), task.Name)

		if err := e.ExecuteTask(task); err != nil {
			return fmt.Errorf("error executing task '%s': %v", task.Name, err)
		}

		color.Green("✓ Task completed: %s", task.Name)
		fmt.Println()
	}

	return nil
}

// executeCommands executes a list of commands
func (e *Executor) executeCommands(task presets.Task) error {
	for i, command := range task.Commands {
		if len(task.Commands) > 1 {
			color.HiBlack("  Command %d/%d: %s", i+1, len(task.Commands), command)
		}

		if err := e.runCommand(command); err != nil {
			return fmt.Errorf("command failed: %s - %v", command, err)
		}
	}
	return nil
}

// executeScript executes a script
func (e *Executor) executeScript(task presets.Task) error {
	// Create temporary script file
	tmpFile, err := os.CreateTemp("", "setup-script-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temp script file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write script content
	if _, err := tmpFile.WriteString(task.Script); err != nil {
		return fmt.Errorf("failed to write script: %v", err)
	}
	tmpFile.Close()

	// Make script executable
	if err := os.Chmod(tmpFile.Name(), 0o755); err != nil {
		return fmt.Errorf("failed to make script executable: %v", err)
	}

	// Execute script
	return e.runCommand(tmpFile.Name())
}

// createFile creates a file with specified content
func (e *Executor) createFile(task presets.Task) error {
	// File tasks expect Commands[0] to be the file path and Script to be the content
	if len(task.Commands) == 0 {
		return fmt.Errorf("file task requires file path in Commands[0]")
	}

	filePath := task.Commands[0]
	content := task.Script

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	// Create or overwrite the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	// Write content to file
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content to file %s: %v", filePath, err)
	}

	// Set file permissions if specified in Commands[1]
	if len(task.Commands) > 1 {
		if perm, err := strconv.ParseUint(task.Commands[1], 8, 32); err == nil {
			if err := os.Chmod(filePath, os.FileMode(perm)); err != nil {
				color.Yellow("Warning: Failed to set permissions on %s: %v", filePath, err)
			}
		}
	}

	color.HiGreen("    ✓ File created: %s", filePath)
	return nil
}

// manageService manages system services
func (e *Executor) manageService(task presets.Task) error {
	// Service tasks expect Commands to contain systemctl operations
	// Format: ["service_name", "action"] where action is start/stop/enable/disable/restart
	if len(task.Commands) < 2 {
		return fmt.Errorf("service task requires service name and action in Commands")
	}

	serviceName := task.Commands[0]
	action := task.Commands[1]

	// Validate action
	validActions := []string{"start", "stop", "enable", "disable", "restart", "reload", "status"}
	isValidAction := false
	for _, validAction := range validActions {
		if action == validAction {
			isValidAction = true
			break
		}
	}

	if !isValidAction {
		return fmt.Errorf("invalid service action: %s. Valid actions: %v", action, validActions)
	}

	// Build systemctl command
	var cmd *exec.Cmd
	if action == "status" {
		// For status, we don't need sudo and we want to capture output
		cmd = exec.Command("systemctl", action, serviceName)
	} else {
		// For other actions, we typically need sudo
		cmd = exec.Command("sudo", "systemctl", action, serviceName)
	}

	// Set up command execution
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	color.HiBlack("    Running: systemctl %s %s", action, serviceName)

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	if err != nil {
		color.Red("    ✗ Service operation failed in %v", duration)
		return fmt.Errorf("systemctl %s %s failed: %v", action, serviceName, err)
	}

	color.HiGreen("    ✓ Service operation completed in %v", duration)

	// Additional actions based on the operation
	switch action {
	case "enable":
		color.HiBlack("    Service %s will start automatically on boot", serviceName)
	case "disable":
		color.HiBlack("    Service %s will not start automatically on boot", serviceName)
	case "start":
		color.HiBlack("    Service %s is now running", serviceName)
	case "stop":
		color.HiBlack("    Service %s is now stopped", serviceName)
	}

	return nil
}

// runCommand runs a single command
func (e *Executor) runCommand(command string) error {
	if strings.TrimSpace(command) == "" {
		return fmt.Errorf("empty command")
	}

	// Use sh -c to properly handle quoted arguments, pipes, and shell syntax
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Set environment variables
	cmd.Env = os.Environ()

	// Run command
	color.HiBlack("    Running: %s", command)

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	if err != nil {
		color.Red("    ✗ Failed in %v", duration)
		return err
	}

	color.HiGreen("    ✓ Completed in %v", duration)
	return nil
}

// dryRunTask simulates task execution without actually running commands
func (e *Executor) dryRunTask(task presets.Task) error {
	color.Yellow("[DRY RUN] Would execute task: %s", task.Name)

	switch task.Type {
	case "command":
		for _, command := range task.Commands {
			color.HiBlack("  [DRY RUN] Command: %s", command)
		}
	case "script":
		color.HiBlack("  [DRY RUN] Script execution")
		// Show first few lines of script
		lines := strings.Split(task.Script, "\n")
		for i, line := range lines {
			if i >= 3 {
				color.HiBlack("  [DRY RUN] ... (%d more lines)", len(lines)-i)
				break
			}
			if strings.TrimSpace(line) != "" {
				color.HiBlack("  [DRY RUN] %s", line)
			}
		}
	case "file":
		color.HiBlack("  [DRY RUN] File creation")
	case "service":
		color.HiBlack("  [DRY RUN] Service management")
	}

	return nil
}

// SetDryRun sets the dry-run mode
func (e *Executor) SetDryRun(dryRun bool) {
	e.dryRun = dryRun
}

// IsDryRun returns whether executor is in dry-run mode
func (e *Executor) IsDryRun() bool {
	return e.dryRun
}

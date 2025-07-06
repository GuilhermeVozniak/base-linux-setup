package executor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	tmpFile, err := ioutil.TempFile("", "setup-script-*.sh")
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
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to make script executable: %v", err)
	}

	// Execute script
	return e.runCommand(tmpFile.Name())
}

// createFile creates a file with specified content
func (e *Executor) createFile(task presets.Task) error {
	// Extract file path and content from task
	// This would need to be implemented based on how file tasks are structured
	return fmt.Errorf("file task type not fully implemented")
}

// manageService manages system services
func (e *Executor) manageService(task presets.Task) error {
	// This would handle systemctl commands
	return fmt.Errorf("service task type not fully implemented")
}

// runCommand runs a single command
func (e *Executor) runCommand(command string) error {
	// Parse command
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// Create command
	cmd := exec.Command(parts[0], parts[1:]...)
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

// ValidatePrerequisites checks if prerequisites are met for task execution
func (e *Executor) ValidatePrerequisites() error {
	// Check if running as root when needed
	if os.Geteuid() == 0 {
		color.Yellow("Warning: Running as root user")
	}

	// Check if we have network connectivity
	if err := e.checkNetworkConnectivity(); err != nil {
		return fmt.Errorf("network connectivity check failed: %v", err)
	}

	// Check available disk space
	if err := e.checkDiskSpace(); err != nil {
		return fmt.Errorf("disk space check failed: %v", err)
	}

	return nil
}

// checkNetworkConnectivity checks internet connectivity
func (e *Executor) checkNetworkConnectivity() error {
	cmd := exec.Command("ping", "-c", "1", "8.8.8.8")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("no internet connectivity")
	}
	return nil
}

// checkDiskSpace checks available disk space
func (e *Executor) checkDiskSpace() error {
	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check disk space: %v", err)
	}

	color.HiBlack("Disk space: %s", strings.TrimSpace(string(output)))
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

// CreateBackup creates a backup of important files before making changes
func (e *Executor) CreateBackup() error {
	backupDir := filepath.Join(os.Getenv("HOME"), ".config", "base-linux-setup", "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Backup important files
	filesToBackup := []string{
		"/etc/fstab",
		"/boot/config.txt",
		"/etc/modules",
		filepath.Join(os.Getenv("HOME"), ".bashrc"),
		filepath.Join(os.Getenv("HOME"), ".profile"),
	}

	timestamp := time.Now().Format("20060102-150405")
	
	for _, file := range filesToBackup {
		if _, err := os.Stat(file); err == nil {
			backupFile := filepath.Join(backupDir, fmt.Sprintf("%s.%s", filepath.Base(file), timestamp))
			if err := e.copyFile(file, backupFile); err != nil {
				color.Yellow("Warning: Failed to backup %s: %v", file, err)
			} else {
				color.HiBlack("Backed up: %s -> %s", file, backupFile)
			}
		}
	}

	return nil
}

// copyFile copies a file from src to dst
func (e *Executor) copyFile(src, dst string) error {
	cmd := exec.Command("cp", src, dst)
	return cmd.Run()
}

// RestoreBackup restores files from backup
func (e *Executor) RestoreBackup(timestamp string) error {
	backupDir := filepath.Join(os.Getenv("HOME"), ".config", "base-linux-setup", "backups")
	
	// List available backups
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %v", err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), timestamp) {
			color.Cyan("Found backup: %s", file.Name())
			// Restore logic would go here
		}
	}

	return nil
} 
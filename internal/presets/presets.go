package presets

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"base-linux-setup/internal/detector"
)

// embeddedFS holds the embedded preset files
var embeddedFS fs.FS

// SetEmbeddedFS sets the embedded filesystem containing preset JSON files
func SetEmbeddedFS(fsys fs.FS) {
	embeddedFS = fsys
}

// Task represents a single setup task
type Task struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"` // "command", "script", "file", "service"
	Commands    []string `json:"commands"`
	Script      string   `json:"script"`
	Elevated    bool     `json:"elevated"` // requires sudo
	Optional    bool     `json:"optional"`
	Condition   string   `json:"condition,omitempty"`  // Shell command to check if task should run
	DependsOn   []string `json:"depends_on,omitempty"` // Array of task names this task depends on
}

// MatchCriteria defines when a preset should be auto-selected based on the detected environment
type MatchCriteria struct {
	Distribution  string `json:"distribution,omitempty"`
	OS            string `json:"os,omitempty"`
	Architecture  string `json:"architecture,omitempty"`
	IsRaspberryPi *bool  `json:"is_raspberry_pi,omitempty"`
}

// Preset represents a collection of tasks for a specific environment
type Preset struct {
	Name        string         `json:"name"`
	Environment string         `json:"environment"`
	Description string         `json:"description"`
	Match       *MatchCriteria `json:"match,omitempty"`
	Tasks       []Task         `json:"tasks"`
}

// LoadExternalPreset loads a preset from an external JSON file
func LoadExternalPreset(filePath string) (*Preset, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", filePath)
	}

	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	// Parse JSON
	var preset Preset
	if err := json.Unmarshal(data, &preset); err != nil {
		return nil, fmt.Errorf("failed to parse configuration JSON: %v", err)
	}

	// Validate preset
	if err := validatePreset(&preset); err != nil {
		return nil, fmt.Errorf("invalid preset configuration: %v", err)
	}

	return &preset, nil
}

// validatePreset validates the preset configuration
func validatePreset(preset *Preset) error {
	if preset.Name == "" {
		return fmt.Errorf("preset name is required")
	}

	if len(preset.Tasks) == 0 {
		return fmt.Errorf("preset must have at least one task")
	}

	// Validate tasks and dependencies
	taskNames := make(map[string]bool)
	for _, task := range preset.Tasks {
		if task.Name == "" {
			return fmt.Errorf("task name is required")
		}
		taskNames[task.Name] = true
	}

	// Validate dependencies
	for _, task := range preset.Tasks {
		for _, dep := range task.DependsOn {
			if !taskNames[dep] {
				return fmt.Errorf("task '%s' depends on non-existent task '%s'", task.Name, dep)
			}
		}
	}

	return nil
}

// CheckTaskCondition checks if a task's condition is met
func CheckTaskCondition(task *Task) (bool, error) {
	if task.Condition == "" {
		return true, nil // No condition means always run
	}

	cmd := exec.Command("sh", "-c", task.Condition)
	err := cmd.Run()

	// If command exits with status 0, condition is met
	return err == nil, nil
}

// SortTasksByDependencies sorts tasks to ensure dependencies are executed first
func SortTasksByDependencies(tasks []Task) ([]Task, error) {
	if len(tasks) == 0 {
		return tasks, nil
	}

	// Create a map of task names to tasks
	taskMap := make(map[string]*Task)
	for i := range tasks {
		taskMap[tasks[i].Name] = &tasks[i]
	}

	// Track visited and recursion stack for cycle detection
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	sorted := make([]Task, 0, len(tasks))

	// DFS function for topological sort
	var dfs func(taskName string) error
	dfs = func(taskName string) error {
		if recStack[taskName] {
			return fmt.Errorf("circular dependency detected involving task '%s'", taskName)
		}
		if visited[taskName] {
			return nil
		}

		visited[taskName] = true
		recStack[taskName] = true

		task := taskMap[taskName]
		for _, dep := range task.DependsOn {
			if err := dfs(dep); err != nil {
				return err
			}
		}

		recStack[taskName] = false
		sorted = append(sorted, *task)
		return nil
	}

	// Sort all tasks
	for taskName := range taskMap {
		if !visited[taskName] {
			if err := dfs(taskName); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}

// loadAllEmbeddedPresets loads all preset JSON files from the embedded filesystem
func loadAllEmbeddedPresets() ([]*Preset, error) {
	if embeddedFS == nil {
		return nil, fmt.Errorf("embedded filesystem not set")
	}

	entries, err := fs.ReadDir(embeddedFS, "scripts")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded scripts directory: %v", err)
	}

	var presets []*Preset
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := fs.ReadFile(embeddedFS, "scripts/"+entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded preset %s: %v", entry.Name(), err)
		}

		var preset Preset
		if err := json.Unmarshal(data, &preset); err != nil {
			return nil, fmt.Errorf("failed to parse embedded preset %s: %v", entry.Name(), err)
		}

		if err := validatePreset(&preset); err != nil {
			return nil, fmt.Errorf("invalid embedded preset %s: %v", entry.Name(), err)
		}

		presets = append(presets, &preset)
	}

	return presets, nil
}

// matchesEnvironment checks if a preset's match criteria fits the detected environment
func matchesEnvironment(match *MatchCriteria, env *detector.Environment) bool {
	if match == nil {
		return false
	}

	if match.Distribution != "" {
		if !strings.Contains(strings.ToLower(env.Distribution), strings.ToLower(match.Distribution)) &&
			!strings.Contains(strings.ToLower(env.OS), strings.ToLower(match.Distribution)) {
			return false
		}
	}

	if match.OS != "" {
		if !strings.Contains(strings.ToLower(env.OS), strings.ToLower(match.OS)) {
			return false
		}
	}

	if match.Architecture != "" {
		if !strings.Contains(strings.ToLower(env.Architecture), strings.ToLower(match.Architecture)) {
			return false
		}
	}

	if match.IsRaspberryPi != nil {
		if *match.IsRaspberryPi != env.IsRaspberryPi {
			return false
		}
	}

	return true
}

// matchSpecificity returns the number of non-empty match fields.
// Higher specificity means a more specific (better) match.
func matchSpecificity(match *MatchCriteria) int {
	if match == nil {
		return 0
	}
	score := 0
	if match.Distribution != "" {
		score++
	}
	if match.OS != "" {
		score++
	}
	if match.Architecture != "" {
		score++
	}
	if match.IsRaspberryPi != nil {
		score++
	}
	return score
}

// GetPreset returns the best matching preset for the given environment.
// It loads all embedded presets, matches them against the environment,
// and returns the most specific match. Falls back to the default preset
// (one without a match field) if no match is found.
func GetPreset(env *detector.Environment) *Preset {
	allPresets, err := loadAllEmbeddedPresets()
	if err != nil {
		return nil
	}

	var bestMatch *Preset
	bestScore := 0
	var defaultPreset *Preset

	for _, preset := range allPresets {
		if preset.Match == nil {
			defaultPreset = preset
			continue
		}

		if matchesEnvironment(preset.Match, env) {
			score := matchSpecificity(preset.Match)
			if score > bestScore {
				bestScore = score
				bestMatch = preset
			}
		}
	}

	if bestMatch != nil {
		return bestMatch
	}
	return defaultPreset
}

// GetAllPresets returns all available embedded presets
func GetAllPresets() []*Preset {
	allPresets, err := loadAllEmbeddedPresets()
	if err != nil {
		return nil
	}
	return allPresets
}

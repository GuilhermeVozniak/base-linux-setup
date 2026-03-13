package presets

import (
	"encoding/json"
	"io/fs"
	"testing"
	"testing/fstest"

	"base-linux-setup/internal/detector"
)

func TestValidatePreset(t *testing.T) {
	tests := []struct {
		name    string
		preset  Preset
		wantErr bool
	}{
		{
			name: "valid preset",
			preset: Preset{
				Name: "test",
				Tasks: []Task{
					{Name: "task1", Type: "command", Commands: []string{"echo hi"}},
				},
			},
			wantErr: false,
		},
		{
			name:    "missing name",
			preset:  Preset{Tasks: []Task{{Name: "t"}}},
			wantErr: true,
		},
		{
			name:    "no tasks",
			preset:  Preset{Name: "test"},
			wantErr: true,
		},
		{
			name: "missing task name",
			preset: Preset{
				Name:  "test",
				Tasks: []Task{{Type: "command"}},
			},
			wantErr: true,
		},
		{
			name: "valid dependency",
			preset: Preset{
				Name: "test",
				Tasks: []Task{
					{Name: "a", Type: "command"},
					{Name: "b", Type: "command", DependsOn: []string{"a"}},
				},
			},
			wantErr: false,
		},
		{
			name: "missing dependency",
			preset: Preset{
				Name: "test",
				Tasks: []Task{
					{Name: "a", Type: "command", DependsOn: []string{"nonexistent"}},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePreset(&tt.preset)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePreset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMatchesEnvironment(t *testing.T) {
	boolTrue := true
	boolFalse := false

	env := &detector.Environment{
		OS:            "Linux",
		Distribution:  "Kali",
		Architecture:  "aarch64",
		IsRaspberryPi: true,
	}

	tests := []struct {
		name  string
		match *MatchCriteria
		want  bool
	}{
		{
			name:  "nil match",
			match: nil,
			want:  false,
		},
		{
			name:  "empty match (matches everything)",
			match: &MatchCriteria{},
			want:  true,
		},
		{
			name:  "distribution match",
			match: &MatchCriteria{Distribution: "kali"},
			want:  true,
		},
		{
			name:  "distribution mismatch",
			match: &MatchCriteria{Distribution: "ubuntu"},
			want:  false,
		},
		{
			name:  "os match",
			match: &MatchCriteria{OS: "linux"},
			want:  true,
		},
		{
			name:  "architecture match",
			match: &MatchCriteria{Architecture: "aarch64"},
			want:  true,
		},
		{
			name:  "architecture mismatch",
			match: &MatchCriteria{Architecture: "x86_64"},
			want:  false,
		},
		{
			name:  "raspberry pi match true",
			match: &MatchCriteria{IsRaspberryPi: &boolTrue},
			want:  true,
		},
		{
			name:  "raspberry pi match false",
			match: &MatchCriteria{IsRaspberryPi: &boolFalse},
			want:  false,
		},
		{
			name: "multiple criteria all match",
			match: &MatchCriteria{
				Distribution:  "kali",
				Architecture:  "aarch64",
				IsRaspberryPi: &boolTrue,
			},
			want: true,
		},
		{
			name: "multiple criteria one fails",
			match: &MatchCriteria{
				Distribution: "kali",
				Architecture: "x86_64",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesEnvironment(tt.match, env)
			if got != tt.want {
				t.Errorf("matchesEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchSpecificity(t *testing.T) {
	boolTrue := true

	tests := []struct {
		name  string
		match *MatchCriteria
		want  int
	}{
		{"nil", nil, 0},
		{"empty", &MatchCriteria{}, 0},
		{"one field", &MatchCriteria{Distribution: "kali"}, 1},
		{"two fields", &MatchCriteria{Distribution: "kali", OS: "linux"}, 2},
		{"all fields", &MatchCriteria{Distribution: "kali", OS: "linux", Architecture: "arm", IsRaspberryPi: &boolTrue}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchSpecificity(tt.match)
			if got != tt.want {
				t.Errorf("matchSpecificity() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSortTasksByDependencies(t *testing.T) {
	t.Run("empty tasks", func(t *testing.T) {
		sorted, err := SortTasksByDependencies(nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(sorted) != 0 {
			t.Errorf("expected 0 tasks, got %d", len(sorted))
		}
	})

	t.Run("no dependencies", func(t *testing.T) {
		tasks := []Task{
			{Name: "a"},
			{Name: "b"},
		}
		sorted, err := SortTasksByDependencies(tasks)
		if err != nil {
			t.Fatal(err)
		}
		if len(sorted) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(sorted))
		}
	})

	t.Run("linear dependency chain", func(t *testing.T) {
		tasks := []Task{
			{Name: "c", DependsOn: []string{"b"}},
			{Name: "a"},
			{Name: "b", DependsOn: []string{"a"}},
		}
		sorted, err := SortTasksByDependencies(tasks)
		if err != nil {
			t.Fatal(err)
		}

		// Build position map
		pos := make(map[string]int)
		for i, task := range sorted {
			pos[task.Name] = i
		}

		if pos["a"] >= pos["b"] {
			t.Errorf("a should come before b, got a=%d b=%d", pos["a"], pos["b"])
		}
		if pos["b"] >= pos["c"] {
			t.Errorf("b should come before c, got b=%d c=%d", pos["b"], pos["c"])
		}
	})

	t.Run("circular dependency", func(t *testing.T) {
		tasks := []Task{
			{Name: "a", DependsOn: []string{"b"}},
			{Name: "b", DependsOn: []string{"a"}},
		}
		_, err := SortTasksByDependencies(tasks)
		if err == nil {
			t.Error("expected error for circular dependency")
		}
	})
}

// makeTestFS creates a test filesystem with preset JSON files
func makeTestFS(presets map[string]Preset) fs.FS {
	mapFS := make(fstest.MapFS)
	for name, preset := range presets {
		data, _ := json.Marshal(preset)
		mapFS["scripts/"+name+".json"] = &fstest.MapFile{Data: data}
	}
	return mapFS
}

func TestGetPreset(t *testing.T) {
	boolTrue := true

	testFS := makeTestFS(map[string]Preset{
		"default": {
			Name:        "Default",
			Environment: "Generic",
			Tasks:       []Task{{Name: "default-task"}},
		},
		"kali": {
			Name:        "Kali Setup",
			Environment: "Kali Linux",
			Match:       &MatchCriteria{Distribution: "kali"},
			Tasks:       []Task{{Name: "kali-task"}},
		},
		"kali-pi": {
			Name:        "Kali Pi Setup",
			Environment: "Kali on Pi",
			Match:       &MatchCriteria{Distribution: "kali", IsRaspberryPi: &boolTrue},
			Tasks:       []Task{{Name: "kali-pi-task"}},
		},
	})

	// Save and restore global state
	oldFS := embeddedFS
	defer func() { embeddedFS = oldFS }()
	embeddedFS = testFS

	t.Run("exact match returns most specific", func(t *testing.T) {
		env := &detector.Environment{
			Distribution:  "Kali",
			OS:            "Linux",
			IsRaspberryPi: true,
		}
		preset, err := GetPreset(env)
		if err != nil {
			t.Fatal(err)
		}
		if preset == nil {
			t.Fatal("expected a preset, got nil")
		}
		if preset.Name != "Kali Pi Setup" {
			t.Errorf("expected 'Kali Pi Setup', got '%s'", preset.Name)
		}
	})

	t.Run("partial match returns less specific", func(t *testing.T) {
		env := &detector.Environment{
			Distribution:  "Kali",
			OS:            "Linux",
			IsRaspberryPi: false,
		}
		preset, err := GetPreset(env)
		if err != nil {
			t.Fatal(err)
		}
		if preset == nil {
			t.Fatal("expected a preset, got nil")
		}
		if preset.Name != "Kali Setup" {
			t.Errorf("expected 'Kali Setup', got '%s'", preset.Name)
		}
	})

	t.Run("no match returns default", func(t *testing.T) {
		env := &detector.Environment{
			Distribution: "Fedora",
			OS:           "Linux",
		}
		preset, err := GetPreset(env)
		if err != nil {
			t.Fatal(err)
		}
		if preset == nil {
			t.Fatal("expected default preset, got nil")
		}
		if preset.Name != "Default" {
			t.Errorf("expected 'Default', got '%s'", preset.Name)
		}
	})

	t.Run("nil embedded FS returns error", func(t *testing.T) {
		embeddedFS = nil
		_, err := GetPreset(&detector.Environment{})
		if err == nil {
			t.Error("expected error when embeddedFS is nil")
		}
	})
}

func TestGetAllPresets(t *testing.T) {
	testFS := makeTestFS(map[string]Preset{
		"one": {Name: "One", Tasks: []Task{{Name: "t1"}}},
		"two": {Name: "Two", Tasks: []Task{{Name: "t2"}}},
	})

	oldFS := embeddedFS
	defer func() { embeddedFS = oldFS }()
	embeddedFS = testFS

	presets, err := GetAllPresets()
	if err != nil {
		t.Fatal(err)
	}
	if len(presets) != 2 {
		t.Errorf("expected 2 presets, got %d", len(presets))
	}
}

func TestLoadExternalPreset(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		_, err := LoadExternalPreset("/nonexistent/file.json")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})
}

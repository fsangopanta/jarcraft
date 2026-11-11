package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"jarcraft/internal/config"
)

func TestLoadUsesDefaultConfiguration(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	cfg, err := config.Load("")

	if err != nil {
		t.Fatal(err)
	}

	if cfg.Build.Output != "./builds" {
		t.Errorf("expected ./builds, got %s", cfg.Build.Output)
	}

	if cfg.Build.Parallel {
		t.Error("expected parallel to be false")
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("expected info, got %s", cfg.Logging.Level)
	}

	if len(cfg.Scanner.Exclude) == 0 {
		t.Error("expected default excluded directories")
	}
}

func TestLoadCustomConfiguration(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "jarcraft.yaml")

	content := []byte(`
build:
  output: "./artifacts"
  parallel: true
  fail_fast: true

scanner:
  exclude:
    - ".git"
    - "target"

logging:
  level: "debug"
`)

	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(configPath)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.Build.Output != "./artifacts" {
		t.Errorf("expected ./artifacts, got %s", cfg.Build.Output)
	}

	if !cfg.Build.Parallel {
		t.Error("expected parallel to be true")
	}

	if !cfg.Build.FailFast {
		t.Error("expected fail_fast to be true")
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("expected debug, got %s", cfg.Logging.Level)
	}

	if len(cfg.Scanner.Exclude) != 2 {
		t.Errorf("expected 2 excluded directories, got %d", len(cfg.Scanner.Exclude))
	}
}

func TestLoadRejectsInvalidLoggingLevel(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "jarcraft.yaml")

	content := []byte(`
logging:
  level: "verbose"
`)

	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := config.Load(configPath)

	if err == nil {
		t.Fatal("expected error for invalid logging level")
	}
}

func TestLoadRejectsMalformedYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "jarcraft.yaml")

	content := []byte(`
build:
  output: [this is not valid yaml
`)

	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := config.Load(configPath)

	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestResolveProjectConfigUsesDefaults(t *testing.T) {
	skipTests := true

	cfg := &config.Config{
		Defaults: config.ProjectConfig{
			SkipTests:   &skipTests,
			JavaVersion: "21",
		},
		Projects: map[string]config.ProjectConfig{},
	}

	result := cfg.ResolveProjectConfig("orders")

	if !result.SkipTests {
		t.Error("expected SkipTests to be true")
	}

	if result.JavaVersion != "21" {
		t.Errorf("expected JavaVersion 21, got %s", result.JavaVersion)
	}
}

func TestResolveProjectConfigOverridesDefaults(t *testing.T) {
	defaultSkipTests := true
	projectSkipTests := false

	cfg := &config.Config{
		Defaults: config.ProjectConfig{
			SkipTests:   &defaultSkipTests,
			JavaVersion: "21",
		},
		Projects: map[string]config.ProjectConfig{
			"orders": {
				SkipTests:   &projectSkipTests,
				JavaVersion: "17",
			},
		},
	}

	result := cfg.ResolveProjectConfig("orders")

	if result.SkipTests {
		t.Error("expected project to execute tests")
	}

	if result.JavaVersion != "17" {
		t.Errorf("expected JavaVersion 17, got %s", result.JavaVersion)
	}
}

func TestResolveProjectConfigUsesDefaultWhenProjectDoesNotExist(t *testing.T) {
	skipTests := false

	cfg := &config.Config{
		Defaults: config.ProjectConfig{
			SkipTests:   &skipTests,
			JavaVersion: "25",
		},
		Projects: map[string]config.ProjectConfig{},
	}

	result := cfg.ResolveProjectConfig("unknown-project")

	if result.SkipTests {
		t.Error("expected SkipTests to be false")
	}

	if result.JavaVersion != "25" {
		t.Errorf("expected JavaVersion 25, got %s", result.JavaVersion)
	}
}

func TestResolveProjectConfigWithOverrides(t *testing.T) {
	defaultSkipTests := false
	projectSkipTests := false
	cliSkipTests := true

	cfg := &config.Config{
		Defaults: config.ProjectConfig{
			SkipTests:   &defaultSkipTests,
			JavaVersion: "21",
		},
		Projects: map[string]config.ProjectConfig{
			"orders": {
				SkipTests:   &projectSkipTests,
				JavaVersion: "17",
			},
		},
	}

	result := cfg.ResolveProjectConfigWithOverrides(
		"orders",
		config.BuildOverrides{
			SkipTests:   &cliSkipTests,
			JavaVersion: "25",
		},
	)

	if !result.SkipTests {
		t.Error("expected CLI to override SkipTests to true")
	}

	if result.JavaVersion != "25" {
		t.Errorf("expected JavaVersion 25, got %s", result.JavaVersion)
	}
}

func TestResolveProjectConfigWithEmptyOverrides(t *testing.T) {
	projectSkipTests := false

	cfg := &config.Config{
		Projects: map[string]config.ProjectConfig{
			"orders": {
				SkipTests:   &projectSkipTests,
				JavaVersion: "17",
			},
		},
	}

	result := cfg.ResolveProjectConfigWithOverrides(
		"orders",
		config.BuildOverrides{},
	)

	if result.SkipTests {
		t.Error("expected project SkipTests configuration to remain false")
	}

	if result.JavaVersion != "17" {
		t.Errorf("expected JavaVersion 17, got %s", result.JavaVersion)
	}
}

package engine_test

import (
	"context"
	"errors"
	"jarcraft/internal/config"
	"os"
	"path/filepath"
	"testing"

	"jarcraft/internal/engine"
	"jarcraft/internal/scanner"
)

func createMavenProject(t *testing.T, root, name string) {
	t.Helper()

	projectPath := filepath.Join(root, name)

	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(
		filepath.Join(projectPath, "pom.xml"),
		[]byte("<project></project>"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEngineBuildsAllProjects(t *testing.T) {
	root := t.TempDir()

	createMavenProject(t, root, "payments")
	createMavenProject(t, root, "customers")

	buildCount := 0

	e := engine.NewWithBuilder(
		func(ctx context.Context, project scanner.Project) error {
			buildCount++
			return nil
		},
	)

	results, err := e.Run(context.Background(), root)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if buildCount != 2 {
		t.Errorf("expected 2 builds, got %d", buildCount)
	}

	for _, result := range results {
		if !result.Success {
			t.Errorf("expected project %s to succeed", result.Project.Name)
		}
	}
}

func TestEngineContinuesAfterBuildFailure(t *testing.T) {
	root := t.TempDir()

	createMavenProject(t, root, "payments")
	createMavenProject(t, root, "customers")

	buildCount := 0

	e := engine.NewWithBuilder(
		func(ctx context.Context, project scanner.Project) error {
			buildCount++

			if project.Name == "payments" {
				return errors.New("compilation failed")
			}

			return nil
		},
	)

	results, err := e.Run(context.Background(), root)

	if err != nil {
		t.Fatal(err)
	}

	if buildCount != 2 {
		t.Errorf("expected 2 builds, got %d", buildCount)
	}

	var failures int

	for _, result := range results {
		if result.Error != nil {
			failures++
		}
	}

	if failures != 1 {
		t.Errorf("expected 1 failure, got %d", failures)
	}
}

func TestEngineReturnsErrorForInvalidRoot(t *testing.T) {
	e := engine.NewWithBuilder(
		func(ctx context.Context, project scanner.Project) error {
			return nil
		},
	)

	_, err := e.Run(context.Background(), filepath.Join(t.TempDir(), "missing"))

	if err == nil {
		t.Fatal("expected error for invalid root")
	}
}

func TestRunWithConfigResolvesProjectConfiguration(t *testing.T) {
	root := t.TempDir()

	projectPath := filepath.Join(root, "orders")

	if err := os.Mkdir(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		filepath.Join(projectPath, "build.gradle"),
		[]byte(""),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	skipTests := true

	cfg := &config.Config{
		Scanner: config.ScannerConfig{
			Exclude: []string{},
		},
		Projects: map[string]config.ProjectConfig{
			"orders": {
				SkipTests:   &skipTests,
				JavaVersion: "21",
			},
		},
	}

	var receivedConfig config.ResolvedProjectConfig

	engine := engine.NewWithConfiguredBuilder(
		func(
			ctx context.Context,
			project scanner.Project,
			options config.ResolvedProjectConfig,
		) error {
			receivedConfig = options
			return nil
		},
	)

	results, err := engine.RunWithConfig(
		context.Background(),
		root,
		cfg,
		config.BuildOverrides{},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if !receivedConfig.SkipTests {
		t.Error("expected SkipTests to be true")
	}

	if receivedConfig.JavaVersion != "21" {
		t.Errorf(
			"expected JavaVersion 21, got %s",
			receivedConfig.JavaVersion,
		)
	}
}

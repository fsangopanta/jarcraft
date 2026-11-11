package scanner__test

import (
	"os"
	"path/filepath"
	"testing"

	"jarcraft/internal/scanner"
)

func TestScanFindsNestedProjects(t *testing.T) {
	root := t.TempDir()

	createProject(t, root, "banking/payment-service", "pom.xml")
	createProject(t, root, "banking/notification-service", "build.gradle")
	createProject(t, root, "ecommerce/backend", "pom.xml")
	createProject(t, root, "personal/kafka-provider", "build.gradle.kts")

	projects, err := scanner.Scan(root)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projects) != 4 {
		t.Fatalf("expected 4 projects, got %d", len(projects))
	}

	expected := map[string]scanner.BuildTool{
		"payment-service":      scanner.Maven,
		"notification-service": scanner.Gradle,
		"backend":              scanner.Maven,
		"kafka-provider":       scanner.Gradle,
	}

	for _, project := range projects {
		buildTool, exists := expected[project.Name]

		if !exists {
			t.Errorf("unexpected project: %s", project.Name)
			continue
		}

		if project.BuildTool != buildTool {
			t.Errorf(
				"project %s: expected %s, got %s",
				project.Name,
				buildTool,
				project.BuildTool,
			)
		}
	}
}

func TestScanIgnoresExcludedDirectories(t *testing.T) {
	root := t.TempDir()

	createProject(t, root, "application", "pom.xml")
	createProject(t, root, ".git/hidden-project", "pom.xml")
	createProject(t, root, ".idea/hidden-project", "build.gradle")

	projects, err := scanner.Scan(root)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	if projects[0].Name != "application" {
		t.Errorf("expected application, got %s", projects[0].Name)
	}
}

func TestScanRejectsFileAsRoot(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "test.txt")

	err := os.WriteFile(file, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = scanner.Scan(file)

	if err == nil {
		t.Fatal("expected error when scanning a file")
	}
}

func createProject(
	t *testing.T,
	root string,
	relativePath string,
	buildFile string,
) {
	t.Helper()

	projectPath := filepath.Join(root, relativePath)

	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	filePath := filepath.Join(projectPath, buildFile)

	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestScanWithCustomExcludes(t *testing.T) {
	root := t.TempDir()

	createProject(t, root, ".git/payment", "pom.xml")
	createProject(t, root, "applications/customer", "pom.xml")

	projects, err := scanner.ScanWithExcludes(
		root,
		[]string{".git"},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	if projects[0].Name != "customer" {
		t.Errorf("expected customer, got %s", projects[0].Name)
	}
}

func TestScanWithEmptyExcludes(t *testing.T) {
	root := t.TempDir()

	createProject(t, root, ".git/payment", "pom.xml")
	createProject(t, root, "applications/customer", "pom.xml")

	projects, err := scanner.ScanWithExcludes(root, []string{})

	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
}

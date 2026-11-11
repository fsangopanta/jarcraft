package scanner__test

import (
	"jarcraft/internal/scanner"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProject(t *testing.T) {
	tempDir := t.TempDir()

	pomPath := filepath.Join(tempDir, "pom.xml")

	err := os.WriteFile(pomPath, []byte("<project></project>"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	project, err := scanner.DetectProject(tempDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if project == nil {
		t.Fatal("expected project, got nil")
	}

	if project.BuildTool != scanner.Maven {
		t.Errorf("expected Maven, got %s", project.BuildTool)
	}
}

func TestDetectMavenWrapper(t *testing.T) {
	tempDir := t.TempDir()

	createProject(t, tempDir, ".", "pom.xml")

	err := os.WriteFile(
		filepath.Join(tempDir, "mvnw.cmd"),
		[]byte(""),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	project, err := scanner.DetectProject(tempDir)

	if err != nil {
		t.Fatal(err)
	}

	if project != nil {
		if project.Wrapper == nil {
			t.Fatal("expected Maven wrapper, got nil")
		}

		if project.Wrapper.Windows != "mvnw.cmd" {
			t.Errorf("expected mvnw.cmd, got %s", project.Wrapper.Windows)
		}

		if project.Wrapper.Unix != "mvnw" {
			t.Errorf("expected mvnw, got %s", project.Wrapper.Unix)
		}
	}

}

func TestDetectGradleWrapper(t *testing.T) {
	tempDir := t.TempDir()

	createProject(t, tempDir, ".", "build.gradle")

	err := os.WriteFile(
		filepath.Join(tempDir, "gradlew.bat"),
		[]byte(""),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	project, err := scanner.DetectProject(tempDir)

	if err != nil {
		t.Fatal(err)
	}

	if project == nil {
		t.Fatal("expected project, got nil")
	}

	if project.Wrapper == nil {
		t.Fatal("expected Gradle wrapper, got nil")
	}

	if project.Wrapper.Windows != "gradlew.bat" {
		t.Errorf("expected gradlew.bat, got %s", project.Wrapper.Windows)
	}

	if project.Wrapper.Unix != "gradlew" {
		t.Errorf("expected gradlew, got %s", project.Wrapper.Unix)
	}
}

func TestDetectProjectWithoutWrapper(t *testing.T) {
	tempDir := t.TempDir()

	createProject(t, tempDir, ".", "pom.xml")

	project, err := scanner.DetectProject(tempDir)

	if err != nil {
		t.Fatal(err)
	}

	if project != nil {
		if project.Wrapper != nil {
			t.Errorf("expected no wrapper, got %+v", project.Wrapper)
		}
	}

}

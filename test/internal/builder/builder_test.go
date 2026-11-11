package builder_test

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"jarcraft/internal/builder"
	"jarcraft/internal/scanner"
)

func TestCreateMavenCommandWithWindowsWrapper(t *testing.T) {
	project := scanner.Project{
		Name:      "payments",
		Path:      `C:\projects\payments`,
		BuildTool: scanner.Maven,
		Wrapper: &scanner.BuildWrapper{
			Windows: "mvnw.cmd",
			Unix:    "mvnw",
		},
	}

	command, err := builder.CreateCommand(project, "windows")

	if err != nil {
		t.Fatal(err)
	}

	if command.Name != "cmd.exe" {
		t.Errorf("expected cmd.exe, got %s", command.Name)
	}

	if len(command.Args) != 3 {
		t.Fatalf("expected 3 arguments, got %d", len(command.Args))
	}


}

func TestCreateGradleCommandWithUnixWrapper(t *testing.T) {
	project := scanner.Project{
		Name:      "inventory",
		Path:      "/projects/inventory",
		BuildTool: scanner.Gradle,
		Wrapper: &scanner.BuildWrapper{
			Windows: "gradlew.bat",
			Unix:    "gradlew",
		},
	}

	command, err := builder.CreateCommand(project, "linux")

	if err != nil {
		t.Fatal(err)
	}

	expected := "\\projects\\inventory\\gradlew"

	if command.Name != expected {
		t.Errorf("expected %s, got %s", expected, command.Name)
	}

	if len(command.Args) != 2 ||
		command.Args[0] != "clean" ||
		command.Args[1] != "build" {
		t.Errorf("unexpected arguments: %v", command.Args)
	}
}

func TestCreateMavenCommandWithoutWrapper(t *testing.T) {
	project := scanner.Project{
		Name:      "payments",
		Path:      "/projects/payments",
		BuildTool: scanner.Maven,
	}

	command, err := builder.CreateCommand(project, "linux")

	if err != nil {
		t.Fatal(err)
	}

	if command.Name != "mvn" {
		t.Errorf("expected mvn, got %s", command.Name)
	}

	expectedArgs := []string{"clean", "package"}

	if !reflect.DeepEqual(command.Args, expectedArgs) {
		t.Errorf("expected %v, got %v", expectedArgs, command.Args)
	}
}
func TestCreateCommandRejectsUnsupportedBuildTool(t *testing.T) {
	project := scanner.Project{
		Name:      "unknown",
		Path:      "/projects/unknown",
		BuildTool: "UNKNOWN",
	}

	_, err := builder.CreateCommand(project, "linux")

	if err == nil {
		t.Fatal("expected error for unsupported build tool")
	}
}

func TestCreateMavenCommandWithSkipTests(t *testing.T) {
	project := scanner.Project{
		Name:      "orders",
		Path:      "/projects/orders",
		BuildTool: scanner.Maven,
	}

	command, err := builder.CreateCommandWithOptions(
		project,
		"linux",
		builder.BuildOptions{
			SkipTests: true,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"clean", "package", "-DskipTests"}

	if !reflect.DeepEqual(command.Args, expected) {
		t.Errorf("expected %v, got %v", expected, command.Args)
	}
}

func TestCreateGradleCommandWithSkipTests(t *testing.T) {
	project := scanner.Project{
		Name:      "orders",
		Path:      "/projects/orders",
		BuildTool: scanner.Gradle,
	}

	command, err := builder.CreateCommandWithOptions(
		project,
		"linux",
		builder.BuildOptions{
			SkipTests: true,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"clean", "build", "-x", "test"}

	if !reflect.DeepEqual(command.Args, expected) {
		t.Errorf("expected %v, got %v", expected, command.Args)
	}
}
func TestBuildEnvironmentUsesSelectedJava(t *testing.T) {
	javaHome := `C:\Java\jdk-21`

	env := builder.BuildEnvironment(javaHome)

	var javaHomeFound bool
	var pathFound bool

	for _, variable := range env {
		if strings.EqualFold(variable, "JAVA_HOME="+javaHome) {
			javaHomeFound = true
		}

		if strings.HasPrefix(strings.ToUpper(variable), "PATH=") {
			pathFound = true

			expectedBin := strings.ToUpper(
				filepath.Join(javaHome, "bin"),
			)

			if !strings.Contains(strings.ToUpper(variable), expectedBin) {
				t.Errorf("Java bin directory not found in PATH: %s", variable)
			}
		}
	}

	if !javaHomeFound {
		t.Error("JAVA_HOME was not configured correctly")
	}

	if !pathFound {
		t.Error("PATH was not configured")
	}
}

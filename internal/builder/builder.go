package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"jarcraft/internal/scanner"
)

type Command struct {
	Name string
	Args []string
	Dir  string
}

type BuildOptions struct {
	SkipTests   bool
	JavaVersion string
	JavaHome    string
}

func CreateCommand(project scanner.Project, goos string) (Command, error) {
	return CreateCommandWithOptions(project, goos, BuildOptions{})
}

func createWindowsWrapperCommand(wrapperPath string, args []string) Command {
	wrapperName := filepath.Base(wrapperPath)

	command := "call " + wrapperName + " " + strings.Join(args, " ")

	return Command{
		Name: "cmd.exe",
		Args: []string{
			"/D",
			"/C",
			command,
		},
		Dir: filepath.Dir(wrapperPath),
	}
}

// CreateCommandWithOptions CreateCommand determines which command should compile a project.
func CreateCommandWithOptions(
	project scanner.Project,
	goos string,
	options BuildOptions,
) (Command, error) {

	var executable string
	var args []string

	switch project.BuildTool {
	case scanner.Maven:
		args = []string{"clean", "package"}

		if options.SkipTests {
			args = append(args, "-DskipTests")
		}

		_, command, err, done := validateWrapper(project, goos, executable, args)
		if done {
			return command, err
		}

		executable = "mvn"
	case scanner.Gradle:
		args = []string{"clean", "build"}

		if options.SkipTests {
			args = append(args, "-x", "test")
		}

		if project.Wrapper != nil {
			if goos == "windows" && project.Wrapper.Windows != "" {
				executable = filepath.Join(project.Path, project.Wrapper.Windows)
				return createWindowsWrapperCommand(executable, args), nil
			}

			if goos != "windows" && project.Wrapper.Unix != "" {
				return Command{
					Name: filepath.Join(project.Path, project.Wrapper.Unix),
					Args: args,
					Dir:  project.Path,
				}, nil
			}
		}

		executable = "gradle"

	default:
		return Command{}, fmt.Errorf(
			"unsupported build tool: %s",
			project.BuildTool,
		)
	}

	if goos == "windows" {
		return Command{
			Name: "cmd.exe",
			Args: []string{
				"/D",
				"/S",
				"/C",
				executable + " " + strings.Join(args, " "),
			},
			Dir: project.Path,
		}, nil
	}

	return Command{
		Name: executable,
		Args: args,
		Dir:  project.Path,
	}, nil
}

func validateWrapper(project scanner.Project, goos string, executable string, args []string) (string, Command, error, bool) {
	if project.Wrapper != nil {
		if goos == "windows" && project.Wrapper.Windows != "" {
			executable = filepath.Join(project.Path, project.Wrapper.Windows)
			return "", createWindowsWrapperCommand(executable, args), nil, true
		}

		if goos != "windows" && project.Wrapper.Unix != "" {
			return "", Command{
				Name: filepath.Join(project.Path, project.Wrapper.Unix),
				Args: args,
				Dir:  project.Path,
			}, nil, true
		}
	}
	return executable, Command{}, nil, false
}
func validateOS(project scanner.Project, goos string, command Command) (Command, error, bool) {
	if project.Wrapper != nil {
		if goos == "windows" && project.Wrapper.Windows != "" {
			command.Name = "cmd"
			command.Args = append(
				[]string{"/C", filepath.Join(project.Path, project.Wrapper.Windows)},
				command.Args...,
			)
			return command, nil, true
		}

		if goos != "windows" && project.Wrapper.Unix != "" {
			command.Name = filepath.Join(project.Path, project.Wrapper.Unix)
			return command, nil, true
		}
	}
	return Command{}, nil, false
}

// Build executes the compilation of a detected project.
func Build(ctx context.Context, project scanner.Project) error {
	command, err := CreateCommand(project, runtime.GOOS)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, command.Name, command.Args...)
	cmd.Dir = command.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"failed to build project %s: %w",
			project.Name,
			err,
		)
	}

	return nil
}

func BuildWithOptions(
	ctx context.Context,
	project scanner.Project,
	options BuildOptions,
) error {
	command, err := CreateCommandWithOptions(
		project,
		runtime.GOOS,
		options,
	)

	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, command.Name, command.Args...)

	if options.JavaHome != "" {
		cmd.Env = BuildEnvironment(options.JavaHome)
	}
	cmd.Dir = command.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("\n========================================\n")
	fmt.Printf("Building project: %s\n", project.Name)
	fmt.Printf("Command: %s %s\n", cmd.Path, strings.Join(cmd.Args[1:], " "))
	fmt.Printf("Build tool: %s\n", project.BuildTool)
	fmt.Printf("Project path: %s\n", project.Path)
	fmt.Printf("========================================\n")

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n[FAILED] %s\n", project.Name)

		return fmt.Errorf(
			"failed to build project %s: %w",
			project.Name,
			err,
		)
	}

	fmt.Printf("\n[SUCCESS] %s compiled successfully!\n", project.Name)

	return nil
}

func BuildEnvironment(javaHome string) []string {
	env := os.Environ()

	pathValue := os.Getenv("PATH")
	javaBin := filepath.Join(javaHome, "bin")

	newPath := javaBin + string(os.PathListSeparator) + pathValue

	env = setEnvironmentVariable(env, "JAVA_HOME", javaHome)
	env = setEnvironmentVariable(env, "PATH", newPath)

	return env
}

func setEnvironmentVariable(env []string, key, value string) []string {
	prefix := key + "="
	result := make([]string, 0, len(env)+1)

	for _, variable := range env {
		if strings.EqualFold(variable[:min(len(variable), len(prefix))], prefix) {
			continue
		}

		result = append(result, variable)
	}

	return append(result, prefix+value)
}

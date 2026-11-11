package scanner

import (
	"os"
	"path/filepath"
)

func DetectProject(path string) (*Project, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make(map[string]bool)

	for _, entry := range entries {
		if !entry.IsDir() {
			files[entry.Name()] = true
		}
	}

	project := &Project{
		Name: filepath.Base(path),
		Path: path,
	}

	switch {
	case files["pom.xml"]:
		project.BuildTool = Maven
		project.BuildFile = "pom.xml"
		project.Wrapper = detectMavenWrapper(files)

	case files["build.gradle"]:
		project.BuildTool = Gradle
		project.BuildFile = "build.gradle"
		project.Wrapper = detectGradleWrapper(files)

	case files["build.gradle.kts"]:
		project.BuildTool = Gradle
		project.BuildFile = "build.gradle.kts"
		project.Wrapper = detectGradleWrapper(files)

	default:
		return nil, nil
	}

	return project, nil
}

func detectMavenWrapper(files map[string]bool) *BuildWrapper {
	wrapper := &BuildWrapper{
		Windows: "mvnw.cmd",
		Unix:    "mvnw",
	}

	if !files[wrapper.Windows] && !files[wrapper.Unix] {
		return nil
	}

	return wrapper
}

func detectGradleWrapper(files map[string]bool) *BuildWrapper {
	wrapper := &BuildWrapper{
		Windows: "gradlew.bat",
		Unix:    "gradlew",
	}

	if !files[wrapper.Windows] && !files[wrapper.Unix] {
		return nil
	}

	return wrapper
}

package javamanager

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func DiscoverInstallations() map[string]string {
	installations := make(map[string]string)

	if runtime.GOOS != "windows" {
		return installations
	}

	var candidates []string

	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		candidates = append(candidates, javaHome)
	}

	programFiles := os.Getenv("ProgramFiles")

	if programFiles != "" {
		candidates = append(
			candidates,
			filepath.Join(programFiles, "Java"),
			filepath.Join(programFiles, "Eclipse Adoptium"),
		)
	}

	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		candidates = append(
			candidates,
			filepath.Join(userProfile, ".jdks"),
		)
	}

	for _, candidate := range candidates {
		entries, err := os.ReadDir(candidate)
		if err != nil {
			// También comprobamos si el candidato es directamente un JDK.
			addInstallation(installations, candidate)
			continue
		}

		addInstallation(installations, candidate)

		for _, entry := range entries {
			if entry.IsDir() {
				addInstallation(
					installations,
					filepath.Join(candidate, entry.Name()),
				)
			}
		}
	}

	return installations
}

func addInstallation(installations map[string]string, home string) {
	if !isValidJDK(home) {
		return
	}

	version := detectVersion(home)

	if version != "" {
		installations[version] = home
	}
}

func isValidJDK(home string) bool {
	javaPath := filepath.Join(home, "bin", "java.exe")
	javacPath := filepath.Join(home, "bin", "javac.exe")

	if _, err := os.Stat(javaPath); err != nil {
		return false
	}

	if _, err := os.Stat(javacPath); err != nil {
		return false
	}

	return true
}

func detectVersion(home string) string {
	javaPath := filepath.Join(home, "bin", "java.exe")

	cmd := exec.Command(javaPath, "-version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	text := string(output)

	// Java 8: version "1.8.0_xxx"
	// Java 17+: version "17.0.x"
	start := strings.Index(text, `version "`)

	if start == -1 {
		return ""
	}

	text = text[start+len(`version "`):]

	end := strings.Index(text, `"`)

	if end == -1 {
		return ""
	}

	version := text[:end]

	if strings.HasPrefix(version, "1.") {
		version = strings.TrimPrefix(version, "1.")
		version = strings.Split(version, ".")[0]
	} else {
		version = strings.Split(version, ".")[0]
	}

	return version
}

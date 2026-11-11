package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

var defaultExcludedDirectories = []string{
	".git",
	".idea",
	".gradle",
	"node_modules",
}

func Scan(root string) ([]Project, error) {
	return ScanWithExcludes(root, defaultExcludedDirectories)
}

func ScanWithExcludes(root string, excludes []string) ([]Project, error) {
	absolutePath, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve root path: %w", err)
	}

	info, err := os.Stat(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("access root directory: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("root path is not a directory: %s", absolutePath)
	}

	excludedDirectories := make(map[string]struct{}, len(excludes))

	for _, name := range excludes {
		if name != "" {
			excludedDirectories[name] = struct{}{}
		}
	}

	var projects []Project

	err = filepath.WalkDir(absolutePath, func(
		path string,
		entry os.DirEntry,
		walkErr error,
	) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() && path != absolutePath {
			if shouldSkipDirectory(entry.Name(), excludedDirectories) {
				return filepath.SkipDir
			}
		}

		if !entry.IsDir() {
			return nil
		}

		project, err := DetectProject(path)
		if err != nil {
			return fmt.Errorf("detect project in %s: %w", path, err)
		}

		if project != nil {
			projects = append(projects, *project)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("scan projects: %w", err)
	}

	return projects, nil
}

func shouldSkipDirectory(
	name string,
	excludedDirectories map[string]struct{},
) bool {
	_, excluded := excludedDirectories[name]
	return excluded
}

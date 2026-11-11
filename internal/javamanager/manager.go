package javamanager

import (
	"fmt"
	"maps"
	"sort"
	"strconv"
)

type Manager struct {
	installations map[string]string
}

type JavaInstallation struct {
	Version string
	Home    string
}

func New(installations map[string]string) *Manager {
	return &Manager{
		installations: installations,
	}
}

func (m *Manager) Resolve(version string) (*JavaInstallation, error) {
	available := DiscoverInstallations()

	maps.Copy(available, m.installations)

	if version == "" {
		versions := make([]int, 0, len(available))

		for v := range available {
			parsed, err := strconv.Atoi(v)
			if err == nil {
				versions = append(versions, parsed)
			}
		}

		if len(versions) == 0 {
			return nil, fmt.Errorf("no Java installations were found")
		}

		sort.Sort(sort.Reverse(sort.IntSlice(versions)))

		version = strconv.Itoa(versions[0])
	}

	home, exists := available[version]
	if !exists {
		return nil, fmt.Errorf(
			"java %s was not found in configured or discovered installations",
			version,
		)
	}

	return validateInstallation(version, home)
}

func validateInstallation(version, home string) (*JavaInstallation, error) {
	if !isValidJDK(home) {
		return nil, fmt.Errorf("invalid JDK installation: %s", home)
	}

	detectedVersion := detectVersion(home)

	if detectedVersion == "" {
		return nil, fmt.Errorf("could not detect Java version at %s", home)
	}

	if detectedVersion != version {
		return nil, fmt.Errorf(
			"requested Java %s but found Java %s at %s",
			version,
			detectedVersion,
			home,
		)
	}

	return &JavaInstallation{
		Version: detectedVersion,
		Home:    home,
	}, nil
}

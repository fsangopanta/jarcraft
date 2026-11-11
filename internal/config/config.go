package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Build    BuildConfig              `mapstructure:"build"`
	Scanner  ScannerConfig            `mapstructure:"scanner"`
	Defaults ProjectConfig            `mapstructure:"defaults"`
	Projects map[string]ProjectConfig `mapstructure:"projects"`
	Logging  LoggingConfig            `mapstructure:"logging"`
	Java     JavaConfig               `mapstructure:"java"`
}

type BuildConfig struct {
	Output   string `mapstructure:"output"`
	Parallel bool   `mapstructure:"parallel"`
	FailFast bool   `mapstructure:"fail_fast"`
}

type ScannerConfig struct {
	Exclude []string `mapstructure:"exclude"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

type ProjectConfig struct {
	SkipTests   *bool  `mapstructure:"skip_tests"`
	JavaVersion string `mapstructure:"java_version"`
}

type ResolvedProjectConfig struct {
	SkipTests   bool
	JavaVersion string
}

type BuildOverrides struct {
	SkipTests   *bool
	JavaVersion string
}

type JavaConfig struct {
	Installations map[string]string `mapstructure:"installations"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	setDefaults(v)

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("jarcraft")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError

		if !errors.As(err, &notFound) || configPath != "" {
			return nil, fmt.Errorf("read configuration: %w", err)
		}
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("build.output", "./builds")
	v.SetDefault("build.parallel", false)
	v.SetDefault("build.fail_fast", false)

	v.SetDefault("scanner.exclude", []string{
		".git",
		".idea",
		".gradle",
		"node_modules",
	})

	v.SetDefault("logging.level", "info")
}

func (c *Config) Validate() error {
	c.Logging.Level = strings.ToLower(strings.TrimSpace(c.Logging.Level))

	switch c.Logging.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf(
			"unsupported logging level %q: expected debug, info, warn or error",
			c.Logging.Level,
		)
	}

	c.Build.Output = strings.TrimSpace(c.Build.Output)

	if c.Build.Output == "" {
		return errors.New("build.output cannot be empty")
	}

	return nil
}

func (c *Config) ResolveProjectConfig(
	relativePath string,
) ResolvedProjectConfig {

	resolved := ResolvedProjectConfig{
		SkipTests:   false,
		JavaVersion: "",
	}

	if c.Defaults.SkipTests != nil {
		resolved.SkipTests = *c.Defaults.SkipTests
	}

	if c.Defaults.JavaVersion != "" {
		resolved.JavaVersion = c.Defaults.JavaVersion
	}

	normalizedPath := filepath.ToSlash(filepath.Clean(relativePath))

	projectConfig, exists := c.Projects[normalizedPath]
	if !exists {
		return resolved
	}

	if projectConfig.SkipTests != nil {
		resolved.SkipTests = *projectConfig.SkipTests
	}

	if projectConfig.JavaVersion != "" {
		resolved.JavaVersion = projectConfig.JavaVersion
	}

	return resolved
}

func (c *Config) ResolveProjectConfigWithOverrides(
	relativePath string,
	overrides BuildOverrides,
) ResolvedProjectConfig {

	resolved := c.ResolveProjectConfig(relativePath)

	if overrides.SkipTests != nil {
		resolved.SkipTests = *overrides.SkipTests
	}

	if overrides.JavaVersion != "" {
		resolved.JavaVersion = overrides.JavaVersion
	}

	return resolved
}

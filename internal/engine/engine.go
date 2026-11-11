package engine

import (
	"context"
	"fmt"
	"jarcraft/internal/config"
	"jarcraft/internal/javamanager"
	"path/filepath"

	"jarcraft/internal/builder"
	"jarcraft/internal/scanner"
)

type BuildResult struct {
	Project scanner.Project
	Success bool
	Error   error
}

type BuildFunc func(context.Context, scanner.Project) error

type Engine struct {
	build           BuildFunc
	configuredBuild ConfiguredBuildFunc
}

type ConfiguredBuildFunc func(
	context.Context,
	scanner.Project,
	config.ResolvedProjectConfig,
) error

func New() *Engine {
	return NewWithBuilder(builder.Build)
}

func NewWithBuilder(build BuildFunc) *Engine {
	return &Engine{
		build: build,
	}
}

func NewWithConfiguredBuilder(build ConfiguredBuildFunc) *Engine {
	return &Engine{
		configuredBuild: build,
	}
}

type ProjectBuild struct {
	Project scanner.Project
	Config  config.ResolvedProjectConfig
}

func (e *Engine) Run(
	ctx context.Context,
	root string,
) ([]BuildResult, error) {
	projects, err := scanner.Scan(root)
	if err != nil {
		return nil, err
	}

	return e.buildProjects(ctx, projects)
}

func (e *Engine) RunWithExcludes(
	ctx context.Context,
	root string,
	excludes []string,
) ([]BuildResult, error) {
	projects, err := scanner.ScanWithExcludes(root, excludes)
	if err != nil {
		return nil, fmt.Errorf("scan projects: %w", err)
	}

	return e.buildProjects(ctx, projects)
}

func (e *Engine) buildProjects(
	ctx context.Context,
	projects []scanner.Project,
) ([]BuildResult, error) {
	results := make([]BuildResult, 0, len(projects))

	for _, project := range projects {
		if err := ctx.Err(); err != nil {
			return results, err
		}

		result := BuildResult{
			Project: project,
		}

		err := e.build(ctx, project)

		if err != nil {
			result.Error = err
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results, nil
}

func (e *Engine) buildProjectsWithConfig(
	ctx context.Context,
	root string,
	projects []scanner.Project,
	cfg *config.Config,
	overrides config.BuildOverrides,
) ([]BuildResult, error) {

	if e.configuredBuild == nil {
		return nil, fmt.Errorf("configured builder is not configured")
	}

	results := make([]BuildResult, 0, len(projects))

	for _, project := range projects {
		if err := ctx.Err(); err != nil {
			return results, err
		}

		relativePath, err := filepath.Rel(root, project.Path)
		if err != nil {
			return results, fmt.Errorf(
				"resolve relative path for project %s: %w",
				project.Name,
				err,
			)
		}

		relativePath = filepath.ToSlash(relativePath)

		projectConfig := cfg.ResolveProjectConfigWithOverrides(
			relativePath,
			overrides,
		)

		result := BuildResult{
			Project: project,
		}

		err = e.configuredBuild(ctx, project, projectConfig)

		if err != nil {
			result.Error = err
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results, nil
}

func (e *Engine) RunWithConfig(
	ctx context.Context,
	root string,
	cfg *config.Config,
	overrides config.BuildOverrides,
) ([]BuildResult, error) {

	projects, err := scanner.ScanWithExcludes(
		root,
		cfg.Scanner.Exclude,
	)
	if err != nil {
		return nil, fmt.Errorf("scan projects: %w", err)
	}

	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve root path: %w", err)
	}

	return e.buildProjectsWithConfig(
		ctx,
		absoluteRoot,
		projects,
		cfg,
		overrides,
	)
}

func NewConfigured(javaConfig config.JavaConfig) *Engine {
	javaManager := javamanager.New(javaConfig.Installations)

	return NewWithConfiguredBuilder(
		func(
			ctx context.Context,
			project scanner.Project,
			options config.ResolvedProjectConfig,
		) error {

			javaInstallation, err := javaManager.Resolve(options.JavaVersion)
			if err != nil {
				return fmt.Errorf(
					"project %s: %w",
					project.Name,
					err,
				)
			}

			return builder.BuildWithOptions(
				ctx,
				project,
				builder.BuildOptions{
					SkipTests:   options.SkipTests,
					JavaVersion: options.JavaVersion,
					JavaHome:    javaInstallation.Home,
				},
			)
		},
	)
}

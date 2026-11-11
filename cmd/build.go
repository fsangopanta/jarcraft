package cmd

import (
	"fmt"
	"jarcraft/internal/config"

	"jarcraft/internal/engine"

	"github.com/spf13/cobra"
)

func newBuildCmd() *cobra.Command {
	var rootPath string

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Scan and build Java projects",
		Long:  "Detects Maven and Gradle projects in a directory and builds them automatically.",

		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := getConfig(cmd)

			skipTestsValue, err := cmd.Flags().GetBool("skip-tests")
			if err != nil {
				return err
			}

			javaVersion, err := cmd.Flags().GetString("java-version")
			if err != nil {
				return err
			}

			overrides := config.BuildOverrides{
				JavaVersion: javaVersion,
			}

			if cmd.Flags().Changed("skip-tests") {
				overrides.SkipTests = &skipTestsValue
			}

			projectEngine := engine.NewConfigured(cfg.Java)

			results, err := projectEngine.RunWithConfig(
				cmd.Context(),
				rootPath,
				cfg,
				overrides,
			)

			if err != nil {
				return err
			}

			if len(results) == 0 {
				cmd.Println("No Java projects found.")
				return nil
			}

			var successful, failed int

			cmd.Println("Build results:")
			cmd.Println("------------------------------")

			for _, result := range results {
				if result.Success {
					successful++
					cmd.Printf("[SUCCESS] %s (%s)\n",
						result.Project.Name,
						result.Project.BuildTool,
					)
				} else {
					failed++
					cmd.Printf("[FAILED] %s (%s)\n",
						result.Project.Name,
						result.Project.BuildTool,
					)
					cmd.Printf("          Error: %v\n", result.Error)
				}
			}

			cmd.Println("------------------------------")
			cmd.Printf("Total: %d | Successful: %d | Failed: %d\n",
				len(results), successful, failed,
			)

			if failed > 0 {
				return fmt.Errorf("%d project(s) failed to build", failed)
			}

			return nil
		},
	}

	buildCmd.Flags().StringVarP(
		&rootPath,
		"path",
		"p",
		".",
		"Root directory to scan for Java projects",
	)

	buildCmd.Flags().Bool("skip-tests", false, "Skip tests during the build")
	buildCmd.Flags().String("java-version", "", "Java version to use")

	return buildCmd
}

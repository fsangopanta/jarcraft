package cmd

import (
	"context"

	"jarcraft/internal/config"

	"github.com/spf13/cobra"
)

type configKey struct{}

func newRootCmd() *cobra.Command {
	var configPath string

	rootCmd := &cobra.Command{
		Use:           "jarcraft",
		Short:         "Automatically build Java projects",
		Long:          "JarCraft scans directories, detects Java projects and builds them using Maven or Gradle.",
		SilenceUsage:  true,
		SilenceErrors: true,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}

			ctx := context.WithValue(cmd.Context(), configKey{}, cfg)
			cmd.SetContext(ctx)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(
		&configPath,
		"config",
		"",
		"Path to configuration file",
	)

	rootCmd.AddCommand(newBuildCmd())

	return rootCmd
}

func Execute() error {
	return newRootCmd().Execute()
}

func getConfig(cmd *cobra.Command) *config.Config {
	cfg, _ := cmd.Context().Value(configKey{}).(*config.Config)
	return cfg
}

package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a static website",
	Long:  `Build a static website.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildToFs(cfg, traceFS(cfg, afero.NewBasePathFs(afero.NewOsFs(), cfg.Build.Out)))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("build.out", "o", "out", "path to output")

	viper.BindPFlags(buildCmd.Flags())
}

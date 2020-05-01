package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/config"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a static website",
	Long:  `Build a static website.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}
		return buildToFs(cfg, afero.NewBasePathFs(afero.NewOsFs(), cfg.Build.Out))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("build.out", "o", "out", "path to output")

	viper.BindPFlags(buildCmd.Flags())
}

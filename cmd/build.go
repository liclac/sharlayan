package cmd

import (
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/render"
	"github.com/liclac/sharlayan/tree"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a static website",
	Long:  `Build a static website.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		meta, err := calibre.Read(cfg.Library)
		if err != nil {
			return err
		}
		fs := osfs.New(cfg.Build.Out)
		return tree.Render(fs, false, "/", render.Root(cfg, meta))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("build.out", "o", "out", "path to output")

	viper.BindPFlags(buildCmd.Flags())
}

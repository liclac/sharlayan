package cmd

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/builder"
	"github.com/liclac/sharlayan/builder/tree"
	"github.com/liclac/sharlayan/calibre"
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
		bld, err := builder.New(cfg)
		if err != nil {
			return err
		}
		root := builder.Root(bld, meta)
		if root == nil {
			return fmt.Errorf("root == nil, nothing to render")
		}
		fs := traceFS(cfg, afero.NewBasePathFs(afero.NewOsFs(), cfg.Build.Out))
		if err := root.Render(fs, tree.ByID, "/_id/"); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("build.out", "o", "out", "path to output")

	viper.BindPFlags(buildCmd.Flags())
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/render"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a static website",
	Long:  `Build a static website.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		lpath := viper.GetString("library")
		if lpath == "" {
			return fmt.Errorf("-l/--library is required")
		}

		var cfg render.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}
		meta, err := calibre.Read(lpath)
		if err != nil {
			return err
		}
		return render.Render(cfg, meta)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("out", "o", "www", "path to output")
	buildCmd.Flags().String("templates", "templates", "path to templates")
	buildCmd.Flags().String("title", "My Library", "title for rendered site")

	viper.BindPFlags(buildCmd.Flags())
}

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
		outPath := viper.GetString("out")
		templatePath := viper.GetString("templates")
		cfg := render.Config{
			Title: viper.GetString("title"),
		}

		meta, err := calibre.Read(lpath)
		if err != nil {
			return err
		}
		return render.Render(outPath, templatePath, meta, cfg)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().String("title", "My Library", "title for rendered site")
	buildCmd.Flags().StringP("out", "o", "www", "path to output")
	buildCmd.Flags().String("templates", "templates", "path to templates")

	viper.BindPFlags(buildCmd.Flags())
}

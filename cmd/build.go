package cmd

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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

		metaStart := time.Now()
		meta, err := calibre.Read(lpath)
		if err != nil {
			return err
		}
		metaTime := time.Since(metaStart)
		log.WithFields(log.Fields{
			"books": len(meta.Books),
			"t":     metaTime,
		}).Debug("Loaded: Calibre database")

		rStart := time.Now()
		r, err := render.New(cfg, meta)
		if err != nil {
			return err
		}
		rTime := time.Since(rStart)
		log.WithField("t", rTime).Debug("Loaded: Templates")

		rootStart := time.Now()
		root := render.Root(cfg, meta)
		rootTime := time.Since(rootStart)
		log.WithField("t", rootTime).Debug("Loaded: Node Tree")

		renderStart := time.Now()
		pages, err := r.Render(root)
		if err != nil {
			return err
		}
		renderTime := time.Since(renderStart)
		log.WithFields(log.Fields{
			"pages": pages,
			"t":     renderTime,
		}).Debug("Rendered!")
		return err
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("out", "o", "www", "path to output")
	buildCmd.Flags().String("templates", "templates", "path to templates")

	buildCmd.Flags().String("title", "My Library", "title for rendered site")
	buildCmd.Flags().Bool("author.no-index", false, "disable author index")
	buildCmd.Flags().Bool("series.no-index", false, "disable series index")
	buildCmd.Flags().Bool("tag.no-index", false, "disable tag index")

	viper.BindPFlags(buildCmd.Flags())
}

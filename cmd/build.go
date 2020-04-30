package cmd

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/builder"
	"github.com/liclac/sharlayan/calibre"
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

		metaStart := time.Now()
		meta, err := calibre.Read(cfg.Library)
		if err != nil {
			return err
		}
		metaTime := time.Since(metaStart)
		log.WithFields(log.Fields{
			"books": len(meta.Books),
			"t":     metaTime,
		}).Debug("Loaded: Calibre Metadata")

		bldStart := time.Now()
		bld, err := builder.New(cfg)
		if err != nil {
			return err
		}
		bldTime := time.Since(bldStart)
		log.WithField("t", bldTime).Debug("Loaded: Builder")

		nodesStart := time.Now()
		nodes := bld.Nodes(meta)
		nodesTime := time.Since(nodesStart)
		log.WithFields(log.Fields{
			"num": len(nodes),
			"t":   nodesTime,
		}).Debug("Loaded: Nodes")

		buildStart := time.Now()
		if err := bld.Build(afero.NewBasePathFs(afero.NewOsFs(), cfg.Out), nodes); err != nil {
			return err
		}
		buildTime := time.Since(buildStart)
		log.WithField("t", buildTime).Debug("Rendered!")
		return err
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("out", "o", "out", "path to output")

	buildCmd.Flags().String("html.templates", "templates", "path to templates")
	buildCmd.Flags().String("html.title", "My Library", "title for rendered site")

	buildCmd.Flags().String("books.path", "/books", "output path to books")
	buildCmd.Flags().String("authors.path", "/authors", "output path to authors")
	buildCmd.Flags().String("series.path", "/series", "output path to series")
	buildCmd.Flags().String("tags.path", "/tags", "output path to tags")

	viper.BindPFlags(buildCmd.Flags())
}

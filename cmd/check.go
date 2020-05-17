package cmd

import (
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/config"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check library consistency",
	Long:  `Check library consistency.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		meta, err := calibre.Read(cfg.Library)
		if err != nil {
			return err
		}
		report, err := meta.Check()
		if err != nil {
			return err
		}

		// Prepare file reports for display, in order.
		var paths []string
		for path := range report.Files {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			status := report.Files[path]
			fn := color.Green
			if status == calibre.FileStatusMissing {
				fn = color.Red
			} else if status == calibre.FileStatusOrphan {
				fn = color.Yellow
			}
			fn("%s %s\n", status, path)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/config"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump database contents",
	Long:  `Dump database contents.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		meta, err := calibre.Read(cfg.Library)
		if err != nil {
			return err
		}
		return dump(meta)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/calibre"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump database contents",
	Long:  `Dump database contents.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		lpath := viper.GetString("library")
		if lpath == "" {
			return fmt.Errorf("-l/--library is required")
		}

		meta, err := calibre.Read(lpath)
		if err != nil {
			return err
		}
		return dump(meta)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}

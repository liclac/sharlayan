package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liclac/sharlayan/calibre"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List database contents",
	Long:  `List database contents.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		lpath := viper.GetString("library")
		if lpath == "" {
			return fmt.Errorf("-l/--library is required")
		}

		lib, err := calibre.Open(lpath)
		if err != nil {
			return err
		}
		books, err := lib.Books()
		if err != nil {
			return err
		}
		return dump(books)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

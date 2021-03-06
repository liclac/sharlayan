/*
Copyright © 2020 embr <hi@liclac.eu>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/liclac/sharlayan/config"
)

var cfg = &config.Config{}

var rootCmd = &cobra.Command{
	Use:   "sharlayan",
	Short: "What's in your library?",
	Long:  `What's in your library?`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.Unmarshal(cfg); err != nil {
			return err
		}

		// Feel free to change the logging setup as you see fit.
		logLevel := zapcore.InfoLevel
		if cfg.Verbose {
			logLevel = zapcore.DebugLevel
		}
		zap.ReplaceGlobals(zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
				MessageKey:     "M",
				LevelKey:       "L",
				TimeKey:        zapcore.OmitKey,
				NameKey:        "N",
				CallerKey:      "C",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.EpochTimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			os.Stderr, logLevel)))
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	defaultCfgDir := filepath.Join(userCfgDir, "sharlayan")

	rootCmd.PersistentFlags().StringP("library", "l", filepath.Join(home, "Calibre Library"), "path to calibre library")
	rootCmd.PersistentFlags().StringVarP(&cfg.Config.Dir, "config.dir", "C", defaultCfgDir, "path to config directory")
	rootCmd.PersistentFlags().StringVarP(&cfg.Config.File, "config.file", "c", "", "path to config file (default ${config.dir}/config.toml)")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable debug logging")
	rootCmd.PersistentFlags().Bool("debug.trace-fs", false, "log all filesystem calls")

	rootCmd.PersistentFlags().String("html.templates", "templates", "path to templates")
	rootCmd.PersistentFlags().String("html.root", "", "public path to library root")
	rootCmd.PersistentFlags().String("html.title", "My Library", "title for rendered site")

	rootCmd.PersistentFlags().String("books.path", "/books", "output path to books")
	rootCmd.PersistentFlags().String("authors.path", "/authors", "output path to authors")
	rootCmd.PersistentFlags().String("series.path", "/series", "output path to series")
	rootCmd.PersistentFlags().String("tags.path", "/tags", "output path to tags")

	viper.BindPFlags(rootCmd.PersistentFlags())
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile := cfg.Config.File; cfgFile != "" {
		viper.SetConfigFile(cfg.Config.File)
	} else {
		viper.AddConfigPath(cfg.Config.Dir)
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, isNotFound := err.(viper.ConfigFileNotFoundError); !isNotFound {
			panic(err)
		}
	}
}

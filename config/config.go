package config

import (
	"time"
)

type Config struct {
	Library string `mapstructure:"library"` // Path to Calibre library.
	Config  struct {
		File string `mapstructure:"file"` // Path to config.
		Dir  string `mapstructure:"dir"`  // Path to config directory.
	} `mapstructure:"config"`

	// Debug options.
	Verbose bool `mapstructure:"verbose"` // Enable debug logging.
	Debug   struct {
		TraceFS bool `mapstructure:"trace-fs"` // Log all filesystem operations.
	} `mapstructure:"debug"`

	// Build command specific.
	Build struct {
		Out    string `mapstructure:"out"` // Output directory.
		ByName struct {
			Enable bool   `mapstructure:"enable"` // Enable human readable names.
			Prefix string `mapstructure:"prefix"` // Prefix in the output tree.
		} `mapstructure:"by-name"`
		ByID struct {
			Enable bool   `mapstructure:"enable"` // Enable stable ID paths.
			Prefix string `mapstructure:"prefix"` // Prefix in the output tree.
		} `mapstructure:"by-id"`
	} `mapstructure:"build"`

	// Serve command specific.
	HTTP struct {
		Enable bool          `mapstructure:"enable"` // Enable the HTTP server.
		Addr   string        `mapstructure:"addr"`   // Address to listen on.
		Grace  time.Duration `mapstructure:"grace"`  // Shutdown grace period.
	} `mapstructure:"http"`

	SSH struct {
		Enable  bool   `mapstructure:"enable"`   // Enable the SSH server.
		Addr    string `mapstructure:"addr"`     // Address to listen on.
		HostKey string `mapstructure:"host-key"` // Path to host private key.
		Trace   bool   `mapstructure:"trace"`    // Log all SSH operations.

		// SSH subsystems.
		SFTP struct {
			Enable bool `mapstructure:"enable"` // Enable the SFTP subsystem.
		}
	}

	// Output formats.
	HTML struct {
		Templates string `mapstructure:"templates"` // Template source directory.
		Root      string `mapstructure:"root"`      // Prefix from the root of your site.
		Title     string `mapstructure:"title"`     // Site title.
	} `mapstructure:"html"`
}

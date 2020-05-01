package config

type Config struct {
	Library string `mapstructure:"library"`

	// Build command specific.
	Build struct {
		Out string `mapstructure:"out"`
	} `mapstructure:"build"`

	// Serve protocol specific.
	HTTP struct {
		Enable bool   `mapstructure:"enable"`
		Addr   string `mapstructure:"addr"`
	} `mapstructure:"http"`

	// Output formats.
	HTML struct {
		Templates string `mapstructure:"templates"`
		Root      string `mapstructure:"root"`
		Title     string `mapstructure:"title"`
	} `mapstructure:"html"`

	// Collections.
	Books struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"books"`

	Authors struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"authors"`

	Series struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"series"`

	Tags struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"tags"`
}

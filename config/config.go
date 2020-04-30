package config

type Config struct {
	Out     string `mapstructure:"out"`
	Library string `mapstructure:"library"`

	// Output formats.
	HTML struct {
		Templates string `mapstructure:"templates"`
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

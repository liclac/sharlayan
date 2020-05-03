package builder

import (
	"github.com/liclac/sharlayan/builder/html"
	"github.com/liclac/sharlayan/config"
)

type Builder struct {
	Cfg  *config.Config
	HTML *html.Builder
}

func New(cfg *config.Config) (*Builder, error) {
	htmlBuilder, err := html.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Builder{
		Cfg:  cfg,
		HTML: htmlBuilder,
	}, nil
}

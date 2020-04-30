package html

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"

	"gopkg.in/russross/blackfriday.v2"

	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/config"
)

type Funcs struct {
	Config *config.Config
}

func NewFuncs(cfg *config.Config) Funcs {
	return Funcs{Config: cfg}
}

func (f Funcs) Map() template.FuncMap {
	return template.FuncMap{
		"cfg":      f.Cfg,
		"markdown": f.Markdown,
		"link":     f.Link,
		"linkTo":   f.LinkTo,
		"linksTo":  f.LinksTo,
	}
}

func (f Funcs) Cfg() *config.Config {
	return f.Config
}

func (f Funcs) Markdown(v string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(v)))
}

func (f Funcs) Link(text string, parts ...string) Link {
	return Link{Href: filepath.Join(parts...), Text: text}
}

func (f Funcs) LinkTo(iv interface{}) (Link, error) {
	switch v := iv.(type) {
	case Link:
		return v, nil
	case *calibre.Book:
		return f.Link(v.Title, f.Config.Books.Path, strconv.Itoa(v.ID)), nil
	case *calibre.Author:
		return f.Link(v.Name, f.Config.Authors.Path, strconv.Itoa(v.ID)), nil
	case *calibre.Series:
		return f.Link(v.Name, f.Config.Series.Path, strconv.Itoa(v.ID)), nil
	case *calibre.Tag:
		return f.Link(v.Name, f.Config.Tags.Path, strconv.Itoa(v.ID)), nil
	}
	return Link{}, fmt.Errorf("linkTo supports Link, *Book, *Author, *Series and *Tag, not %T", iv)
}

func (f Funcs) LinksTo(ivs ...interface{}) ([]Link, error) {
	links := make([]Link, len(ivs))
	for i, iv := range links {
		link, err := f.LinkTo(iv)
		if err != nil {
			return nil, err
		}
		links[i] = link
	}
	return links, nil
}

package render

import (
	"fmt"
	"html/template"
	"strconv"

	"gopkg.in/russross/blackfriday.v2"

	"github.com/liclac/sharlayan/calibre"
)

type Funcs struct {
	Config *Config
}

func NewFuncs(cfg *Config) Funcs {
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

func (f Funcs) Cfg() *Config {
	return f.Config
}

func (f Funcs) Markdown(v string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(v)))
}

func (f Funcs) Link(path, text string) Link {
	return Link{Href: path, Text: text}
}

func (f Funcs) LinkTo(iv interface{}) (Link, error) {
	switch v := iv.(type) {
	case Link:
		return v, nil
	case *calibre.Book:
		return f.Link("/book/"+strconv.Itoa(v.ID), v.Title), nil
	case *calibre.Author:
		return f.Link("/author/"+strconv.Itoa(v.ID), v.Name), nil
	case *calibre.Series:
		return f.Link("/series/"+strconv.Itoa(v.ID), v.Name), nil
	case *calibre.Tag:
		return f.Link("/tag/"+strconv.Itoa(v.ID), v.Name), nil
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

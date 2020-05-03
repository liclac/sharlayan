package html

import (
	"fmt"
	"html/template"

	"gopkg.in/russross/blackfriday.v2"

	"github.com/liclac/sharlayan/builder/tree"
	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/config"
)

// A Link used by generic '_nav' lists.
type Link struct {
	f     *Funcs
	Abs   bool
	Infos []tree.NodeInfo
}

func (l Link) Href() string {
	href := tree.Path(l.f.Naming, l.Infos...)
	if l.Abs {
		href = "/" + href
	}
	return href
}

func (l Link) Text() string {
	return l.Infos[len(l.Infos)-1].Name
}

type Funcs struct {
	Config *config.Config
	Naming tree.NamingScheme
}

func NewFuncs(cfg *config.Config) Funcs {
	return Funcs{Config: cfg}
}

func (f Funcs) Map() template.FuncMap {
	return template.FuncMap{
		"cfg":      f.Cfg,
		"markdown": f.Markdown,
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

func (f *Funcs) LinkTo(iv interface{}) (Link, error) {
	switch v := iv.(type) {
	case Link:
		return v, nil
	case tree.NodeInfo:
		return Link{f, false, []tree.NodeInfo{v}}, nil
	case *calibre.Book:
		return Link{f, true, []tree.NodeInfo{tree.BooksDirInfo, tree.BookInfo(v)}}, nil
	case *calibre.Author:
		return Link{f, true, []tree.NodeInfo{tree.AuthorsDirInfo, tree.AuthorInfo(v)}}, nil
	case *calibre.Series:
		return Link{f, true, []tree.NodeInfo{tree.SeriesDirInfo, tree.SeriesInfo(v)}}, nil
	case *calibre.Tag:
		return Link{f, true, []tree.NodeInfo{tree.TagsDirInfo, tree.TagInfo(v)}}, nil
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

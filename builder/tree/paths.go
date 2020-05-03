package tree

import (
	"path/filepath"
	"strconv"

	"github.com/liclac/sharlayan/calibre"
)

var (
	BooksDirInfo   = NodeInfo{ID: "books", Name: "Books"}
	AuthorsDirInfo = NodeInfo{ID: "authors", Name: "Authors"}
	SeriesDirInfo  = NodeInfo{ID: "series", Name: "Series"}
	TagsDirInfo    = NodeInfo{ID: "tags", Name: "Tags"}
)

func BookInfo(b *calibre.Book) NodeInfo     { return NodeInfo{ID: strconv.Itoa(b.ID), Name: b.Title} }
func AuthorInfo(a *calibre.Author) NodeInfo { return NodeInfo{ID: strconv.Itoa(a.ID), Name: a.Name} }
func SeriesInfo(s *calibre.Series) NodeInfo { return NodeInfo{ID: strconv.Itoa(s.ID), Name: s.Name} }
func TagInfo(t *calibre.Tag) NodeInfo       { return NodeInfo{ID: strconv.Itoa(t.ID), Name: t.Name} }

func Path(ns NamingScheme, infos ...NodeInfo) string {
	parts := make([]string, len(infos))
	for i, info := range infos {
		parts[i] = info.Filename(ns)
	}
	return filepath.Join(parts...)
}

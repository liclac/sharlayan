package render

import (
	"strconv"

	"github.com/liclac/sharlayan/calibre"
)

// Nodes make up the hierarchy of pages to be rendered. Each can be an item and/or a collection.
// A Node generally corresponds to an output directory containing only a rendered 'index.html'.
// Nodes with no Item or Children (eg. the 'series' node if no series are defined) are skipped.
type Node struct {
	Filename string // Required for all but the index page.
	Template string // Required, use '_nav' for a generic list.
	Item     interface{}
	Children []Node
}

// A Link used by generic '_nav' lists.
type Link struct {
	Href string
	Text string
}

// Helper for constructing Nodes whose items are a list of links.
type Dir struct {
	Filename string
	Template string
	Links    []Link
	Children []Node
}

func (d Dir) Node() Node {
	return Node{Filename: d.Filename, Template: d.Template, Item: d.Links, Children: d.Children}
}

// Builds a root node, and a tree based on your Calibre metadata.
func Root(cfg Config, meta *calibre.Metadata) Node {
	var children []Node
	var links []Link

	if len(meta.Books) > 0 {
		books := Dir{Filename: "book", Template: "_nav"}
		for _, book := range meta.Books {
			id := strconv.Itoa(book.ID)
			books.Children = append(books.Children, Node{Filename: id, Template: "book", Item: book})
			books.Links = append(books.Links, Link{Href: id, Text: book.Title})
		}
		children = append(children, books.Node())
		links = append(links, Link{Href: books.Filename, Text: "Books"})
	}

	if !cfg.Author.NoIndex && len(meta.Authors) > 0 {
		authors := Dir{Filename: "author", Template: "_nav"}
		for _, v := range meta.Authors {
			id := strconv.Itoa(v.ID)
			authors.Children = append(authors.Children, Node{Filename: id, Template: "author", Item: v})
			authors.Links = append(authors.Links, Link{Href: id, Text: v.Name})
		}
		children = append(children, authors.Node())
		links = append(links, Link{Href: authors.Filename, Text: "Authors"})
	}
	if !cfg.Series.NoIndex && len(meta.Series) > 0 {
		series := Dir{Filename: "series", Template: "_nav"}
		for _, v := range meta.Series {
			id := strconv.Itoa(v.ID)
			series.Children = append(series.Children, Node{Filename: id, Template: "series", Item: v})
			series.Links = append(series.Links, Link{Href: id, Text: v.Name})
		}
		children = append(children, series.Node())
		links = append(links, Link{Href: series.Filename, Text: "Series"})
	}
	if !cfg.Tag.NoIndex && len(meta.Tags) > 0 {
		tags := Dir{Filename: "tags", Template: "_nav"}
		for _, v := range meta.Tags {
			id := strconv.Itoa(v.ID)
			tags.Children = append(tags.Children, Node{Filename: id, Template: "tag", Item: v})
			tags.Links = append(tags.Links, Link{Href: id, Text: v.Name})
		}
		children = append(children, tags.Node())
		links = append(links, Link{Href: tags.Filename, Text: "Tags"})
	}

	return Node{Template: "_nav", Item: links, Children: children}
}

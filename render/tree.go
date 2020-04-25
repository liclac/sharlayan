package render

import (
	"strconv"

	"github.com/liclac/sharlayan/calibre"
)

// Nodes make up the hierarchy of pages to be rendered. Each can be an item and/or a collection.
// A Node generally corresponds to an output directory containing only a rendered 'index.html'.
// Nodes with no Item or Items (eg. the 'series' node if no series are defined) are skipped.
type Node struct {
	Filename string // Required for all but the index page.
	Template string // Required, defaults to '_nav' for collections.

	// If set, this is an Item node, and will be rendered with this as the template context.
	// An Item node can still have children. If you'd like a listing of them, either include
	// '_nav/list' or have a look at that template for inspiration.
	Item interface{}

	// Child nodes to be rendered under the node.
	Items []Node
}

// Builds a root node, and a tree based on your Calibre metadata.
func Root(cfg Config, meta *calibre.Metadata) Node {
	books := Node{Filename: "book", Items: make([]Node, len(meta.Books))}
	for i, book := range meta.Books {
		books.Items[i] = Node{Filename: strconv.Itoa(book.ID), Template: "book", Item: book}
	}

	authors := Node{Filename: "author"}
	if !cfg.Author.NoIndex {
		authors.Items = make([]Node, len(meta.Authors))
		for i, v := range meta.Authors {
			authors.Items[i] = Node{Filename: strconv.Itoa(v.ID), Template: "author", Item: v}
		}
	}
	series := Node{Filename: "series"}
	if !cfg.Series.NoIndex {
		series.Items = make([]Node, len(meta.Series))
		for i, v := range meta.Series {
			series.Items[i] = Node{Filename: strconv.Itoa(v.ID), Template: "series", Item: v}
		}
	}
	tags := Node{Filename: "tags"}
	if !cfg.Tag.NoIndex {
		tags.Items = make([]Node, len(meta.Tags))
		for i, v := range meta.Tags {
			tags.Items[i] = Node{Filename: strconv.Itoa(v.ID), Template: "tag", Item: v}
		}
	}

	return Node{Template: "index", Items: []Node{books, authors, series, tags}}
}

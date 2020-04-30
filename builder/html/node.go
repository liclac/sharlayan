package html

import (
	"github.com/liclac/sharlayan/builder/tree"
	"github.com/spf13/afero"
)

var _ tree.Node = Page{}
var _ tree.Node = NavPage{}

// Nodes make up the hierarchy of pages to be rendered. Each can be an item and/or a collection.
// A Node generally corresponds to an output directory containing only a rendered 'index.html'.
// Nodes with no Item or Children (eg. the 'series' node if no series are defined) are skipped.
type Page struct {
	Builder *Builder
	tree.NodeInfo
	Template string // Required, use '_nav' for a generic list.
	Item     interface{}
}

func (p Page) Info() tree.NodeInfo { return p.NodeInfo }

func (p Page) Render(fs afero.Fs, path string) error {
	return p.Builder.Render(fs, path, p.Template, p.Item)
}

// A Link used by generic '_nav' lists.
type Link struct {
	Href string
	Text string
}

// Convenience wrapper for Page for building "_nav" pages that wrap a slice of Links.
type NavPage struct {
	Builder *Builder
	tree.NodeInfo
	Links []Link
}

func (p NavPage) Info() tree.NodeInfo { return p.NodeInfo }

func (p NavPage) Render(fs afero.Fs, path string) error {
	return p.Builder.Render(fs, path, "_nav", p.Links)
}

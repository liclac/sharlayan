package render

import (
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

// Returns a root node built from your Calibre metadata.
func Root(cfg Config, meta *calibre.Metadata) Node {
	return Node{Template: "index"}
}

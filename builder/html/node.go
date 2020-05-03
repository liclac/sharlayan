package html

import (
	"github.com/spf13/afero"

	"github.com/liclac/sharlayan/builder/tree"
)

var _ tree.Node = PageNode{}

// An HTML page to be rendered.
type PageNode struct {
	tree.NodeInfo
	Builder  *Builder
	Template string // Required, use '_nav' for a generic list.
	Item     interface{}
}

func Page(b *Builder, id, tmpl string, item interface{}) *PageNode {
	return &PageNode{
		NodeInfo: tree.NodeInfo{ID: id},
		Builder:  b,
		Template: tmpl,
		Item:     item,
	}
}

func (p PageNode) Info() tree.NodeInfo { return p.NodeInfo }

func (p PageNode) Render(fs afero.Fs, ns tree.NamingScheme, path string) error {
	return p.Builder.Render(fs, ns, path, p.Template, p.Item)
}

func Index(b *Builder, nodes ...tree.Node) *PageNode {
	infos := make([]tree.NodeInfo, 0, len(nodes))
	for _, node := range nodes {
		if node != nil {
			infos = append(infos, node.Info())
		}
	}
	return Page(b, "index.html", "_nav", infos)
}
func AddIndex(b *Builder, nodes ...tree.Node) []tree.Node {
	if len(nodes) != 0 {
		nodes = append(nodes, Index(b, nodes...))
	}
	return nodes
}

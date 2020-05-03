package builder

import (
	"github.com/liclac/sharlayan/builder/tree"
	"github.com/liclac/sharlayan/calibre"
)

func ShadowRoot(b *Builder, meta *calibre.Metadata, realPath string, realNS tree.NamingScheme) *tree.DirNode {
	return tree.Dir("", "")
}

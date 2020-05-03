package builder

import (
	"github.com/liclac/sharlayan/builder/html"
	"github.com/liclac/sharlayan/builder/tree"
	"github.com/liclac/sharlayan/calibre"
)

func Root(b *Builder, meta *calibre.Metadata) *tree.DirNode {
	return tree.Dir("", "", html.AddIndex(b.HTML,
		Books(b, meta.Books),
	)...)
}

func Books(b *Builder, books []*calibre.Book) tree.Node {
	return tree.DirInfo(tree.BooksDirInfo, html.AddIndex(b.HTML, BookNodes(b, books)...)...)
}

func BookNodes(b *Builder, books []*calibre.Book) []tree.Node {
	nodes := make([]tree.Node, len(books))
	for i, book := range books {
		nodes[i] = BookNode(b, book)
	}
	return nodes
}

func BookNode(b *Builder, book *calibre.Book) tree.Node {
	// TODO: Copy/Link data files.
	return tree.DirInfo(tree.BookInfo(book), html.Page(b.HTML, "index.html", "book", book))
}

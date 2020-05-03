package builder

import (
	"github.com/liclac/sharlayan/builder/html"
	"github.com/liclac/sharlayan/builder/tree"
	"github.com/liclac/sharlayan/calibre"
)

func Root(b *Builder, meta *calibre.Metadata) *tree.DirNode {
	return tree.Dir("", "", html.AddIndex(b.HTML,
		BookDir(b, meta.Books),
		AuthorDir(b, meta.Authors),
		SeriesDir(b, meta.Series),
		TagDir(b, meta.Tags),
	)...)
}

func BookDir(b *Builder, books []*calibre.Book) tree.Node {
	nodes := make([]tree.Node, len(books))
	for i, book := range books {
		nodes[i] = BookNode(b, book)
	}
	return tree.DirInfo(tree.BookDirInfo, html.AddIndex(b.HTML, nodes...)...)
}

func BookNode(b *Builder, book *calibre.Book) tree.Node {
	// TODO: Copy/Link data files.
	return tree.DirInfo(tree.BookInfo(book), html.Page(b.HTML, "index.html", "book", book))
}

func AuthorDir(b *Builder, authors []*calibre.Author) tree.Node {
	nodes := make([]tree.Node, len(authors))
	for i, author := range authors {
		nodes[i] = AuthorNode(b, author)
	}
	return tree.DirInfo(tree.AuthorDirInfo, html.AddIndex(b.HTML, nodes...)...)
}

func AuthorNode(b *Builder, author *calibre.Author) tree.Node {
	return tree.DirInfo(tree.AuthorInfo(author), html.Page(b.HTML, "index.html", "author", author))
}

func SeriesDir(b *Builder, series []*calibre.Series) tree.Node {
	nodes := make([]tree.Node, len(series))
	for i, series := range series {
		nodes[i] = SeriesNode(b, series)
	}
	return tree.DirInfo(tree.SeriesDirInfo, html.AddIndex(b.HTML, nodes...)...)
}

func SeriesNode(b *Builder, series *calibre.Series) tree.Node {
	return tree.DirInfo(tree.SeriesInfo(series), html.Page(b.HTML, "index.html", "series", series))
}

func TagDir(b *Builder, tags []*calibre.Tag) tree.Node {
	nodes := make([]tree.Node, len(tags))
	for i, tag := range tags {
		nodes[i] = TagNode(b, tag)
	}
	return tree.DirInfo(tree.TagDirInfo, html.AddIndex(b.HTML, nodes...)...)
}

func TagNode(b *Builder, tag *calibre.Tag) tree.Node {
	return tree.DirInfo(tree.TagInfo(tag), html.Page(b.HTML, "index.html", "tag", tag))
}

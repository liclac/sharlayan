package builder

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/liclac/sharlayan/builder/html"
	"github.com/liclac/sharlayan/builder/tree"
	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/config"
	"github.com/spf13/afero"
)

type Builder struct {
	Cfg  config.Config
	HTML *html.Builder
}

func New(cfg config.Config) (*Builder, error) {
	htmlBuilder, err := html.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Builder{
		Cfg:  cfg,
		HTML: htmlBuilder,
	}, nil
}

func (b *Builder) Build(fs afero.Fs, nodes []tree.Node) error {
	for _, node := range nodes {
		nodeInfo := node.Info()
		if err := node.Render(fs, nodeInfo.Path); err != nil {
			return fmt.Errorf("couldn't render %s: %w", nodeInfo.Path, err)
		}
	}
	return nil
}

// Build a tree of tree.Nodes based on your Calibre metadata.
func (b *Builder) Nodes(meta *calibre.Metadata) (nodes []tree.Node) {
	nodes = append(nodes, b.BookNodes(meta)...)
	nodes = append(nodes, b.AuthorNodes(meta)...)
	nodes = append(nodes, b.IndexNodes(meta)...)
	return nodes
}

func (b *Builder) BookNodes(meta *calibre.Metadata) (nodes []tree.Node) {
	path := b.Cfg.Books.Path
	index := html.NavPage{
		Builder:  b.HTML,
		NodeInfo: tree.NodeInfo{Path: filepath.Join(path, "index.html")},
	}
	for _, book := range meta.Books {
		id := strconv.Itoa(book.ID)
		nodes = append(nodes, html.Page{
			Builder:  b.HTML,
			NodeInfo: tree.NodeInfo{Path: filepath.Join(path, id, "index.html")},
			Template: "book",
			Item:     book,
		})
		index.Links = append(index.Links,
			html.Link{Href: filepath.Join(path, id), Text: book.Title})
	}
	return append(nodes, index)
}

func (b *Builder) AuthorNodes(meta *calibre.Metadata) (nodes []tree.Node) {
	path := b.Cfg.Authors.Path
	index := html.NavPage{
		Builder:  b.HTML,
		NodeInfo: tree.NodeInfo{Path: filepath.Join(path, "index.html")},
	}
	for _, author := range meta.Authors {
		id := strconv.Itoa(author.ID)
		nodes = append(nodes, html.Page{
			Builder:  b.HTML,
			NodeInfo: tree.NodeInfo{Path: filepath.Join(path, id, "index.html")},
			Template: "author",
			Item:     author,
		})
		index.Links = append(index.Links,
			html.Link{Href: filepath.Join(path, id), Text: author.Name})
	}
	return append(nodes, index)
}

func (b *Builder) SeriesNodes(meta *calibre.Metadata) (nodes []tree.Node) {
	path := b.Cfg.Series.Path
	index := html.NavPage{
		Builder:  b.HTML,
		NodeInfo: tree.NodeInfo{Path: filepath.Join(path, "index.html")},
	}
	for _, series := range meta.Series {
		id := strconv.Itoa(series.ID)
		nodes = append(nodes, html.Page{
			Builder:  b.HTML,
			NodeInfo: tree.NodeInfo{Path: filepath.Join(path, id, "index.html")},
			Template: "series",
			Item:     series,
		})
		index.Links = append(index.Links,
			html.Link{Href: filepath.Join(path, id), Text: series.Name})
	}
	return append(nodes, index)
}

func (b *Builder) TagNodes(meta *calibre.Metadata) (nodes []tree.Node) {
	path := b.Cfg.Tags.Path
	index := html.NavPage{
		Builder:  b.HTML,
		NodeInfo: tree.NodeInfo{Path: filepath.Join(path, "index.html")},
	}
	for _, tag := range meta.Tags {
		id := strconv.Itoa(tag.ID)
		nodes = append(nodes, html.Page{
			Builder:  b.HTML,
			NodeInfo: tree.NodeInfo{Path: filepath.Join(path, id, "index.html")},
			Template: "tag",
			Item:     tag,
		})
		index.Links = append(index.Links,
			html.Link{Href: filepath.Join(path, id), Text: tag.Name})
	}
	return append(nodes, index)
}

func (b *Builder) IndexNodes(meta *calibre.Metadata) []tree.Node {
	return []tree.Node{
		html.NavPage{
			Builder:  b.HTML,
			NodeInfo: tree.NodeInfo{Path: "/index.html"},
			Links: []html.Link{
				html.Link{Href: b.Cfg.Books.Path, Text: "Books"},
				html.Link{Href: b.Cfg.Authors.Path, Text: "Authors"},
				html.Link{Href: b.Cfg.Series.Path, Text: "Series"},
				html.Link{Href: b.Cfg.Tags.Path, Text: "Tags"},
			},
		},
	}
}

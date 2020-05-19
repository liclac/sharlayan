package tree

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

// A NodeWrapper is a Node + spatial awareness, used by the Tree.
type NodeWrapper struct {
	Node
	NodeInfo
	Path     string         // Full path, relative to Tree.Path.
	Children []*NodeWrapper // Child nodes, if any.
}

func WrapNode(basePath string, n Node) *NodeWrapper {
	if n == nil {
		return nil
	}

	w := &NodeWrapper{Node: n, NodeInfo: n.Info()}
	w.Path = filepath.Join(basePath, w.Filename)

	if n, ok := n.(NodeCollection); ok {
		children := n.Children()
		w.Children = make([]*NodeWrapper, len(children))
		for i, child := range children {
			w.Children[i] = WrapNode(w.Path, child)
		}
	}

	return w
}

// A Tree is an abstract representation of a filesystem hierarchy.
type Tree struct {
	billy.Filesystem
	Path string       // Path in the surrounding filesystem.
	Root *NodeWrapper // The root node.

	ByPath   map[string]*NodeWrapper
	ByLinkID map[string]*NodeWrapper
}

// Creates a new Tree, for rendering the root node to a given path in the passed filesystem.
func New(fs billy.Filesystem, path string, root Node) (*Tree, error) {
	t := &Tree{
		Filesystem: fs,
		Path:       path,
		Root:       WrapNode("/", root),
		ByPath:     map[string]*NodeWrapper{},
		ByLinkID:   map[string]*NodeWrapper{},
	}
	return t, t.walkNode(t.Root)
}

// Shorthand for creating a Tree and calling Render() on it.
func Render(fs billy.Filesystem, path string, root Node) error {
	t, err := New(fs, path, root)
	if err != nil {
		return err
	}
	return t.Render()
}

// Recursive function for walking a NodeWrapper when initialising a Tree.
func (t *Tree) walkNode(w *NodeWrapper) error {
	if w == nil {
		return nil
	}

	t.ByPath[w.Path] = w

	if w.LinkID != "" {
		if old, hit := t.ByLinkID[w.LinkID]; hit {
			return fmt.Errorf("duplicate LinkID '%s' for %s, previously used by: %s",
				w.LinkID, w.Path, old.Path)
		}
		t.ByLinkID[w.LinkID] = w
	}

	for i, cw := range w.Children {
		if cw == nil {
			continue
		}
		if cw.Filename == "" {
			return fmt.Errorf("child %d has no Filename: %s", i, w.Path)
		}
		if err := t.walkNode(cw); err != nil {
			return err
		}
	}

	return nil
}

func (t Tree) Render() error {
	if err := t.MkdirAll(t.Path, 0755); err != nil {
		return err
	}
	return t.renderNode(t.Root)
}

func (t Tree) renderNode(w *NodeWrapper) error {
	if w == nil {
		return nil
	}
	if err := w.Render(t, filepath.Join(t.Path, w.Path)); err != nil {
		return fmt.Errorf("rendering %s: %w", w.Path, err)
	}
	for _, cw := range w.Children {
		if err := t.renderNode(cw); err != nil {
			return err
		}
	}
	return nil
}

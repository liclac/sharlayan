package tree

import (
	"os"

	"github.com/go-git/go-billy/v5/util"
)

var _ NodeCollection = DirNode{}

// A Node is the basic building block of a Tree.
type Node interface {
	Info() NodeInfo
	Render(t Tree, path string) error
}

// A NodeCollection is a node that contains other nodes, eg. a directory.
// These nodes are expected to create their directories in Render().
type NodeCollection interface {
	Node
	Children() []Node
}

// NodeInfo represents basic information about a node.
type NodeInfo struct {
	Filename string // Required: Filename relative to parent node.

	// Optional: Unique ID for the node. The LinkID can be resolved to a path in a Tree.
	// This is used for two things: declaring symlinks that are resolved at render time,
	// and create automatic symlinks to deduplicate two merged Trees.
	LinkID string
}

// A DirNode represents a directory of other nodes. Implements NodeCollection.
type DirNode struct {
	NodeInfo
	Nodes []Node
}

// Returns a DirNode containing the given nodes.
func Dir(filename string, nodes ...Node) DirNode {
	return DirID(filename, "", nodes...)
}

// Returns a DirNode with a LinkID, containing the given nodes.
func DirID(filename, linkID string, nodes ...Node) DirNode {
	return DirNode{NodeInfo{Filename: filename, LinkID: linkID}, nodes}
}

func (n DirNode) Info() NodeInfo   { return n.NodeInfo }
func (n DirNode) Children() []Node { return n.Nodes }

func (n DirNode) Render(t Tree, path string) error {
	return t.MkdirAll(path, 0755)
}

// A ConstNode writes some static data to a file.
type ConstNode struct {
	NodeInfo
	Mode os.FileMode
	Data []byte
}

// Returns a ConstNode that writes a string to a file.
func String(filename string, s string) ConstNode {
	return ConstNode{NodeInfo{Filename: filename}, 0644, []byte(s)}
}

func (n ConstNode) Info() NodeInfo { return n.NodeInfo }

func (n ConstNode) Render(t Tree, path string) error {
	return util.WriteFile(t, path, n.Data, n.Mode)
}

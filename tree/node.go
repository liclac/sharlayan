package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5/util"
)

// A Node is the basic building block of a Tree.
type Node interface {
	Info() NodeInfo
	Render(t Tree, self *NodeWrapper, path string) error
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
	Alias    string // Optional: Display name used in listings.

	// Optional: Unique ID for the node. The LinkID can be resolved to a path in a Tree.
	// This is used for two things: declaring symlinks that are resolved at render time,
	// and create automatic symlinks to deduplicate two merged Trees.
	LinkID string
}

type NodeInfoOpt func(info *NodeInfo)

// Functional builder for NodeInfo structs, see NodeInfoOpt.
func Named(filename string, opts ...NodeInfoOpt) NodeInfo {
	info := NodeInfo{Filename: filename}
	for _, opt := range opts {
		opt(&info)
	}
	return info
}

// NodeInfoOpt that sets NodeInfo.Alias.
func Alias(alias string) NodeInfoOpt {
	return func(info *NodeInfo) { info.Alias = alias }
}

// NodeInfoOpt that sets NodeInfo.LinkID.
func LinkID(id ...string) NodeInfoOpt {
	return func(info *NodeInfo) { info.LinkID = strings.Join(id, ":") }
}

// Automatically implement Node.Info() on node types that embed a NodeInfo.
func (i NodeInfo) Info() NodeInfo { return i }

var _ NodeCollection = DirNode{}

// A DirNode represents a directory of other nodes. Implements NodeCollection.
type DirNode struct {
	NodeInfo
	Nodes []Node
}

// Returns a DirNode containing the given nodes.
func Dir(info NodeInfo, nodes ...Node) DirNode {
	return DirNode{info, nodes}
}

func (n DirNode) Children() []Node { return n.Nodes }

func (n DirNode) Render(t Tree, self *NodeWrapper, path string) error {
	return t.MkdirAll(path, 0755)
}

var _ Node = ConstNode{}

// A ConstNode writes some static data to a file.
type ConstNode struct {
	NodeInfo
	Mode os.FileMode
	Data []byte
}

// Returns a ConstNode that writes a string to a file.
func String(info NodeInfo, s string) ConstNode {
	return ConstNode{info, 0644, []byte(s)}
}

func (n ConstNode) Render(t Tree, self *NodeWrapper, path string) error {
	return util.WriteFile(t, path, n.Data, n.Mode)
}

var _ Node = SymlinkNode{}

// A SymlinkNode creates a symlink to a target identified by its LinkID.
// The resulting link will use the shortest possible relative path.
type SymlinkNode struct {
	NodeInfo
	TargetLinkID string
}

// Returns a SymlinkNode that links to the given LinkID.
func Symlink(info NodeInfo, targetLinkID ...string) SymlinkNode {
	return SymlinkNode{info, strings.Join(targetLinkID, ":")}
}

func (n SymlinkNode) Render(t Tree, self *NodeWrapper, path string) error {
	dst, dstOK := t.ByLinkID[n.TargetLinkID]
	if !dstOK {
		return fmt.Errorf("unknown link ID: %s", n.TargetLinkID)
	}
	relPath, err := filepath.Rel(filepath.Dir(self.Path), dst.Path)
	if err != nil {
		return fmt.Errorf("couldn't find relative path from %s [link] to %s [linkID='%s']: %w",
			self.Path, dst.Path, n.TargetLinkID, err)
	}
	return t.Symlink(relPath, path)
}

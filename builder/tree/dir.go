package tree

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"go.uber.org/zap"
)

var _ Node = DirNode{}

// A DirNode is a node that can contain other nodes.
// Rendering a DirNode with no children is a no-op - it won't create empty directories.
type DirNode struct {
	NodeInfo
	Nodes []Node
}

// Create a DirNode containing the given nodes, skipping any nils.
func Dir(id, name string, nodes ...Node) *DirNode {
	return DirInfo(NodeInfo{ID: id, Name: name}, nodes...)
}

func DirInfo(info NodeInfo, nodes ...Node) *DirNode {
	nonNilNodes := make([]Node, 0, len(nodes))
	for _, node := range nodes {
		nonNilNodes = append(nonNilNodes, node)
	}
	return &DirNode{info, nonNilNodes}
}

func (dir DirNode) Info() NodeInfo { return dir.NodeInfo }

func (dir DirNode) Render(fs afero.Fs, ns NamingScheme, path string) error {
	L := zap.L().With(zap.String("path", path))
	if len(dir.Nodes) == 0 {
		L.Warn("Skipping empty directory")
		return nil
	}
	if err := fs.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("couldn't mkdir: %w", err)
	}
	for _, node := range dir.Nodes {
		info := node.Info()
		filename := info.Filename(ns)
		L.Debug("Dir: Rendering...", zap.String("filename", filename))
		if err := node.Render(fs, ns, filepath.Join(path, filename)); err != nil {
			return fmt.Errorf("Dir(%s): %w", filename, err)
		}
	}
	return nil
}

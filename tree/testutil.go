package tree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/stretchr/testify/require"
)

type fsDumpNode struct {
	Filename string
	Mode     os.FileMode
	Data     []byte
	Nodes    []fsDumpNode
}

func fsDump(t *testing.T, fs billy.Filesystem, path string) fsDumpNode {
	node, err := fsDump_(fs, path)
	require.NoError(t, err)
	return node
}

// Dumps a filesystem into a tree of fsDumpNodes.
func fsDump_(fs billy.Filesystem, path string) (fsDumpNode, error) {
	stat, err := fs.Stat(path)
	if err != nil {
		return fsDumpNode{}, fmt.Errorf("stat: %s: %w", path, err)
	}
	node := fsDumpNode{
		Filename: stat.Name(),
		Mode:     stat.Mode(),
	}

	if stat.IsDir() {
		childInfos, err := fs.ReadDir(path)
		if err != nil {
			return node, fmt.Errorf("readdir: %s: %w", path, err)
		}
		for _, info := range childInfos {
			child, err := fsDump_(fs, filepath.Join(path, info.Name()))
			if err != nil {
				return node, err
			}
			node.Nodes = append(node.Nodes, child)
		}
		sort.Slice(node.Nodes, func(i, j int) bool {
			return node.Nodes[i].Filename < node.Nodes[j].Filename
		})
	} else {
		f, err := fs.Open(path)
		if err != nil {
			return node, fmt.Errorf("open: %s: %w", path, err)
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return node, fmt.Errorf("read: %s: %w", path, err)
		}
		node.Data = data
	}

	return node, nil
}

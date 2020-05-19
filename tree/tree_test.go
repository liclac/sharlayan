package tree

import (
	"os"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Build a tree consisting only of a single, non-directory node.
func TestTreeSingleNode(t *testing.T) {
	fs := memfs.New()
	root := WrapNode("/", String("test.txt", "abc123"))

	tree, err := New(fs, "/prefix", root.Node)
	require.NoError(t, err)
	assert.Equal(t, &Tree{
		Filesystem: fs,
		Path:       "/prefix",
		Root:       root,
		ByPath: map[string]*NodeWrapper{
			"/test.txt": root,
		},
		ByLinkID: map[string]*NodeWrapper{},
	}, tree)

	require.NoError(t, tree.Render())
	assert.Equal(t, fsDumpNode{Filename: "/", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
		fsDumpNode{Filename: "prefix", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
			fsDumpNode{Filename: "test.txt", Mode: 0644, Data: []byte("abc123")},
		}},
	}}, fsDump(t, fs, "/"))
}

// Building a tree with a nil root node should render an empty directory.
func TestTreeNilRoot(t *testing.T) {
	fs := memfs.New()
	tree, err := New(fs, "/prefix", nil)
	require.NoError(t, err)
	assert.Equal(t, &Tree{
		Filesystem: fs,
		Path:       "/prefix",
		ByPath:     map[string]*NodeWrapper{},
		ByLinkID:   map[string]*NodeWrapper{},
	}, tree)
	require.NoError(t, tree.Render())
	assert.Equal(t, fsDumpNode{Filename: "/", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
		fsDumpNode{Filename: "prefix", Mode: 0755 | os.ModeDir},
	}}, fsDump(t, fs, "/"))
}

// Build a tree that requires correctly recursing.
func TestTreeDirs(t *testing.T) {
	fs := memfs.New()
	root := WrapNode("/",
		Dir("",
			String("test1.txt", "abc123"),
			Dir("sub",
				String("test2.txt", "def456"))))

	tree, err := New(fs, "/prefix", root.Node)
	require.NoError(t, err)
	assert.Equal(t, &Tree{
		Filesystem: fs,
		Path:       "/prefix",
		Root:       root,
		ByPath: map[string]*NodeWrapper{
			"/":              root,
			"/test1.txt":     root.Children[0],
			"/sub":           root.Children[1],
			"/sub/test2.txt": root.Children[1].Children[0],
		},
		ByLinkID: map[string]*NodeWrapper{},
	}, tree)

	require.NoError(t, tree.Render())
	assert.Equal(t, fsDumpNode{Filename: "/", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
		fsDumpNode{Filename: "prefix", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
			fsDumpNode{Filename: "sub", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
				fsDumpNode{Filename: "test2.txt", Mode: 0644, Data: []byte("def456")},
			}},
			fsDumpNode{Filename: "test1.txt", Mode: 0644, Data: []byte("abc123")},
		}},
	}}, fsDump(t, fs, "/"))
}

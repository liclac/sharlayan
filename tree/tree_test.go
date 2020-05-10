package tree

import (
	"os"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	t.Run("File Root", func(t *testing.T) {
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
		}}, fsDumpT(t, fs, "/"))
	})
}

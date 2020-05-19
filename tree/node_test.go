package tree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirNode(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, fsDumpNode{Filename: "dir", Mode: 0755 | os.ModeDir},
			renderDump(t, Dir(Named("dir"))))
	})

	t.Run("One File", func(t *testing.T) {
		assert.Equal(t,
			fsDumpNode{Filename: "dir", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
				fsDumpNode{Filename: "test.txt", Mode: 0644, Data: "abc123"},
			}}, renderDump(t, Dir(Named("dir"), String(Named("test.txt"), "abc123"))))

		t.Run("nil", func(t *testing.T) {
			assert.Equal(t, fsDumpNode{Filename: "dir", Mode: 0755 | os.ModeDir},
				renderDump(t, Dir(Named("dir"), nil)))
		})
	})

	t.Run("Two Files", func(t *testing.T) {
		assert.Equal(t,
			fsDumpNode{Filename: "dir", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
				fsDumpNode{Filename: "test1.txt", Mode: 0644, Data: "abc123"},
				fsDumpNode{Filename: "test2.txt", Mode: 0644, Data: "def456"}}},
			renderDump(t, Dir(Named("dir"),
				String(Named("test1.txt"), "abc123"),
				String(Named("test2.txt"), "def456"))))
	})
}

func TestConstNode(t *testing.T) {
	assert.Equal(t,
		fsDumpNode{Filename: "test.txt", Mode: 0644, Data: "abc123"},
		renderDump(t, String(Named("test.txt"), "abc123")))
}

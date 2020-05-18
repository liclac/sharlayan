package tree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that a dir node creates a directory.
func TestDirNode(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, fsDumpNode{Filename: "dir", Mode: 0755 | os.ModeDir},
			renderDump(t, Dir("dir")))
	})

	t.Run("One File", func(t *testing.T) {
		assert.Equal(t,
			fsDumpNode{Filename: "dir", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
				fsDumpNode{Filename: "test.txt", Mode: 0644, Data: []byte("abc123")},
			}}, renderDump(t, Dir("dir", String("test.txt", "abc123"))))
	})

	t.Run("Two Files", func(t *testing.T) {
		assert.Equal(t,
			fsDumpNode{Filename: "dir", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
				fsDumpNode{Filename: "test1.txt", Mode: 0644, Data: []byte("abc123")},
				fsDumpNode{Filename: "test2.txt", Mode: 0644, Data: []byte("def456")},
			}},
			renderDump(t, Dir("dir", String("test1.txt", "abc123"), String("test2.txt", "def456"))))
	})
}

// Test that a const node writes a file.
func TestConstNode(t *testing.T) {
	assert.Equal(t, fsDumpNode{Filename: "test.txt", Mode: 0644, Data: []byte("abc123")},
		renderDump(t, String("test.txt", "abc123")))
}

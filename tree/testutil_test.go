package tree

import (
	"os"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_fsDump(t *testing.T) {
	fs := memfs.New()
	require.NoError(t, fs.MkdirAll("/a/b", 0755))
	require.NoError(t, util.WriteFile(fs, "/a/b/file1.txt", []byte("test 1"), 0644))
	require.NoError(t, util.WriteFile(fs, "/a/b/file2.txt", []byte("test 2"), 0644))
	require.NoError(t, fs.MkdirAll("/a/c", 0755))
	require.NoError(t, util.WriteFile(fs, "/a/c/file3.txt", []byte("test 3"), 0644))
	require.NoError(t, fs.MkdirAll("/a/b/d", 0755))
	require.NoError(t, util.WriteFile(fs, "/a/b/d/file4.txt", []byte("test 4"), 0644))
	require.NoError(t, fs.MkdirAll("/a/e", 0755))

	assert.Equal(t, fsDumpNode{
		Filename: "/",
		Mode:     0755 | os.ModeDir,
		Nodes: []fsDumpNode{
			fsDumpNode{Filename: "a", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
				fsDumpNode{Filename: "b", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
					fsDumpNode{Filename: "d", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
						fsDumpNode{Filename: "file4.txt", Mode: 0644, Data: []byte("test 4")},
					}},
					fsDumpNode{Filename: "file1.txt", Mode: 0644, Data: []byte("test 1")},
					fsDumpNode{Filename: "file2.txt", Mode: 0644, Data: []byte("test 2")},
				}},
				fsDumpNode{Filename: "c", Mode: 0755 | os.ModeDir, Nodes: []fsDumpNode{
					fsDumpNode{Filename: "file3.txt", Mode: 0644, Data: []byte("test 3")},
				}},
				fsDumpNode{Filename: "e", Mode: 0755 | os.ModeDir},
			}},
		},
	}, fsDumpT(t, fs, "/"))
}

package ssh

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (i testFileInfo) Name() string       { return i.name }
func (i testFileInfo) Size() int64        { return i.size }
func (i testFileInfo) Mode() os.FileMode  { return i.mode }
func (i testFileInfo) ModTime() time.Time { return i.modTime }
func (i testFileInfo) IsDir() bool        { return i.isDir }
func (i testFileInfo) Sys() interface{}   { return i }

func Test_sftpLister(t *testing.T) {
	testdata := map[string]struct {
		outLen int
		offset int64
		n      int
		EOF    bool
		out    []os.FileInfo
		lister sftpLister
	}{
		"Empty":        {0, 0, 0, true, sftpLister{}, []os.FileInfo{}},
		"Empty/Offset": {0, 100, 0, true, sftpLister{}, []os.FileInfo{}},

		"Exact": {2, 0, 2, true,
			[]os.FileInfo{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}},
			sftpLister{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}}},
		"Exact/Offset": {1, 1, 1, true,
			[]os.FileInfo{testFileInfo{name: "test2.txt"}},
			sftpLister{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}}},

		"Larger": {9, 0, 2, true,
			[]os.FileInfo{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}},
			sftpLister{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}}},
		"Larger/Offset": {9, 1, 1, true,
			[]os.FileInfo{testFileInfo{name: "test2.txt"}},
			sftpLister{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}}},
		"Larger/Offset/PastEnd": {9, 2, 0, true,
			[]os.FileInfo{},
			sftpLister{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}}},

		"Smaller": {1, 0, 1, false,
			[]os.FileInfo{testFileInfo{name: "test1.txt"}},
			sftpLister{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}}},
		"Smaller/End": {1, 1, 1, true,
			[]os.FileInfo{testFileInfo{name: "test2.txt"}},
			sftpLister{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}}},
		"Smaller/End/Past": {1, 2, 0, true,
			[]os.FileInfo{},
			sftpLister{testFileInfo{name: "test1.txt"}, testFileInfo{name: "test2.txt"}}},
	}
	for name, tdata := range testdata {
		t.Run(name, func(t *testing.T) {
			out := make([]os.FileInfo, tdata.outLen)
			n, err := tdata.lister.ListAt(out, tdata.offset)
			assert.Equal(t, n, tdata.n)
			if tdata.EOF {
				assert.Equal(t, io.EOF, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package tree

import (
	"github.com/spf13/afero"
)

type NodeInfo struct {
	Path string
}

type Node interface {
	Info() NodeInfo
	Render(fs afero.Fs, path string) error
}

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"

	"github.com/liclac/sharlayan/builder"
	"github.com/liclac/sharlayan/builder/tree"
	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/config"
)

func dump(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func loadMeta(cfg config.Config) (*calibre.Metadata, error) {
	return calibre.Read(cfg.Library)
}

func createBuilder(cfg config.Config) (*builder.Builder, error) {
	return builder.New(cfg)
}

func createNodes(bld *builder.Builder, meta *calibre.Metadata) []tree.Node {
	return bld.Nodes(meta)
}

func buildToFs(bld *builder.Builder, nodes []tree.Node, fs afero.Fs) error {
	return bld.Build(fs, nodes)
}

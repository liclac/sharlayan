package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"

	"github.com/liclac/sharlayan/builder"
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

func buildToFs(cfg config.Config, fs afero.Fs) error {
	meta, err := calibre.Read(cfg.Library)
	if err != nil {
		return err
	}
	bld, err := builder.New(cfg)
	if err != nil {
		return err
	}
	return bld.Build(fs, bld.Nodes(meta))
}

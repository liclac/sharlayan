package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/liclac/sharlayan/afhack"
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

func traceFS(cfg *config.Config, fs afero.Fs) afero.Fs {
	if !cfg.Debug.TraceFS {
		return fs
	}
	L := zap.L().Named("fs")
	L.Debug("FS Tracing enabled")
	return afhack.NewTraceFs(L, fs)
}

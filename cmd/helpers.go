package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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
	metaStart := time.Now()
	meta, err := calibre.Read(cfg.Library)
	metaTime := time.Since(metaStart)
	log.WithFields(log.Fields{
		"books": len(meta.Books),
		"t":     metaTime,
	}).Debug("Loaded: Calibre Metadata")
	return meta, err
}

func createBuilder(cfg config.Config) (*builder.Builder, error) {
	bldStart := time.Now()
	bld, err := builder.New(cfg)
	bldTime := time.Since(bldStart)
	log.WithField("t", bldTime).Debug("Loaded: Builder")
	return bld, err
}

func createNodes(bld *builder.Builder, meta *calibre.Metadata) []tree.Node {
	nodesStart := time.Now()
	nodes := bld.Nodes(meta)
	nodesTime := time.Since(nodesStart)
	log.WithFields(log.Fields{
		"num": len(nodes),
		"t":   nodesTime,
	}).Debug("Loaded: Nodes")
	return nodes
}

func buildToFs(bld *builder.Builder, nodes []tree.Node, fs afero.Fs) error {
	buildStart := time.Now()
	err := bld.Build(fs, nodes)
	buildTime := time.Since(buildStart)
	log.WithField("t", buildTime).Debug("Rendered!")
	return err
}

package render

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/liclac/sharlayan/calibre"
)

type Config struct {
	Title string
}

type TemplateContext struct {
	Cfg  Config
	Meta *calibre.Metadata
}

func Render(outPath, templatePath string, meta *calibre.Metadata, cfg Config) error {
	// Parse the whole template file tree into named templates.
	t := template.New("")
	if err := filepath.Walk(templatePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Find the relative template path, eg. 'templates/book/list.tmpl' -> 'book/list.tmpl'.
			filename, err := filepath.Rel(templatePath, path)
			if err != nil {
				return fmt.Errorf("couldn't get relative template path: %s: %w", path, err)
			}

			// Skip directories and anything that doesn't end in .tmpl.
			if info.IsDir() {
				log.WithField("filename", filename).Debug("template init: directory...")
				return nil
			}
			if !strings.HasSuffix(filename, ".tmpl") {
				log.WithField("filename", filename).Debug("template init: not a template")
				return nil
			}

			// Parse it in, named after the base name without extension.
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("couldn't read template: %s: %w", filename, err)
			}
			name := strings.TrimSuffix(info.Name(), ".tmpl")
			if _, err := t.New(name).Parse(string(data)); err != nil {
				return fmt.Errorf("couldn't parse template: %s: %w", filename, err)
			}
			log.WithFields(log.Fields{
				"filename": filename,
				"name":     name,
			}).Debug("template init: loaded template")
			return err
		},
	); err != nil {
		return err
	}

	// Render!
	tctx := &TemplateContext{Cfg: cfg, Meta: meta}
	return renderToFile(filepath.Join(outPath, "index.html"), t, "index", tctx)
}

func renderToFile(path string, t *template.Template, name string, tctx interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.ExecuteTemplate(f, name, tctx)
}

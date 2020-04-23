package render

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".tmpl") {
				return nil
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("couldn't read template: %s: %w", path, err)
			}
			templateName := strings.TrimSuffix(info.Name(), ".tmpl")
			if _, err := t.New(templateName).Parse(string(data)); err != nil {
				return fmt.Errorf("couldn't parse template: %s: %w", path, err)
			}
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

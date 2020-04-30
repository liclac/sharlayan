package html

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/liclac/sharlayan/calibre"
	"github.com/liclac/sharlayan/config"
)

type Builder struct {
	Config    config.Config
	Meta      *calibre.Metadata
	Base      *template.Template
	Templates map[string]*template.Template
}

func New(cfg config.Config) (*Builder, error) {
	b := &Builder{
		Config:    cfg,
		Base:      template.New(""),
		Templates: make(map[string]*template.Template),
	}
	return b, b.loadTemplates()
}

func (b *Builder) loadTemplates() error {
	b.Base = b.Base.Funcs(NewFuncs(&b.Config).Map())
	names, err := b.listTemplates()
	if err != nil {
		return fmt.Errorf("couldn't list templates: %w", err)
	}

	// The docs do not make it super clear how the semantics work when it comes to bundles
	// of multiple templates, but the gist of it is essentially that you can't load them all
	// into one big template and call them by name - you have to create a "base" bundle of
	// your shared layouts, functions and partials, then use Clone() to create separate ones
	// for your "leaf" templates, which can then access the shared stuff in the base bundle.
	leaves := []string{}
	for _, name := range names {
		if name == "layout" || strings.ContainsRune(name, '/') {
			if err := b.loadTemplate(b.Base, name); err != nil {
				return fmt.Errorf("couldn't load %s: %w", name, err)
			}
		} else {
			leaves = append(leaves, name)
		}
	}
	for _, name := range leaves {
		t, err := b.Base.Clone()
		if err != nil {
			return fmt.Errorf("couldn't clone base template to load %s: %w", name, err)
		}
		if err := b.loadTemplate(t, name); err != nil {
			return fmt.Errorf("couldn't load %s: %w", name, err)
		}
		b.Templates[name] = t
	}
	return nil
}

func (b *Builder) loadTemplate(base *template.Template, name string) error {
	data, err := ioutil.ReadFile(filepath.Join(b.Config.HTML.Templates, name+".tmpl"))
	if err != nil {
		return err
	}
	_, err = base.New(name).Parse(string(data))
	return err
}

func (b *Builder) listTemplates() ([]string, error) {
	var names []string
	err := filepath.Walk(b.Config.HTML.Templates,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			filename, err := filepath.Rel(b.Config.HTML.Templates, path)
			if err != nil {
				return fmt.Errorf("couldn't get relative template path: %s: %w", path, err)
			}
			if !info.IsDir() && strings.HasSuffix(filename, ".tmpl") {
				names = append(names, strings.TrimSuffix(filename, ".tmpl"))
			}
			return nil
		},
	)
	return names, err
}

func (b *Builder) Render(fs afero.Fs, path, name string, tctx interface{}) error {
	if err := fs.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if t, ok := b.Templates[name]; ok {
		return t.ExecuteTemplate(f, name, tctx)
	}
	return fmt.Errorf("no such template: %s", name)
}

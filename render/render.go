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
	// Input.
	Out       string `mapstructure:"out"`
	Templates string `mapstructure:"templates"`

	// Output.
	Title string `mapstructure:"title"`

	Author struct {
		NoIndex bool `mapstructure:"no-index"`
	} `mapstructure:"author"`
	Series struct {
		NoIndex bool `mapstructure:"no-index"`
	} `mapstructure:"series"`
	Tag struct {
		NoIndex bool `mapstructure:"no-index"`
	} `mapstructure:"tag"`
}

type Renderer struct {
	Config    Config
	Meta      *calibre.Metadata
	Base      *template.Template
	Templates map[string]*template.Template
}

func New(cfg Config, meta *calibre.Metadata) (*Renderer, error) {
	r := &Renderer{
		Config:    cfg,
		Meta:      meta,
		Base:      template.New(""),
		Templates: make(map[string]*template.Template),
	}
	return r, r.loadTemplates()
}

func (r *Renderer) loadTemplates() error {
	r.Base = r.Base.Funcs(template.FuncMap{
		"cfg": func() *Config { return &r.Config },
	})
	names, err := r.listTemplates()
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
			if err := r.loadTemplate(r.Base, name); err != nil {
				return fmt.Errorf("couldn't load %s: %w", name, err)
			}
		} else {
			leaves = append(leaves, name)
		}
	}
	for _, name := range leaves {
		t, err := r.Base.Clone()
		if err != nil {
			return fmt.Errorf("couldn't clone base template to load %s: %w", name, err)
		}
		if err := r.loadTemplate(t, name); err != nil {
			return fmt.Errorf("couldn't load %s: %w", name, err)
		}
		r.Templates[name] = t
	}
	return nil
}

func (r *Renderer) loadTemplate(base *template.Template, name string) error {
	data, err := ioutil.ReadFile(filepath.Join(r.Config.Templates, name+".tmpl"))
	if err != nil {
		return err
	}
	_, err = base.New(name).Parse(string(data))
	return err
}

func (r *Renderer) listTemplates() ([]string, error) {
	var names []string
	err := filepath.Walk(r.Config.Templates,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			filename, err := filepath.Rel(r.Config.Templates, path)
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

func (r *Renderer) Render(node Node) (int, error) {
	return r.renderTree(node, nil)
}

func (r *Renderer) renderTree(node Node, parents []string) (int, error) {
	// Ignore empty nodes, except for the root.
	if len(parents) > 0 && node.Item == nil && len(node.Items) == 0 {
		return 0, nil
	}

	// All nodes except for the root must have a Filename.
	parentsAndSelf := append(parents, node.Filename)
	relPath := filepath.Join(parentsAndSelf...)
	if node.Filename == "" && len(parents) > 0 {
		return 0, fmt.Errorf("child has no filename: %s", relPath)
	}

	// Items must have a Template, collections get a default of "_nav".
	tname := node.Template
	if tname == "" && node.Item == nil {
		tname = "_nav"
	}
	if tname == "" {
		return 0, fmt.Errorf("no template set for: %s", relPath)
	}

	// Render children, count how many pages were actually rendered.
	var rendered int
	for _, child := range node.Items {
		childRendered, err := r.renderTree(child, parentsAndSelf)
		if err != nil {
			return rendered, err
		}
		rendered += childRendered
	}

	// If this is an Item, or any children rendered pages, render an index page and count it.
	if node.Item != nil || rendered > 0 {
		indexPath := filepath.Join(r.Config.Out, relPath, "index.html")
		if err := r.render(indexPath, tname, node.Item); err != nil {
			return rendered, fmt.Errorf("error rendering %s (template: '%s'): %w", relPath, tname, err)
		}
		rendered++
	}
	return rendered, nil
}

func (r *Renderer) render(path string, name string, tctx interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if t, ok := r.Templates[name]; ok {
		return t.ExecuteTemplate(f, name, tctx)
	}
	return fmt.Errorf("no such template: %s", name)
}

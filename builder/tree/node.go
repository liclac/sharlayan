package tree

import (
	"github.com/spf13/afero"
)

// There are two supported naming schemes: IDs and Names.
type NamingScheme string

const (
	ByID   NamingScheme = "id"   // eg. "/authors/4"
	ByName NamingScheme = "name" // eg. "/Authors/Terry Pratchett"
)

// Returns the name matching the scheme.
// If no appropriate name is given, defaults to the ID.
func (ns NamingScheme) Name(id, name string) string {
	if ns == ByName && name != "" {
		return name
	}
	return id
}

// Basic information about a node. Only the ID field is required.
type NodeInfo struct {
	ID   string // Filename in the ByID scheme, eg. "authors", "4".
	Name string // Filename in the ByName scheme, eg. "Authors", "Terry Pratchett".
}

// Returns the node's filename in the given naming scheme.
// If ns is ByName, but no Name is specified, ID is used.
func (n NodeInfo) Filename(ns NamingScheme) string {
	return ns.Name(n.ID, n.Name)
}

type Node interface {
	Info() NodeInfo
	Render(fs afero.Fs, ns NamingScheme, path string) error
}

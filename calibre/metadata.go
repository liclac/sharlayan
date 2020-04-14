package calibre

import (
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Metadata struct {
	Books   []*Book   `json:"books"`
	Authors []*Author `json:"authors"`
}

func Read(path string) (*Metadata, error) {
	dbpath := filepath.Join(path, "metadata.db")
	db, err := sqlx.Connect("sqlite3", "file:"+dbpath+"?mode=ro")
	if err != nil {
		return nil, err
	}

	m := Metadata{}
	if err := db.Select(&m.Books, `SELECT * FROM books`); err != nil {
		return nil, err
	}
	if err := db.Select(&m.Authors, `SELECT * FROM authors`); err != nil {
		return nil, err
	}

	return &m, nil
}

package calibre

import (
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Library struct {
	Path string
	DB   *sqlx.DB
}

func New(path string, db *sqlx.DB) *Library {
	return &Library{Path: path, DB: db}
}

func Open(path string) (*Library, error) {
	dbpath := filepath.Join(path, "metadata.db")
	db, err := sqlx.Open("sqlite3", "file:"+dbpath+"?mode=ro")
	if err != nil {
		return nil, err
	}
	return New(path, db), nil
}

func (l *Library) Books() (out []Book, err error) {
	err = l.DB.Select(&out, `SELECT * FROM books`)
	return
}

package calibre

import (
	"time"
)

type Book struct {
	ID           int        `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Sort         string     `json:"sort" db:"sort"`
	Timestamp    *time.Time `json:"timestamp" db:"timestamp"`
	PubDate      *time.Time `json:"pubdate" db:"pubdate"`
	SeriesIndex  float64    `json:"series_index" db:"series_index"`
	AuthorSort   string     `json:"author_sort" db:"author_sort"`
	ISBN         string     `json:"isbn" db:"isbn"`
	LCCN         string     `json:"lccn" db:"lccn"`
	Path         string     `json:"path" db:"path"`
	Flags        int        `json:"flags" db:"flags"`
	UUID         string     `json:"uuid" db:"uuid"`
	HasCover     bool       `json:"has_cover" db:"has_cover"`
	LastModified time.Time  `json:"last_modified" db:"last_modified"`

	// Inlined: comments (id, book, text), UNIQUE on book.
	Comment string `json:"_comment" db:"_comment"`
	// Inlined: languages (id, lang_code), UNIQUE lang_code.
	Languages []string `json:"_languages" db:"-"`

	Data    []*Data   `json:"_data" db:"-"`    // many-to-one
	Authors []*Author `json:"_authors" db:"-"` // many-to-many
}

type Data struct {
	ID               int    `json:"id" db:"id"`
	Book             int    `json:"book" db:"book"`
	Format           string `json:"format" db:"format"`
	UncompressedSize int    `json:"uncompressed_size" db:"uncompressed_size"`
	Name             string `json:"name" db:"name"`
}

type Author struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Sort string `json:"sort" db:"sort"`
	Link string `json:"link" db:"link"`

	// Relationships.
	Books []*Book `json:"-" db:"-"`
}

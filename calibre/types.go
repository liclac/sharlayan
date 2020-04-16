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

	// Inlined: comments (id, UNIQUE book, text).
	Comment string `json:"_comment" db:"_comment"`

	// I: ratings (id, UNIQUE rating), _link (id, UNIQUE(book, rating)).
	// Yes, that's a many-to-many link of score (0-10) proxies to books.
	Rating *int `json:"_rating" db:"_rating"`

	// I: languages (id, UNIQUE lang_code), _link (id, UNIQUE(book, lang_code), item_order).
	// Note: _link.lang_code actually references lang.id, not lang.lang_code.
	Languages []string `json:"_languages" db:"-"`

	Authors IDs `json:"_authors" db:"_authors"` // many-to-many
	Series  IDs `json:"_series" db:"_series"`   // many-to-many
	Tags    IDs `json:"_tags" db:"_tags"`       // many-to-many

	Data       []*Data           `json:"_data" db:"-"` // many-to-one
	PluginData []*BookPluginData `json:"_plugin_data" db:"-"`
}

type Data struct {
	ID               int    `json:"id" db:"id"`
	Book             int    `json:"book" db:"book"`
	Format           string `json:"format" db:"format"`
	UncompressedSize int    `json:"uncompressed_size" db:"uncompressed_size"`
	Name             string `json:"name" db:"name"`
}

type BookPluginData struct {
	ID   int    `json:"id" db:"id"`
	Book int    `json:"book" db:"book"`
	Name string `json:"name" db:"name"`
	Val  string `json:"val" db:"val"`
}

type Author struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE
	Sort string `json:"sort" db:"sort"`
	Link string `json:"link" db:"link"`

	Books IDs `json:"_books" db:"_books"` // many-to-many
}

type Series struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE
	Sort string `json:"sort" db:"sort"`

	Books IDs `json:"_books" db:"_books"` // many-to-many
}

type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE

	Books IDs `json:"_books" db:"_books"` // many-to-many
}

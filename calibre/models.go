package calibre

import (
	"html/template"
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
	// The comment is edited with a WYSIWYG editor in the UI, and stored as HTML.
	CommentRaw template.HTML `json:"_comment_raw" db:"_comment"`
	// CommentRaw has HTML formatted for the Calibre UI's stylesheets, with formatting
	// and classes that don't always work for us. So we run it through a HTML-to-Markdown
	// filter, and then render the Markdown back into HTML.
	Comment string `json:"_comment" db:"-"`

	// I: ratings (id, UNIQUE rating), _link (id, UNIQUE(book, rating)).
	// Yes, that's a many-to-many link of score (0-10) proxies to books.
	Rating *int `json:"_rating" db:"_rating"`

	// I: languages (id, UNIQUE lang_code), _link (id, UNIQUE(book, lang_code), item_order).
	// Note: _link.lang_code actually references lang.id, not lang.lang_code.
	Languages []string `json:"_languages" db:"-"`

	// Many-to-Many
	AuthorIDs IDs       `json:"_author_ids" db:"_authors"`
	Authors   []*Author `json:"-" db:"-"`
	SeriesIDs IDs       `json:"_series_ids" db:"_series"`
	Series    []*Series `json:"-" db:"-"`
	TagIDs    IDs       `json:"_tag_ids" db:"_tags"`
	Tags      []*Tag    `json:"-" db:"-"`

	// Many-to-One
	Data       []*Data       `json:"_data" db:"-"`
	PluginData []*PluginData `json:"_plugin_data" db:"-"`
}

type Data struct {
	ID               int    `json:"id" db:"id"`
	BookID           int    `json:"book_id" db:"book"`
	Format           string `json:"format" db:"format"`
	UncompressedSize int    `json:"uncompressed_size" db:"uncompressed_size"`
	Name             string `json:"name" db:"name"`
}

type PluginData struct {
	ID     int    `json:"id" db:"id"`
	BookID int    `json:"book_id" db:"book"`
	Name   string `json:"name" db:"name"`
	Val    string `json:"val" db:"val"`
}

type Author struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE
	Sort string `json:"sort" db:"sort"`
	Link string `json:"link" db:"link"`

	BookIDs IDs     `json:"_book_ids" db:"_books"` // many-to-many
	Books   []*Book `json:"-" db:"-"`
}

type Series struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE
	Sort string `json:"sort" db:"sort"`

	BookIDs IDs     `json:"_book_ids" db:"_books"` // many-to-many
	Books   []*Book `json:"-" db:"-"`
}

type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE

	BookIDs IDs     `json:"_book_ids" db:"_books"` // many-to-many
	Books   []*Book `json:"-" db:"-"`
}

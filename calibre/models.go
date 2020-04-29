package calibre

import (
	"database/sql"
	"html/template"
	"time"
)

// A single book in your Calibre library.
type Book struct {
	ID          int          `json:"id" db:"id"`                    // Autoincrementing ID.
	UUID        string       `json:"uuid" db:"uuid"`                // Random UUID (v4).
	ISBN        string       `json:"isbn" db:"isbn"`                // Deprecated; see Identifiers.
	LCCN        string       `json:"lccn" db:"lccn"`                // Deprecated; see Identifiers.
	Identifiers []Identifier `json:"identifiers" db:"_identifiers"` // ISBNs, ASINs, etc.
	Flags       int          `json:"flags" db:"flags"`              // Not sure what these do.
	Timestamp   *time.Time   `json:"timestamp" db:"timestamp"`      // When it was added (editable).
	Path        string       `json:"path" db:"path"`                // Path to book's data on disk.
	HasCover    bool         `json:"has_cover" db:"has_cover"`      // A file named cover.jpg in Path.
	Data        []*Data      `json:"data" db:"-"`                   // List of files in Path.

	Title      string        `json:"title" db:"title"`             // eg. "The Fifth Elephant"
	Sort       string        `json:"sort" db:"sort"`               // eg. "Fifth Elephant, The"
	PubDate    *time.Time    `json:"pubdate" db:"pubdate"`         // Publication date.
	Rating     sql.NullInt32 `json:"rating" db:"_rating"`          // Rating (0-10).
	Languages  []string      `json:"languages" db:"-"`             // 3-letter ISO codes, eg. "eng".
	AuthorSort string        `json:"author_sort" db:"author_sort"` // eg. "Pratchett, Terry"
	AuthorIDs  IDs           `json:"authors" db:"_authors"`
	Authors    []*Author     `json:"-" db:"-"`

	// The raw comment, Calibre's "Download metadata" function uses this for a synopsis.
	// This field contains the raw HTML (usually produced by a WYSIWYG editor), and is
	// normally formatted for the Calibre UI's styesheets. See also Comment below.
	CommentRaw template.HTML `json:"comment_raw" db:"_comment"`
	// The same text as CommentRaw, ran through a HTML-to-Markdown filter. Rendering this
	// back into HTML produces more consistent markup than using CommentRaw directly.
	Comment string `json:"comment" db:"-"`

	// A book can be in multiple series, but only has a single index shared between them.
	// This is mostly used for sorting, so it works as long as they're same-order subseries.
	SeriesIndex float64   `json:"series_index" db:"series_index"`
	SeriesIDs   IDs       `json:"series" db:"_series"`
	Series      []*Series `json:"-" db:"-"`

	// A book can have zero or more tags, sometimes incorrectly referred to as categories.
	TagIDs IDs    `json:"tags" db:"_tags"`
	Tags   []*Tag `json:"-" db:"-"`

	PluginData   []*PluginData `json:"plugin_data" db:"-"`
	LastModified time.Time     `json:"last_modified" db:"last_modified"`
}

// An ISBN, MOBI-ASIN, Google Books ID, etc.
type Identifier struct {
	ID     int    `json:"id" db:"id"`
	BookID int    `json:"book_id" db:"book"`
	Type   string `json:"type" db:"type"` // eg. "isbn", "mobi-asin", "google".
	Val    string `json:"val" db:"val"`   // eg. "9781407035208".
}

// A data file inside a book's data directory (Book.Path).
type Data struct {
	ID               int    `json:"id" db:"id"`
	BookID           int    `json:"book_id" db:"book"`
	Format           string `json:"format" db:"format"`
	UncompressedSize int    `json:"uncompressed_size" db:"uncompressed_size"`
	Name             string `json:"name" db:"name"`
}

// Usually a blob of JSON data added by a plugin.
type PluginData struct {
	ID     int    `json:"id" db:"id"`
	BookID int    `json:"book_id" db:"book"`
	Name   string `json:"name" db:"name"`
	Val    string `json:"val" db:"val"`
}

// An author is someone who wrote one or more books. Books can have multiple co-authors.
// The order of co-authors is only loosely preserved by book<->author links' ID sequence.
type Author struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE
	Sort string `json:"sort" db:"sort"`
	Link string `json:"link" db:"link"` // Calibre: Edit Metadata > Manage Authors!

	BookIDs IDs     `json:"books" db:"_books"`
	Books   []*Book `json:"-" db:"-"`
}

// A series of books. Can be used to mean anything, but usually implies continuity.
// Note: A book can belong to multiple series, but it has only one Book.SeriesValue.
type Series struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE
	Sort string `json:"sort" db:"sort"`

	BookIDs IDs     `json:"books" db:"_books"` // many-to-many
	Books   []*Book `json:"-" db:"-"`
}

// Tags that can be applied to books, sometimes incorrectly referred to as categories.
type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // UNIQUE

	BookIDs IDs     `json:"books" db:"_books"` // many-to-many
	Books   []*Book `json:"-" db:"-"`
}

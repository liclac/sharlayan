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
}

type Author struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Sort string `json:"sort" db:"sort"`
	Link string `json:"link" db:"link"`
}

package calibre

import (
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Metadata struct {
	Authors []*Author `json:"authors"`
	Books   []*Book   `json:"books"`
}

func (m Metadata) Author(id int) *Author {
	for _, a := range m.Authors {
		if a.ID == id {
			return a
		}
	}
	return nil
}

func (m Metadata) Book(id int) *Book {
	for _, b := range m.Books {
		if b.ID == id {
			return b
		}
	}
	return nil
}

func Read(path string) (*Metadata, error) {
	dbpath := filepath.Join(path, "metadata.db")
	db, err := sqlx.Connect("sqlite3", "file:"+dbpath+"?mode=ro")
	if err != nil {
		return nil, err
	}

	m := Metadata{}
	if err := db.Select(&m.Authors, `SELECT * FROM authors`); err != nil {
		return nil, err
	}

	// Comments are UNIQUE for a book, so we can inline them right here.
	if err := db.Select(&m.Books, `
        SELECT books.*,
            IFNULL(comments.text, '') AS _comment
        FROM books
        LEFT JOIN comments ON books.id = comments.book
        ORDER BY id
    `); err != nil {
		return nil, err
	}
	for _, book := range m.Books {
		if err := db.Select(&book.Data, `SELECT * FROM data WHERE book = ?`, book.ID); err != nil {
			return nil, err
		}
		if err := db.Select(&book.PluginData, `SELECT * FROM books_plugin_data WHERE book = ?`, book.ID); err != nil {
			return nil, err
		}

		// Confusingly, link.lang_code actually refs lang.id, not lang.lang_code.
		if err := db.Select(&book.Languages, `
            SELECT languages.lang_code
            FROM books_languages_link AS link
            LEFT JOIN languages ON link.lang_code = languages.id
            WHERE book = ?
            ORDER BY item_order ASC
        `, book.ID); err != nil {
			return nil, err
		}
	}

	// Hook up many-to-many links without duplicating objects.
	rows, err := db.Query(`SELECT book, author FROM books_authors_link`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var bookID, authorID int
		if err := rows.Scan(&bookID, &authorID); err != nil {
			return nil, err
		}
		book := m.Book(bookID)
		author := m.Author(authorID)
		if book != nil && author != nil {
			book.Authors = append(book.Authors, author)
			author.Books = append(author.Books, book)
		}
	}

	return &m, nil
}

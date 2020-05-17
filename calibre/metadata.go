package calibre

import (
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"jaytaylor.com/html2text"
)

type Metadata struct {
	Path    string    `json:"path"`
	Tags    []*Tag    `json:"tags"`
	Series  []*Series `json:"series"`
	Authors []*Author `json:"authors"`
	Books   []*Book   `json:"books"`
}

func Read(path string) (*Metadata, error) {
	start := time.Now()
	L := zap.L().Named("calibre")

	// Make sure there's a trailing slash on the path; fixes symlinks.
	path = filepath.Clean(path) + "/"
	L.Info("Reading library...", zap.String("path", path))

	dbpath := filepath.Join(path, "metadata.db")
	dburi := "file:" + dbpath + "?mode=ro"
	L.Debug("Opening database...", zap.String("uri", dburi))
	db, err := sqlx.Connect("sqlite3", dburi)
	if err != nil {
		return nil, err
	}

	m := Metadata{Path: path}

	L.Debug("Loading: Authors...")
	startAuthors := time.Now()
	if err := db.Select(&m.Authors, `SELECT * FROM authors INNER JOIN (
		SELECT author AS id, group_concat(book) _books FROM books_authors_link
		GROUP BY author) AS _books USING (id)
	`); err != nil {
		return nil, err
	}
	L.Info("Loaded: Authors", zap.Int("num", len(m.Authors)),
		zap.Duration("t", time.Since(startAuthors)))

	L.Debug("Loading: Series...")
	startSeries := time.Now()
	if err := db.Select(&m.Series, `SELECT * FROM series INNER JOIN (
		SELECT series AS id, group_concat(book) _books FROM books_series_link
		GROUP BY series) AS _books USING (id)
	`); err != nil {
		return nil, err
	}
	L.Info("Loaded: Series", zap.Int("num", len(m.Series)),
		zap.Duration("t", time.Since(startSeries)))

	L.Debug("Loading: Tags...")
	startTags := time.Now()
	if err := db.Select(&m.Tags, `SELECT * FROM tags INNER JOIN (
		SELECT tag AS id, group_concat(book) _books FROM books_tags_link
		GROUP BY tag) AS _books USING (id)
	`); err != nil {
		return nil, err
	}
	L.Info("Loaded: Tags", zap.Int("num", len(m.Tags)),
		zap.Duration("t", time.Since(startTags)))

	L.Debug("Loading: Books...")
	startBooks := time.Now()
	// Comments are UNIQUE for a book, so we can inline them right here.
	//
	// Ratings, meanwhile, use an odd system where each rating (0-10) is an object
	// in a separate table, many-to-many linked to books... which means a book can
	// technically have more than one rating, although the UI doesn't allow this.
	// Should this still happen somehow, a LEFT JOIN would duplicate the book, so
	// we use a subquery to deduplicate the links before joining with it.
	//
	// The way we get IDs for related top-level items (authors, etc) is a kludge
	// to emulate postgres' array_agg(), we really need a better way to do this.
	if err := db.Select(&m.Books, `
        SELECT books.*,
            COALESCE(comments.text, '') AS _comment,
            ratings.rating AS _rating,
            _authors.authors AS _authors,
            _series.series AS _series,
            _tags.tags AS _tags
        FROM books
        LEFT JOIN comments ON comments.book = books.id
        LEFT JOIN (
            SELECT link.book, ratings.rating
            FROM books_ratings_link AS link
            LEFT JOIN ratings ON link.rating
            GROUP BY link.book
        ) AS ratings ON ratings.book = books.id
        LEFT JOIN (SELECT book, group_concat(author) authors FROM books_authors_link
                   GROUP BY book) AS _authors ON _authors.book = books.id
        LEFT JOIN (SELECT book, group_concat(series) series FROM books_series_link
                   GROUP BY book) AS _series ON _series.book = books.id
        LEFT JOIN (SELECT book, group_concat(tag) tags FROM books_tags_link
                   GROUP BY book) AS _tags ON _tags.book = books.id
        ORDER BY id
    `); err != nil {
		return nil, err
	}
	L.Debug("Loaded: Book objects", zap.Int("num", len(m.Books)),
		zap.Duration("t", time.Since(startBooks)))

	L.Debug("Loading: Book associations...")
	startBookAssocs := time.Now()
	var numBookIdent, numBookData, numBookPluginData, numBookLang int
	for _, book := range m.Books {
		if err := db.Select(&book.Identifiers,
			`SELECT * FROM identifiers WHERE book = ?`, book.ID); err != nil {
			return nil, err
		}
		numBookIdent += len(book.Identifiers)
		if err := db.Select(&book.Data,
			`SELECT * FROM data WHERE book = ?`, book.ID); err != nil {
			return nil, err
		}
		numBookData += len(book.Data)
		if err := db.Select(&book.PluginData,
			`SELECT * FROM books_plugin_data WHERE book = ?`, book.ID); err != nil {
			return nil, err
		}
		numBookPluginData += len(book.PluginData)

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
		numBookLang += len(book.Languages)

		// Link up Many-to-Many relationships.
		for _, id := range book.AuthorIDs {
			author := m.GetAuthor(id)
			book.Authors = append(book.Authors, author)
			author.Books = append(author.Books, book)
		}
		for _, id := range book.SeriesIDs {
			series := m.GetSeries(id)
			book.Series = append(book.Series, series)
			series.Books = append(series.Books, book)
		}
		for _, id := range book.TagIDs {
			tag := m.GetTag(id)
			book.Tags = append(book.Tags, tag)
			tag.Books = append(tag.Books, book)
		}
	}
	L.Debug("Loaded: Book associations",
		zap.Int("idents", numBookIdent), zap.Int("files", numBookData),
		zap.Int("plugin_data", numBookPluginData), zap.Int("langs", numBookLang),
		zap.Duration("t", time.Since(startBookAssocs)))

	L.Debug("Sanitising comments...")
	startBookComments := time.Now()
	var numComments int
	for _, book := range m.Books {
		if book.CommentRaw == "" {
			continue
		}
		// Run CommentRaw through a HTML-to-Markdown filter. If this misbehaves, don't hesitate
		// to do this in-tree; we're dealing with the output of a WYSIWYG editor with known
		// behaviour and quirks on one end, and control the rendering on the other end.
		comment, err := html2text.FromString(string(book.CommentRaw))
		if err != nil {
			L.Warn("Couldn't sanitise RawComment, Comment left blank", zap.Error(err))
		} else {
			book.Comment = comment
		}
		numComments++
	}
	L.Debug("Sanitised comments", zap.Int("num", numComments),
		zap.Duration("t", time.Since(startBookComments)))

	L.Info("Loaded: Books", zap.Int("num", len(m.Books)),
		zap.Duration("t", time.Since(startBooks)))

	L.Debug("Done", zap.Duration("t", time.Since(start)))
	return &m, nil
}

func (m Metadata) GetTag(id int) *Tag {
	for _, t := range m.Tags {
		if t.ID == id {
			return t
		}
	}
	return nil
}

func (m Metadata) GetSeries(id int) *Series {
	for _, s := range m.Series {
		if s.ID == id {
			return s
		}
	}
	return nil
}

func (m Metadata) GetAuthor(id int) *Author {
	for _, a := range m.Authors {
		if a.ID == id {
			return a
		}
	}
	return nil
}

func (m Metadata) GetBook(id int) *Book {
	for _, b := range m.Books {
		if b.ID == id {
			return b
		}
	}
	return nil
}

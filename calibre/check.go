package calibre

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileStatus int

const (
	FileStatusMissing FileStatus = iota
	FileStatusOK
	FileStatusOrphan
)

func (fs FileStatus) String() string {
	switch fs {
	case FileStatusMissing:
		return "-"
	case FileStatusOK:
		return " "
	case FileStatusOrphan:
		return "+"
	default:
		return strconv.Itoa(int(fs))
	}
}

func (fs FileStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(fs.String())
}

type CheckReport struct {
	Files map[string]FileStatus
}

func (m *Metadata) Check() (*CheckReport, error) {
	report := &CheckReport{
		Files: make(map[string]FileStatus),
	}

	// List files on disk, start by assuming everything is an orphan.
	if err := filepath.Walk(m.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() ||
			strings.HasSuffix(path, ".db") || strings.HasSuffix(path, ".json") ||
			strings.HasSuffix(path, ".opf") || strings.HasSuffix(path, ".jpg") {
			return err
		}
		report.Files[strings.TrimPrefix(path, m.Path)] = FileStatusOrphan
		return nil
	}); err != nil {
		return nil, err
	}

	// List files in your library; a match with disk is OK, else it's missing.
	for _, book := range m.Books {
		for _, d := range book.Data {
			filename := d.Name + "." + strings.ToLower(d.Format)
			path := filepath.Join(book.Path, filename)
			if _, ok := report.Files[path]; ok {
				report.Files[path] = FileStatusOK
			} else {
				report.Files[path] = FileStatusMissing
			}
		}
	}

	return report, nil
}

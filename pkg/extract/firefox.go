package extract

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"uuid"
)

type FirefoxBookmark struct {
	Guid         string            `json:"guid"`
	Title        string            `json:"title"`
	Index        int64             `json:"index"`
	DateAdded    int64             `json:"dateAdded"`
	LastModified int64             `json:"lastModified"`
	Id           int64             `json:"id"`
	TypeCode     int               `json:"typeCode"`
	Type         string            `json:"type"`
	Root         string            `json:"root"`
	Children     []FirefoxBookmark `json:"children"`
	Tags         string            `json:"tags"`
	Uri          string            `json:"uri"`
	IconUri      string            `json:"iconUri"`
}

// decompressFirefoxBookmarks executes the `dejsonlz4` programme to decompress a
// `.jsonlz4` file (Firefox's slightly-unique compression file-type) into `.json`.
func decompressFirefoxBookmarks(inputJSONlz4 string, outputJSON string) error {
	return exec.Command("./dejsonlz4/dejsonlz4", inputJSONlz4, outputJSON).Run()

}

// findFirefoxBookmarkFile searches the expected configuration directory for the
// latest backup bookmark file.
func findFirefoxBookmarkFile() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	} else {
		configDir = filepath.Join(configDir, "Mozilla", "Firefox", "Profiles")
		profiles, err := os.ReadDir(configDir)
		if err != nil {
			log.Fatal(err)
		}
		for _, profile := range profiles {
			if profile.IsDir() {
				profileDir := filepath.Join(configDir, profile.Name())
				entries, err := os.ReadDir(profileDir)
				if err != nil {
					log.Fatal(err)
				}
				for _, entry := range entries {
					info, err := entry.Info()
					if err != nil {
						log.Fatal(err)
					}
					if info.IsDir() && info.Name() == "bookmarkbackups" {
						backupDir := filepath.Join(profileDir, info.Name())
						bookmarkFiles, err := os.ReadDir(backupDir)
						if err != nil {
							log.Fatal(err)
						}
						if len(bookmarkFiles) == 0 {
							log.Fatalf("No bookmark files found in %s", profileDir)
						}
						return filepath.Join(backupDir, bookmarkFiles[len(bookmarkFiles)-1].Name())
					}
				}

			}

		}
	}
	log.Fatal("No backup bookmark file found.")
	return ""
}

// ExtractFirefoxBookmarksFrom extracts a (collection of) FirefoxBookmark from
// a particular Firefox backup bookmark file.
// If it is compressed (i.e. .jsonlz4) it will automatically be decompressed.
func ExtractFirefoxBookmarksFrom(bookmarkFile string) (bookmarks *FirefoxBookmark, err error) {
	if strings.HasSuffix(bookmarkFile, ".jsonlz4") {
		compressed := bookmarkFile
		bookmarkFile = filepath.Join(os.TempDir(), uuid.New().String())
		err = decompressFirefoxBookmarks(compressed, bookmarkFile)
	}
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(bookmarkFile, os.O_RDONLY, 0666)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error opening the extracted bookmark file:\n%s", err))
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	bookmarks = new(FirefoxBookmark)
	err = json.UnmarshalRead(f, bookmarks)
	return bookmarks, err
}

// ExtractFirefoxBookmarks extracts all the bookmarks from Firefox, automatically
// locating the right backup bookmark file.
// See also ExtractFirefoxBookmarksFrom
func ExtractFirefoxBookmarks() (*FirefoxBookmark, error) {
	bookmarkFile := findFirefoxBookmarkFile()
	return ExtractFirefoxBookmarksFrom(bookmarkFile)
}

package main

import (
	"encoding/json/v2"
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
	Index        uint64            `json:"index"`
	DateAdded    uint64            `json:"dateAdded"`
	LastModified uint64            `json:"lastModified"`
	Id           uint64            `json:"id"`
	TypeCode     uint8             `json:"typeCode"`
	Type         string            `json:"type"`
	Root         string            `json:"root"`
	Children     []FirefoxBookmark `json:"children"`
	Tags         string            `json:"tags"`
	Uri          string            `json:"uri"`
	IconUri      string            `json:"iconUri"`
}

// decompressFirefoxBookmarks executes the `dejsonlz4` programme to decompress a
// `.jsonlz4` file (Firefox's slightly-unique compression file-type) into `.json`.
func decompressFirefoxBookmarks(inputJSONlz4 string, outputJSON string) {
	err := exec.Command("./dejsonlz4/dejsonlz4", inputJSONlz4, outputJSON).Run()
	if err != nil {
		log.Fatal(err)
	}
}

// findFirefoxBookmarkFile searches the expected configuration directory for the
// latest backup bookmark file.
func findFirefoxBookmarkFile() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	} else {
		configDir = filepath.Join(configDir, "Roaming", "Mozilla", "Firefox", "Profiles")
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

// extractFirefoxBookmarksFrom extracts a (collection of) FirefoxBookmark from
// a particular Firefox backup bookmark file.
// If it is compressed (i.e. .jsonlz4) it will automatically be decompressed.
func extractFirefoxBookmarksFrom(bookmarkFile string) *FirefoxBookmark {
	if strings.HasSuffix(bookmarkFile, ".jsonlz4") {
		compressed := bookmarkFile
		bookmarkFile = filepath.Join(os.TempDir(), uuid.New().String())
		decompressFirefoxBookmarks(compressed, bookmarkFile)
	}

	f, err := os.OpenFile(bookmarkFile, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal("Error opening the extracted bookmark file:", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	bm := new(FirefoxBookmark)
	err = json.UnmarshalRead(f, bm)
	if err != nil {
		log.Fatal("Error unmarshalling the bookmark file into a Go struct:", err)
	}

	// Write to a new JSON file.
	// TODO: Do something useful here instead of proving extraction works.
	err = os.MkdirAll("data/firefox", os.ModePerm)
	if err != nil {
		log.Fatal("Error creating directories for the bookmark file:", err)
	}

	out, err := os.OpenFile("./data/firefox/bookmarks.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Error while creating bookmark file:", err)
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	err = json.MarshalWrite(out, bm)
	if err != nil {
		log.Fatal("Error marshalling new bookmark file:", err)
	}
	return bm
}

// extractFirefoxBookmarks extracts all the bookmarks from Firefox, automatically
// locating the right backup bookmark file.
// See also extractFirefoxBookmarksFrom
func extractFirefoxBookmarks() *FirefoxBookmark {
	bookmarkFile := findFirefoxBookmarkFile()
	return extractFirefoxBookmarksFrom(bookmarkFile)
}

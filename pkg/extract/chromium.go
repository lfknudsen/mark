package extract

import (
	"encoding/json/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Browsers = map[string]string{
	"Brave":  filepath.Join("BraveSoftware", "Brave-Browser"),
	"Edge":   filepath.Join("Microsoft", "Edge"),
	"Chrome": filepath.Join("Google", "Chrome"),
}

type ChromiumRoots struct {
	BookmarkBar ChromiumBookmark `json:"bookmark_bar"`
	Other       ChromiumBookmark `json:"other"`
	Synced      ChromiumBookmark `json:"synced"`
}

type ChromiumBookmarkCollection struct {
	Checksum string        `json:"checksum"`
	Roots    ChromiumRoots `json:"roots"`
	Version  float32       `json:"version"`
}

type ChromiumMetaInfo struct {
	PowerBookmarkMeta string `json:"power_bookmark_meta"`
}

type ChromiumBookmark struct {
	Checksum     string             `json:"checksum"`
	Children     []ChromiumBookmark `json:"children"`
	Guid         string             `json:"guid"`
	DateAdded    int64              `json:"date_added"`
	DateLastUsed int64              `json:"date_last_used"`
	DateModified int64              `json:"date_modified"`
	Id           int64              `json:"id"`
	Type         string             `json:"type"`
	Url          string             `json:"url"`
	Name         string             `json:"name"`
	ShowIcon     bool               `json:"show_icon"`
	Source       string             `json:"source"`
	MetaInfo     ChromiumMetaInfo   `json:"meta_info"`
	VisitCount   int64              `json:"visit_count"`
}

func (cb *ChromiumBookmark) SplitChildren() (urls []ChromiumBookmark, folders []ChromiumBookmark, other []ChromiumBookmark) {
	for _, child := range cb.Children {
		if child.Type == "folder" {
			folders = append(folders, child)
		} else if child.Type == "url" {
			urls = append(urls, child)
		} else {
			other = append(other, child)
		}
	}
	return urls, folders, other
}

func (c *ChromiumBookmarkCollection) JoinRoots() []ChromiumBookmark {
	root1Length := len(c.Roots.BookmarkBar.Children)
	root2Length := len(c.Roots.Other.Children)
	root3Length := len(c.Roots.Synced.Children)
	output := make([]ChromiumBookmark, 0, root1Length+root2Length+root3Length)
	output = append(output, c.Roots.BookmarkBar.Children...)
	output = append(output, c.Roots.Other.Children...)
	output = append(output, c.Roots.Synced.Children...)
	return output
}

// findChromiumBookmarkFile searches the expected configuration directory for the
// latest backup bookmark file.
func findChromiumBookmarkFile(browser string, profile string) string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("No backup bookmark file found.")
	}
	browser = strings.ToTitle(browser)
	softwareDir := Browsers[browser]
	if softwareDir == "" {
		log.Fatalf("Browser %s not found. Acceptable choices are:\n%v\n", browser, Browsers)
	}
	configDirRoot := filepath.Join(configDir, "Local")
	if _, err := os.Stat(configDirRoot); err != nil {
		configDirRoot = filepath.Join(configDir, "Roaming")
		if _, err := os.Stat(configDirRoot); err != nil {
			configDirRoot = configDir
		}
	}
	return filepath.Join(configDirRoot, softwareDir, "User Data", profile, "Bookmarks")
}

// ExtractChromiumBookmarksFrom extracts an ChromiumBookmarkCollection from
// a particular Chromium backup bookmark file (these are always named simply "Bookmarks",
// but are actually JSON files).
func ExtractChromiumBookmarksFrom(bookmarkFile string) (bookmarks *ChromiumBookmarkCollection, err error) {
	f, err := os.OpenFile(bookmarkFile, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	bm := new(ChromiumBookmarkCollection)
	err = json.UnmarshalRead(f, bm)
	return bm, err
}

// ExtractChromiumBookmarks extracts all the bookmarks from a Chromium browser,
// automatically locating the right backup bookmark file.
//
// `browser` must be one of Browsers.
//
// `profile` is the name of the profile whose bookmarks you wish to extract.
// If you only have one profile, try "Default".
//
// See also ExtractChromiumBookmarksFrom
func ExtractChromiumBookmarks(browser string, profile string) (*ChromiumBookmarkCollection, error) {
	if profile == "" {
		profile = "Default"
	}
	bookmarkFile := findChromiumBookmarkFile(browser, profile)
	return ExtractChromiumBookmarksFrom(bookmarkFile)
}

func extractChromiumBookmarksDefault(browser string) (*ChromiumBookmarkCollection, error) {
	bookmarkFile := findChromiumBookmarkFile(browser, "Default")
	return ExtractChromiumBookmarksFrom(bookmarkFile)
}

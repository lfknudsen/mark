package main

import (
	"path"
	"time"
	"uuid"
)

var CurrentVersionBookmark = "1.0.0"
var CurrentVersionBookmarkCollection = "1.0.0"

type Bookmark struct {
	Version      string
	Url          string
	Name         string
	Desc         string
	DateAdded    int64
	DateModified int64
	DateLastUsed int64
	Guid         string
	Icon         string
	Tags         []string
}

type BookmarkCollection struct {
	Version             string
	Name                string
	Guid                string
	DateAdded           int64
	DateModified        int64
	DateAddedContent    int64
	DateModifiedContent int64
	DateLastUsedContent int64
	Icon                string
	Desc                string
	Tags                []string
	Children            []Bookmark
	Subfolders          []BookmarkCollection
}

// NewBookmark creates a new Bookmark struct from a URL, extracting a suitable name and
// leaving the description empty.
func NewBookmark(url string) Bookmark {
	now := time.Now().Unix()
	return Bookmark{
		Url:          url,
		Name:         path.Base(url),
		Desc:         "",
		DateAdded:    now,
		DateModified: now,
		DateLastUsed: now,
		Guid:         uuid.New().String(),
		Icon:         "",
		Tags:         []string{},
	}
}

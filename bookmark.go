package main

import (
	"log"
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
		Version:      CurrentVersionBookmark,
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

func NewBookmarkNamed(url string, name string) Bookmark {
	now := time.Now().Unix()
	return Bookmark{
		Version:      CurrentVersionBookmark,
		Url:          url,
		Name:         name,
		Desc:         "",
		DateAdded:    now,
		DateModified: now,
		DateLastUsed: now,
		Guid:         uuid.New().String(),
		Icon:         "",
		Tags:         []string{},
	}
}

func BookmarkFromChromium(cb *ChromiumBookmark) Bookmark {
	return Bookmark{
		Version:      CurrentVersionBookmark,
		Url:          cb.Url,
		Name:         cb.Name,
		Desc:         "",
		DateAdded:    cb.DateAdded,
		DateModified: cb.DateModified,
		DateLastUsed: cb.DateLastUsed,
		Guid:         cb.Guid,
		Icon:         "",
		Tags:         []string{},
	}
}

func BookmarkCollectionFromChromium(cb *ChromiumBookmark) *BookmarkCollection {
	out := BookmarkCollection{
		Version:      CurrentVersionBookmarkCollection,
		Name:         cb.Name,
		Guid:         cb.Guid,
		DateAdded:    cb.DateAdded,
		DateModified: cb.DateModified,
		Icon:         "",
		Desc:         "",
		Tags:         []string{},
		Children:     nil,
		Subfolders:   nil,
	}
	urls, folders, other := cb.SplitChildren()
	if len(other) > 0 {
		log.Fatalf("Other wasn't empty: %v", other)
	}

	out.Children = make([]Bookmark, 0, len(urls))
	for _, url := range urls {
		out.Children = append(out.Children, BookmarkFromChromium(&url))
	}
	out.Subfolders = make([]BookmarkCollection, 0, len(folders))
	for _, folder := range folders {
		out.Subfolders = append(out.Subfolders, *BookmarkCollectionFromChromium(&folder))
	}
	return &out
}

func BookmarkCollectionFromChromiumCollection(cb *ChromiumBookmarkCollection) *BookmarkCollection {
	now := time.Now().Unix()
	out := BookmarkCollection{
		Version:      CurrentVersionBookmarkCollection,
		Name:         "",
		Guid:         uuid.New().String(),
		DateAdded:    now,
		DateModified: now,
		Icon:         "",
		Desc:         "",
		Tags:         []string{},
		Children:     []Bookmark{},
		Subfolders:   make([]BookmarkCollection, 0, 3),
	}
	out.Subfolders = append(out.Subfolders, *BookmarkCollectionFromChromium(&cb.Roots.BookmarkBar))
	out.Subfolders = append(out.Subfolders, *BookmarkCollectionFromChromium(&cb.Roots.Other))
	out.Subfolders = append(out.Subfolders, *BookmarkCollectionFromChromium(&cb.Roots.Synced))
	return &out
}

func BookmarkCollectionFromFirefox(ff *FirefoxBookmark) *BookmarkCollection {
	now := time.Now().Unix()
	return &BookmarkCollection{
		Version:             CurrentVersionBookmarkCollection,
		Name:                ff.Title,
		Guid:                ff.Guid,
		DateAdded:           ff.DateAdded,
		DateModified:        ff.LastModified,
		DateAddedContent:    now,
		DateModifiedContent: now,
		DateLastUsedContent: ff.LastModified,
		Icon:                "",
		Desc:                "",
		Tags:                []string{},
		Children:            []Bookmark{},
		Subfolders:          []BookmarkCollection{},
	}
}

package main

import (
	"fmt"
	"log"
	"path"
	"reflect"
	"time"
	"uuid"
)

var CurrentVersionBookmark = "1.0.0"
var CurrentVersionBookmarkCollection = "1.0.0"

type Bookmark struct {
	Version      string   `sql:"version" json:"version,omitempty"`
	Url          string   `sql:"url" json:"url,omitempty"`
	Name         string   `sql:"name" json:"name,omitempty"`
	Desc         string   `sql:"desc" json:"desc,omitempty"`
	DateAdded    int64    `sql:"date_added" json:"date_added,omitempty"`
	DateModified int64    `sql:"date_modified" json:"date_modified,omitempty"`
	DateLastUsed int64    `sql:"date_last_used" json:"date_last_used,omitempty"`
	Guid         string   `sql:"guid" json:"guid,omitempty"`
	Icon         string   `sql:"icon" json:"icon,omitempty"`
	Tags         []string `sql:"-" json:"tags,omitempty"`
}

func getBookmarkFieldNames() []string {
	fields := reflect.VisibleFields(reflect.TypeFor[Bookmark]())
	exportedFields := make([]string, 0, len(fields))
	for _, field := range fields {
		tag := field.Tag.Get("sql")
		if tag != "-" {
			exportedFields = append(exportedFields, tag)
		}
	}
	return exportedFields
}

func (b Bookmark) getValues() []string {
	v := reflect.ValueOf(b)
	fields := reflect.VisibleFields(reflect.TypeFor[Bookmark]())
	output := make([]string, 0, len(fields))
	for _, field := range fields {
		tag := field.Tag.Get("sql")
		if tag != "-" {
			output = append(output, fmt.Sprint(v.FieldByName(field.Name)))
		}
	}
	return output
}

func (b Bookmark) String() string {
	return fmt.Sprintf("(%s, %s, %s, %s, %d, %d, %d, %s, %s)",
		b.Version, b.Url, b.Name, b.Desc, b.DateAdded, b.DateModified, b.DateLastUsed, b.Guid, b.Icon)
}

type BookmarkCollection struct {
	Version             string               `sql:"version" json:"version,omitempty"`
	Name                string               `sql:"name" json:"name,omitempty"`
	Guid                string               `sql:"guid" json:"guid,omitempty"`
	DateAdded           int64                `sql:"date_added" json:"date_added,omitempty"`
	DateModified        int64                `sql:"date_modified" json:"date_modified,omitempty"`
	DateAddedContent    int64                `sql:"date_added_content" json:"date_added_content,omitempty"`
	DateModifiedContent int64                `sql:"date_modified_content" json:"date_modified_content,omitempty"`
	DateLastUsedContent int64                `sql:"date_last_used_content" json:"date_last_used_content,omitempty"`
	Icon                string               `sql:"icon" json:"icon,omitempty"`
	Desc                string               `sql:"desc" json:"desc,omitempty"`
	Tags                []string             `sql:"-" json:"tags,omitempty"`
	Children            []Bookmark           `sql:"-" json:"children,omitempty"`
	Subfolders          []BookmarkCollection `sql:"-" json:"subfolders,omitempty"`
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

func BookmarkFromFirefox(ff *FirefoxBookmark) Bookmark {
	return Bookmark{
		Version:      CurrentVersionBookmark,
		Url:          ff.Uri,
		Name:         ff.Title,
		Desc:         "",
		DateAdded:    ff.DateAdded,
		DateModified: ff.LastModified,
		DateLastUsed: 0,
		Guid:         ff.Guid,
		Icon:         "",
		Tags:         []string{},
	}
}

func BookmarkCollectionFromFirefox(ff *FirefoxBookmark) *BookmarkCollection {
	now := time.Now().Unix()
	out := BookmarkCollection{
		Version:             CurrentVersionBookmarkCollection,
		Name:                ff.Title,
		Guid:                ff.Guid,
		DateAdded:           ff.DateAdded,
		DateModified:        ff.LastModified,
		DateAddedContent:    now,
		DateModifiedContent: now,
		DateLastUsedContent: 0,
		Icon:                "",
		Desc:                "",
		Tags:                []string{},
		Children:            []Bookmark{},
		Subfolders:          []BookmarkCollection{},
	}

	if len(ff.Children) > 0 {
		for _, child := range ff.Children {
			if len(child.Children) == 0 {
				out.Children = append(out.Children, BookmarkFromFirefox(&child))
			} else {
				out.Subfolders = append(out.Subfolders, *BookmarkCollectionFromFirefox(&child))
			}
		}
	}
	return &out
}

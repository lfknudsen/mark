package pkg

import (
	"log"
	"time"
	"uuid"

	. "mark/pkg/extract"
)

const CurrentVersionBookmarkCollection = "1.0.0"

type BookmarkCollection struct {
	Version             string                `sql:"version" json:"version,omitempty"`
	Name                string                `sql:"name" json:"name,omitempty"`
	Guid                string                `sql:"guid" json:"guid,omitempty"`
	DateAdded           int64                 `sql:"date_added" json:"date_added,omitempty"`
	DateModified        int64                 `sql:"date_modified" json:"date_modified,omitempty"`
	DateAddedContent    int64                 `sql:"date_added_content" json:"date_added_content,omitempty"`
	DateModifiedContent int64                 `sql:"date_modified_content" json:"date_modified_content,omitempty"`
	DateLastUsedContent int64                 `sql:"date_last_used_content" json:"date_last_used_content,omitempty"`
	Icon                string                `sql:"icon" json:"icon,omitempty"`
	Desc                string                `sql:"desc" json:"desc,omitempty"`
	Tags                []string              `sql:"-" json:"tags,omitempty"`
	Children            []Bookmark            `sql:"-" json:"children,omitempty"`
	Subfolders          []*BookmarkCollection `sql:"-" json:"subfolders,omitempty"`
}

func NewBookmarkCollection(name string) *BookmarkCollection {
	now := time.Now().Unix()
	out := &BookmarkCollection{
		Version:             CurrentVersionBookmarkCollection,
		Name:                name,
		Guid:                uuid.New().String(),
		DateAdded:           now,
		DateModified:        now,
		DateAddedContent:    now,
		DateModifiedContent: now,
		DateLastUsedContent: now,
		Icon:                "",
		Desc:                "",
		Tags:                []string{},
		Children:            make([]Bookmark, 0),
		Subfolders:          []*BookmarkCollection{},
	}
	return out
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
	out.Subfolders = make([]*BookmarkCollection, 0, len(folders))
	for _, folder := range folders {
		out.Subfolders = append(out.Subfolders, BookmarkCollectionFromChromium(&folder))
	}
	return &out
}

func BookmarkCollectionFromChromiumCollection(cb *ChromiumBookmarkCollection) *BookmarkCollection {
	if cb == nil {
		return nil
	}
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
		Subfolders:   make([]*BookmarkCollection, 0, 3),
	}
	out.Subfolders = append(out.Subfolders, BookmarkCollectionFromChromium(&cb.Roots.BookmarkBar))
	out.Subfolders = append(out.Subfolders, BookmarkCollectionFromChromium(&cb.Roots.Other))
	out.Subfolders = append(out.Subfolders, BookmarkCollectionFromChromium(&cb.Roots.Synced))
	return &out
}

func BookmarkCollectionFromFirefox(ff *FirefoxBookmark) *BookmarkCollection {
	if ff == nil {
		return nil
	}
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
		Subfolders:          []*BookmarkCollection{},
	}

	if len(ff.Children) > 0 {
		for _, child := range ff.Children {
			if len(child.Children) == 0 {
				out.Children = append(out.Children, BookmarkFromFirefox(&child))
			} else {
				out.Subfolders = append(out.Subfolders, BookmarkCollectionFromFirefox(&child))
			}
		}
	}
	return &out
}

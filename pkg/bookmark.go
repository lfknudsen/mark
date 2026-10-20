package pkg

import (
	"fmt"
	. "mark/pkg/extract"
	"path"
	"reflect"
	"time"
	"uuid"
)

const CurrentVersionBookmark = "1.0.0"

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

func GetBookmarkFieldNames() []string {
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

func (b Bookmark) GetValues() []string {
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

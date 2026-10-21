package model

import (
	"mark/pkg"
	"mark/pkg/db"
)

type Bookmark pkg.Bookmark

type Folder pkg.BookmarkCollection

type Directory db.SimpleFolder

type FlatFolder db.FlatFolder

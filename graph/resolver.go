package graph

//go:generate go tool gqlgen generate

import "mark/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	bookmarks []*model.Bookmark
}

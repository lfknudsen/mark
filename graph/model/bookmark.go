package model

type Bookmark struct {
	ID   string `json:"id"`
	Url  string `json:"url"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

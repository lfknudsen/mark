package main

import (
	"encoding/json/v2"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var dir = "data"

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "extract" {
		if len(args) > 1 {
			extractChromiumBookmarksFrom(args[1])
		} else {
			extractChromiumBookmarks("Edge", "Default")
		}
		return
	}

	bookmarks := make([]Bookmark, 0, len(args))
	for _, arg := range args {
		bookmarks = append(bookmarks, NewBookmark(arg))
	}
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	for _, mark := range bookmarks {
		out, err := json.Marshal(mark)
		if err != nil {
			log.Println(err)
		}
		filename := filepath.Join(dir, fmt.Sprintf("%s.json", mark.Name))
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.ModePerm)
		if err != nil {
			i := 1
			for err != nil {
				filename = filepath.Join(dir, fmt.Sprintf("%s (%d).json", mark.Name, i))
				f, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.ModePerm)
			}
		}
		_, err = f.Write(out)
		if err != nil {
			log.Println(err)
		}
		_ = f.Close()
	}
}

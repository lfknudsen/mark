package pkg

import (
	"errors"
	. "mark/pkg/extract"
	"strings"
)

func extract(file string, browser string, profile string) (bookmarks *BookmarkCollection, err error) {
	if file != "" {
		if browser == "" {
			err = errors.New("no browser name specified")
		}
		switch strings.ToLower(browser) {
		case "edge", "brave", "chrome", "chromium":
			var imported *ChromiumBookmarkCollection
			imported, err = ExtractChromiumBookmarksFrom(file)
			if err == nil {
				bookmarks = BookmarkCollectionFromChromiumCollection(imported)
			}
		case "firefox":
			var imported *FirefoxBookmark
			imported, err = ExtractFirefoxBookmarksFrom(file)
			if err == nil {
				bookmarks = BookmarkCollectionFromFirefox(imported)
			}
		default:
			err = errors.New("did not recognise the browser name " + browser)
		}
	} else if browser != "" {
		switch strings.ToLower(browser) {
		case "edge", "brave", "chrome", "chromium":
			var imported *ChromiumBookmarkCollection
			imported, err = ExtractChromiumBookmarks(browser, profile)
			if err == nil {
				bookmarks = BookmarkCollectionFromChromiumCollection(imported)
			}
		case "firefox":
			var imported *FirefoxBookmark
			imported, err = ExtractFirefoxBookmarks()
			if err == nil {
				bookmarks = BookmarkCollectionFromFirefox(imported)
			}
		default:
			err = errors.New("did not recognise the browser name " + browser)
		}
	}
	return bookmarks, nil
}

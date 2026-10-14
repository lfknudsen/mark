package main

import (
	"encoding/json/v2"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var dir = "data"

var programmeVersion = "0.0.1"

var importCommand = "extract"
var fileShorthand = "-f"
var fileLongform = "--file"
var browserShorthand = "-b"
var browserLongform = "--browser"
var profileShorthand = "-p"
var profileLongform = "--profile"
var defaultProfileName = "Default"

var HELP = fmt.Sprintf("Mark v%s\n", programmeVersion) +
	"Usage:\n" +
	"  $ mark <URLs...>\n" +
	fmt.Sprintf("  $ mark %s (%s|%s) <browser name> [Options]\n", importCommand, browserShorthand, browserLongform) +
	"Extract options:\n" +
	fmt.Sprintf("%s <file>\n%s <file>\t\t\tImport from a specific bookmark file\n", fileLongform, fileShorthand) +
	"\n" +
	fmt.Sprintf("%s <profile name>\n%s <profile name>\t\tImport from this profile. If unsure, or on Edge with a Microsoft account, leave blank.", profileLongform, profileShorthand)

func parseExtractCommand(args []string) {
	file := ""
	browser := ""
	profile := defaultProfileName

	if len(args) == 0 || len(args)%2 != 0 {
		fmt.Print(HELP)
		os.Exit(1)
	}
	if len(args) == 1 {
		browser = args[0]
	}
	for len(args) >= 2 {
		switch args[0] {
		case fileShorthand, fileLongform:
			file = args[1]
		case browserShorthand, browserLongform:
			browser = args[1]
		case profileShorthand, profileLongform:
			profile = args[1]
		}
		args = args[2:]
	}
	if file != "" {
		if browser == "" {
			log.Fatal("Please specify a browser name.")
		}
		switch strings.ToLower(browser) {
		case "edge", "brave", "chrome", "chromium":
			extractChromiumBookmarksFrom(file)
		case "firefox":
			extractFirefoxBookmarksFrom(file)
		}
		return
	}
	if browser != "" {
		switch strings.ToLower(browser) {
		case "edge", "brave", "chrome", "chromium":
			extractChromiumBookmarks(browser, profile)
		case "firefox":
			extractFirefoxBookmarks()
		}
		return
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || (len(args) >= 1 && (args[0] == "-h" || args[0] == "--help")) {
		fmt.Print(HELP)
		os.Exit(0)
	}
	if len(args) > 0 && args[0] == importCommand {
		parseExtractCommand(args[1:])
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

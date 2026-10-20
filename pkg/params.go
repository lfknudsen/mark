package pkg

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const ImportCommand = "extract"
const fileShorthand = "-f"
const fileLongform = "--file"
const browserShorthand = "-b"
const browserLongform = "--browser"
const profileShorthand = "-p"
const profileLongform = "--profile"
const defaultProfileName = "Default"

func getExecPath() string {
	p, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(p)
}

const VERSION = "0.8.0"

var execPath = getExecPath()
var markDir = filepath.Join(execPath, ".mark")
var DataDir = filepath.Join(markDir, ".data")
var Port = flag.Uint("port", 8081, "Port to listen on")

var HELP = "Mark v" + VERSION +
	"\nUsage:\n" +
	"  $ mark <URLs...>\n" +
	"                 Add <URL(s)> to your bookmarks.\n" +
	"  $ mark list\n" +
	"                 Print all database entries to standard output.\n" +
	"  $ mark [-port <number>] server\n" +
	"                 Start the GraphQL server." +
	fmt.Sprintf("  $ mark %s (%s|%s) <browser name> [Options]\n", ImportCommand, browserShorthand, browserLongform) +
	"Extract options:\n" +
	fmt.Sprintf("%s <file>\n%s <file>\t\t\tImport from a specific bookmark file\n", fileLongform, fileShorthand) +
	"\n" +
	fmt.Sprintf("%s <profile name>\n%s <profile name>\t\tImport from this profile. If unsure, or on Edge with a Microsoft account, leave blank.", profileLongform, profileShorthand)

func ParseExtractCommand(args []string) (*BookmarkCollection, error) {
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
	return extract(file, browser, profile)
}

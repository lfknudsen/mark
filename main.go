package main

import (
	"encoding/json/v2"
	"fmt"
	"log"
	"mark/graph"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

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
var dataDir = filepath.Join(markDir, ".data")

func main() {
	args := os.Args[1:]
	if len(args) == 0 || (len(args) >= 1 && (args[0] == "-h" || args[0] == "--help")) {
		fmt.Print(HELP)
		os.Exit(0)
	}

	initDirectories()
	err := initDatabase()
	if err != nil {
		log.Fatal(err)
	}

	if args[0] == importCommand {
		imported := parseExtractCommand(args[1:])
		fmt.Print("Inserting")
		wg := new(sync.WaitGroup)
		var rows int64
		wg.Go(
			func() { rows, err = insertCollection(imported) },
		)
		wg.Wait()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(rows, "rows affected in total.")
		return
	}

	if args[0] == "list" {
		bookmarks, err := listBookmarks()
		if err != nil {
			log.Fatal(err)
		}
		fields := getBookmarkFieldNames()
		lines := make([]string, 0, len(bookmarks)+1)
		lines = append(lines, strings.Join(fields, "\t"))
		for _, bookmark := range bookmarks {
			values := bookmark.getValues()
			values[1] = capStringLength(values[1])
			values[2] = capStringLength(values[2])
			values[3] = capStringLength(values[3])
			lines = append(lines, strings.Join(values, "\t"))
		}
		fmt.Println("Found", len(bookmarks), "bookmarks:")
		text := strings.Join(lines, "\n")
		tabWriter := tabwriter.NewWriter(os.Stdout, 3, 1, 3, ' ', 0)
		_, err = fmt.Fprintln(tabWriter, text)
		if err != nil {
			log.Fatal(err)
		}
		_ = tabWriter.Flush()
		os.Exit(0)
	}

	if args[0] == "server" {
		startServer()
	}

	bookmarks := make([]Bookmark, 0, len(args))
	for _, arg := range args {
		bookmarks = append(bookmarks, NewBookmark(arg))
	}
	rows, err := insert(bookmarks...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows, "rows affected in total.")
}

func createJsonBookmarks(bookmarks ...Bookmark) {
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	for _, mark := range bookmarks {
		out, err := json.Marshal(mark)
		if err != nil {
			log.Println(err)
		}
		filename := filepath.Join(dataDir, fmt.Sprintf("%s.json", mark.Name))
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.ModePerm)
		if err != nil {
			i := 1
			for err != nil {
				filename = filepath.Join(dataDir, fmt.Sprintf("%s (%d).json", mark.Name, i))
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

func capStringLength(str string) string {
	if len(str) > 33 {
		str = str[:30] + "..."
	}
	return str
}

func initDirectories() {
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

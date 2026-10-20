package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"modernc.org/sqlite"
)

const connectCommand = "connect"
const connectionString = "file:data/mark.db"

const insertStatement = `INSERT INTO bookmarks
    (version, url, name, desc, date_added, date_modified, date_last_used, guid, icon) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? )`

func closeDB(db *sql.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func abort(transaction *sql.Tx) {
	if transaction != nil {
		_ = transaction.Rollback()
	}
}

// connect connects to the bookmark database.
func connect() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer closeDB(db)
	tran, err := db.Begin()
	if err != nil {
		return err
	}
	data, err := os.ReadFile("model.sql")
	if err != nil {
		return err
	}
	res, err := tran.Exec(string(data))
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("Rows affected: %d\n", rowsAffected)
	err = tran.Commit()
	if err != nil {
		return err
	}
	return nil
}

func openDB() (*sql.DB, error) {
	conn, err := sqlite.NewConnector(connectionString)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(conn)
	return db, db.Ping()
}

func insert(bookmarks ...Bookmark) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer closeDB(db)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmnt, err := tx.Prepare(insertStatement)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer closeStmt(stmnt)
	for _, b := range bookmarks {
		res, err := stmnt.Exec(b.Version, b.Url, b.Name, b.Desc, b.DateAdded,
			b.DateModified, b.DateLastUsed, b.Guid, b.Icon)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		fmt.Printf("%d rows were affected by inserting %s.\n", rowsAffected, b.Url)
	}
	return tx.Commit()
}

func closeStmt(stmnt *sql.Stmt) {
	if stmnt != nil {
		_ = stmnt.Close()
	}
}

func rollback(transaction *sql.Tx) {
	if transaction != nil {
		_ = transaction.Rollback()
	}
}

func listBookmarks() ([]Bookmark, error) {
	output := make([]Bookmark, 0, 10)
	db, err := openDB()
	if err != nil {
		return output, err
	}
	defer closeDB(db)
	res, err := db.Query("SELECT * FROM bookmarks")
	if err != nil {
		return output, err
	}
	defer closeResult(res)
	for res.Next() {
		var bookmark Bookmark
		err = res.Scan(&bookmark.Version, &bookmark.Url, &bookmark.Name,
			&bookmark.Desc, &bookmark.DateAdded, &bookmark.DateModified,
			&bookmark.DateLastUsed, &bookmark.Guid, &bookmark.Icon)
		if err != nil {
			continue
		}
		output = append(output, bookmark)
	}
	return output, err
}

func closeResult(transaction *sql.Rows) {
	if transaction != nil {
		_ = transaction.Close()
	}
}

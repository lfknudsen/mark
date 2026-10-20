package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"modernc.org/sqlite"
)

const connectionString = "file:data/mark.db"

const insertStatement = `INSERT INTO bookmarks
    (version, url, name, desc, date_added, date_modified, date_last_used, guid, icon) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? )`
const insertCollectionStatement = `INSERT INTO folders
    (version, name, desc, date_added, date_modified, date_last_used, guid, icon) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ? )`

const insertSubfolderRelationStatement = `INSERT INTO subfolder_relations
    (parent, child)
    VALUES (?, ?)`

const insertChildRelationStatement = `INSERT INTO child_relations
    (parent, child)
    VALUES (?, ?)`

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

// initDatabase initialises the tables in the bookmark database.
func initDatabase() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer closeDB(db)
	tran, err := db.Begin()
	if err != nil {
		return err
	}
	defer rollback(tran)
	data, err := os.ReadFile("model.sql")
	if err != nil {
		return err
	}
	stmn, err := tran.Prepare(string(data))
	if err != nil {
		return err
	}
	defer closeStmt(stmn)
	res, err := stmn.Exec()
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

func insert(bookmarks ...Bookmark) (rowsAffected int64, err error) {
	db, err := openDB()
	if err != nil {
		return rowsAffected, err
	}
	defer closeDB(db)
	tx, err := db.Begin()
	if err != nil {
		return rowsAffected, err
	}
	stmnt, err := tx.Prepare(insertStatement)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	defer closeStmt(stmnt)
	for _, b := range bookmarks {
		res, err := stmnt.Exec(b.Version, b.Url, b.Name, b.Desc, b.DateAdded,
			b.DateModified, b.DateLastUsed, b.Guid, b.Icon)
		if err != nil {
			_ = tx.Rollback()
			return rowsAffected, err
		}
		affected, err := res.RowsAffected()
		rowsAffected += affected
		if err != nil {
			_ = tx.Rollback()
			return rowsAffected, err
		}
	}
	return rowsAffected, tx.Commit()
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

func insertCollection(collections ...*BookmarkCollection) (rowsAffected int64, err error) {
	db, err := openDB()
	if err != nil {
		return rowsAffected, err
	}
	defer closeDB(db)
	tx, err := db.Begin()
	if err != nil {
		return rowsAffected, err
	}
	stmnt, err := tx.Prepare(insertCollectionStatement)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	defer closeStmt(stmnt)
	for _, c := range collections {
		res, err := stmnt.Exec(c.Version, c.Name, c.Desc, c.DateAdded,
			c.DateModified, c.DateLastUsedContent, c.Guid, c.Icon)
		if err != nil {
			_ = tx.Rollback()
			return rowsAffected, err
		}
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return rowsAffected, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return rowsAffected, err
	}
	closeStmt(stmnt)
	for _, c := range collections {
		affected, err := insertCollection(c.Subfolders...)
		rowsAffected += affected
		if err != nil {
			return rowsAffected, err
		}
		for _, sub := range c.Subfolders {
			affected, err = setSubfolder(c, sub)
			rowsAffected += affected
			if err != nil {
				return rowsAffected, err
			}
		}
		affected, err = insert(c.Children...)
		rowsAffected += affected
		if err != nil {
			return rowsAffected, err
		}
		for _, child := range c.Children {
			affected, err = setChild(c, &child)
			rowsAffected += affected
			if err != nil {
				return rowsAffected, err
			}
		}
	}
	return rowsAffected, nil
}

func setSubfolder(parent *BookmarkCollection, child *BookmarkCollection) (rowsAffected int64, err error) {
	db, err := openDB()
	if err != nil {
		return rowsAffected, err
	}
	defer closeDB(db)
	tx, err := db.Begin()
	if err != nil {
		return rowsAffected, err
	}
	stmnt, err := tx.Prepare(insertSubfolderRelationStatement)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	defer closeStmt(stmnt)
	res, err := stmnt.Exec(parent.Guid, child.Guid)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	affected, err := res.RowsAffected()
	rowsAffected += affected
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	return rowsAffected, tx.Commit()
}

func setChild(parent *BookmarkCollection, child *Bookmark) (rowsAffected int64, err error) {
	db, err := openDB()
	if err != nil {
		return rowsAffected, err
	}
	defer closeDB(db)
	tx, err := db.Begin()
	if err != nil {
		return rowsAffected, err
	}
	stmnt, err := tx.Prepare(insertChildRelationStatement)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	defer closeStmt(stmnt)
	res, err := stmnt.Exec(parent.Guid, child.Guid)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	affected, err := res.RowsAffected()
	rowsAffected += affected
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	return rowsAffected, tx.Commit()
}

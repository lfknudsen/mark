package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"mark/pkg"
)

var connectionString = fmt.Sprintf("file:%s/%s", pkg.MarkDir, pkg.DatabaseName)

const insertStatement = `INSERT INTO bookmarks
    (version, url, name, desc, date_added, date_modified, date_last_used, guid, icon) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? )`
const insertCollectionStatement = `INSERT INTO folders
    (version, name, desc, date_added, date_modified, date_last_used, guid, icon) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ? )`

const updateChildRelationStatement = `UPDATE child_relations SET parent = ? WHERE child = ?`
const updateSubfolderRelationStatement = `UPDATE subfolder_relations SET parent = ? WHERE child = ?`

const ROOT_UUID = "0"

func executeScript(filename string, tran *sql.Tx) (rowsAffected int64, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	stmn, err := tran.Prepare(string(data))
	if err != nil {
		return 0, err
	}
	res, err := stmn.Exec()
	if err != nil {
		closeStmt(stmn)
		return 0, err
	}
	closeStmt(stmn)
	return res.RowsAffected()
}

func ExecuteScript(filename string) (rowsAffected int64, err error) {
	db, err := openDB()
	if err != nil {
		return 0, err
	}
	defer closeDB(db)
	tran, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer rollback(tran)
	rowsAffected, err = executeScript(filename, tran)
	return rowsAffected, tran.Commit()
}

// InitDatabase initialises the tables in the bookmark database.
func InitDatabase() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer closeDB(db)
	tran, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete contents of database
	rows1, err := executeScript("clean.sql", tran)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error executing clean.sql:")
		return err
	}

	// Initialise tables
	rows2, err := executeScript("model.sql", tran)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error executing model.sql:")
		return err
	}

	// Inserting root folder
	stmn, err := tran.Prepare(insertCollectionStatement)
	defer closeStmt(stmn)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	rows, err := stmn.Exec(pkg.CurrentVersionBookmarkCollection, ":root", "", now, now, now, "0", "")
	if err != nil {
		return err
	}
	rows3, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	// Create triggers
	// Must be created after root folder has been inserted to avoid it being a child of itself
	rows4, err := executeScript("triggers.sql", tran)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error executing triggers.sql:")
		return err
	}

	// Seed the database with some sample data.
	rows5, err := executeScript("seed.sql", tran)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error executing seed.sql:")
		return err
	}

	fmt.Printf("Rows affected: %d\n", rows1+rows2+rows3+rows4+rows5)
	return tran.Commit()
}

func Insert(bookmarks ...pkg.Bookmark) (rowsAffected int64, err error) {
	return InsertChild(ROOT_UUID, bookmarks...)
}

func ListBookmarks() ([]pkg.Bookmark, error) {
	output := make([]pkg.Bookmark, 0, 10)
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
		var bookmark pkg.Bookmark
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

func insertCollectionAsSub(tx *sql.Tx, parentUuid string, collections ...*pkg.BookmarkCollection) (rowsAffected int64, err error) {
	fmt.Println("Beginning transaction...")
	fmt.Println("Preparing statement 1...")
	stmnt1, err := tx.Prepare(insertCollectionStatement)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	defer closeStmt(stmnt1)

	fmt.Println("Preparing statement 2...")
	stmnt2, err := tx.Prepare(updateSubfolderRelationStatement)
	if err != nil {
		closeStmt(stmnt1)
		_ = tx.Rollback()
		return rowsAffected, err
	}
	defer closeStmt(stmnt2)

	for _, c := range collections {
		fmt.Println("Executing statement 1...")
		res, err := stmnt1.Exec(c.Version, c.Name, c.Desc, c.DateAdded,
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
		fmt.Println("Executing statement 2...")
		res, err = stmnt2.Exec(parentUuid, c.Guid)
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
	closeStmt(stmnt1)
	closeStmt(stmnt2)
	for _, c := range collections {
		affected, err := insertCollectionAsSub(tx, c.Guid, c.Subfolders...)
		rowsAffected += affected
		if err != nil {
			return rowsAffected, err
		}

		affected, err = insertChild(tx, ROOT_UUID, c.Children...)
		rowsAffected += affected
		if err != nil {
			return rowsAffected, err
		}
		for _, child := range c.Children {
			affected, err = setChild(tx, c.Guid, child.Guid)
			rowsAffected += affected
			if err != nil {
				return rowsAffected, err
			}
		}
	}
	fmt.Println("Finished inserting collection...")
	return rowsAffected, nil
}

func InsertCollectionAsSub(parentUuid string, collections ...*pkg.BookmarkCollection) (rowsAffected int64, err error) {
	fmt.Println("Opening the database...")
	db, err := openDB()
	if err != nil {
		return rowsAffected, err
	}
	defer closeDB(db)
	fmt.Println("Inserting collections...")
	tx, _ := db.Begin()

	rowsAffected, err = insertCollectionAsSub(tx, parentUuid, collections...)
	if err != nil {
		return rowsAffected, tx.Rollback()
	}
	return rowsAffected, tx.Commit()
}

func InsertCollection(collections ...*pkg.BookmarkCollection) (rowsAffected int64, err error) {
	return InsertCollectionAsSub(ROOT_UUID, collections...)
}

func setSubfolder(parentUuid string, childUuid string) (rowsAffected int64, err error) {
	db, err := openDB()
	if err != nil {
		return rowsAffected, err
	}
	defer closeDB(db)
	tx, err := db.Begin()
	if err != nil {
		return rowsAffected, err
	}
	stmnt, err := tx.Prepare(updateSubfolderRelationStatement)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	defer closeStmt(stmnt)
	res, err := stmnt.Exec(parentUuid, childUuid)
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

func setChild(tx *sql.Tx, parentUuid string, childUuid string) (rowsAffected int64, err error) {
	statement, err := tx.Prepare(updateChildRelationStatement)
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	defer closeStmt(statement)
	res, err := statement.Exec(parentUuid, childUuid)
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
	return rowsAffected, nil
}

func SetChild(parentUuid string, childUuid string) (rowsAffected int64, err error) {
	db, err := openDB()
	if err != nil {
		return rowsAffected, err
	}
	defer closeDB(db)
	tx, _ := db.Begin()
	rowsAffected, err = setChild(tx, parentUuid, childUuid)
	if err != nil {
		return rowsAffected, tx.Rollback()
	}
	return rowsAffected, tx.Commit()
}

func insertChild(tx *sql.Tx, parentUuid string, bookmarks ...pkg.Bookmark) (rowsAffected int64, err error) {
	fmt.Println("Begun a transaction.")
	stmnt, err := tx.Prepare(insertStatement)
	fmt.Println("Prepared a statement.")
	if err != nil {
		_ = tx.Rollback()
		return rowsAffected, err
	}
	for _, b := range bookmarks {
		res, err := stmnt.Exec(b.Version, b.Url, b.Name, b.Desc, b.DateAdded,
			b.DateModified, b.DateLastUsed, b.Guid, b.Icon)
		fmt.Println("Executed a statement.")
		if err != nil {
			closeStmt(stmnt)
			_ = tx.Rollback()
			return rowsAffected, err
		}
		affected, err := res.RowsAffected()
		rowsAffected += affected
		if err != nil {
			closeStmt(stmnt)
			_ = tx.Rollback()
			return rowsAffected, err
		}
	}
	closeStmt(stmnt)
	fmt.Println("Closed a statement.")
	stmnt, err = tx.Prepare(updateChildRelationStatement)
	fmt.Println("Prepared an update statement.")
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	defer closeStmt(stmnt)
	for _, b := range bookmarks {
		res, err := stmnt.Exec(parentUuid, b.Guid)
		fmt.Println("Executed an update statement.")
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
		affected, err := res.RowsAffected()
		rowsAffected += affected
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}
	fmt.Println("Closed the update statement. Now committing.")
	return rowsAffected, nil
}

func InsertChild(parentUuid string, bookmarks ...pkg.Bookmark) (rowsAffected int64, err error) {
	db, err := openDB()
	fmt.Println("Opened the database.")
	if err != nil {
		return rowsAffected, err
	}
	defer closeDB(db)
	tx, err := db.Begin()
	if err != nil {
		return rowsAffected, err
	}
	rowsAffected, err = insertChild(tx, parentUuid, bookmarks...)
	if err != nil {
		return rowsAffected, tx.Rollback()
	}
	return rowsAffected, tx.Commit()
}

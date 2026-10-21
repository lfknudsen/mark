package db

import (
	"database/sql"
	"log"

	"modernc.org/sqlite"
)

func openDB() (*sql.DB, error) {
	conn, err := sqlite.NewConnector(connectionString)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(conn)
	return db, db.Ping()
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

func closeResult(transaction *sql.Rows) {
	if transaction != nil {
		_ = transaction.Close()
	}
}

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

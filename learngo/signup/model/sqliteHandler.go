package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteHandler struct {
	db *sql.DB
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func newSqliteHandler(filepath string) DBHandler {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}

	statement, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS User(
			userId VARCHAR(20) PRIMARY KEY,
			name VARCHAR(20),
			password VARCHAR(20),
			birthDate VARCHAR(20),
			gender VARCHAR(10),
			phone VARCHAR(20),
			location VARCHAR(20)
			);

		CREATE TABLE Item(
			name VARCHAR(20) PRIMARY KEY,
			price VARCHAR(20),
			market VARCHAR(16),
			created_at VARCHAR(20),
			updated_at VARCHAR(20),
			FOREIGN KEY(userId)
				REFERENCES User(userId)
				);
		`)

	statement.Exec()
	return &sqliteHandler{db: db}
}

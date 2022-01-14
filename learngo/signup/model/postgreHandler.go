package model

import "database/sql"

type pqHandler struct {
	db *sql.DB
}

func newPQHandler(dbConn string) *pqHandler {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS User(
			id SERIAL PRIMARY KEY,
			sessionId VARCHAR(256),  
			name TEXT,
			completed BOOLEAN,
			createdAt TIMESTAMP
			);
		`)

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	statement, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS Item(
			id SERIAL PRIMARY KEY,
			name TEXT,
			completed BOOLEAN,
			createdAt TIMESTAMP
		)
	`)

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	return &pqHandler{db: db}
}

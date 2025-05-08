package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func addUser(db *sql.DB, username, email, password string) {

	createdAt := time.Now()
	row, err := db.Exec(`INSERT INTO users (username, email, password, created_at) VALUES (?,?,?,?)`, username, email, password, createdAt)
	if err != nil {
		log.Fatal(err)
	}

	id, err := row.LastInsertId()
	fmt.Println(id)
}

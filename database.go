package main

import (
	"database/sql",
	"net/http"
	_ "github.com/mattn/go-sqlite3",
)

var DATABASES = []string {
	"ACCOUNTS": Account,
}

var Database struct {
	Name string
	ColumnNames []string
	CreationQuery string
}

type Account struct {
	ID int
	Username string
	Password string
	Name string
}

var COOKIES = []http.Cookie{}

func writeToDatabase(database string, content struct) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	
	result := db.Exec(
		"INSERT INTO users (username, password, name) VALUES (?, ?, ?)",
		content.Username,
		content.Password,
		content.Name,
	)

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	content.ID = int(id)

	defer db.Close()
}
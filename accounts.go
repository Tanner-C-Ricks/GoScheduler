package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"
)

var ACCOUNTS = Database{
	Name: "Accounts",
	ColumnNames: []string{
		"accID",
		"accUsername",
		"accPassword",
		"accName",
	},
	FileURL: "./db/accounts.db",
	CreationQuery: `
	CREATE TABLE IF NOT EXISTS Accounts (
	accID INTEGER PRIMARY KEY AUTOINCREMENT,
	accUsername TEXT UNIQUE,
	accPassword TEXT,
	accName TEXT
	);
	`,
}

type Account struct {
	ID       int
	Username string
	Password string
	Name     string
}

func createAccount(content Account) any {
	db, err := sql.Open("sqlite", ACCOUNTS.FileURL)
	if err != nil {
		fmt.Println("In the first one")
		return err
	}

	_, tableCreationErr := db.Exec(ACCOUNTS.CreationQuery)
	if tableCreationErr != nil {
		return tableCreationErr
	}
	fmt.Printf("In the second one, %d", string(ACCOUNTS.FileURL))

	result, err := db.Exec(
		"INSERT INTO Accounts (accUsername, accPassword, accName) VALUES (?, ?, ?)",
		content.Username,
		content.Password,
		content.Name,
	)

	if err != nil {
		fmt.Printf("In the second one, %d", ACCOUNTS.FileURL)

		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("In the third one")

		return err
	}

	content.ID = int(id)
	defer db.Close()
	return nil
}

func getAccount(acc Account) Account {
	account := Account{}

	db, err := sql.Open("sqlite", ACCOUNTS.FileURL)
	if err != nil {
		fmt.Println(err)
		return account
	}

	query := `
	SELECT * FROM Accounts WHERE accUsername = ? AND accPassword = ?
	`
	result, err := db.Query(query, acc.Username, acc.Password)

	if err != nil {
		fmt.Println(err)
	}
	for result.Next() {
		err = result.Scan(&account.ID, &account.Username, &account.Password, &account.Name)

		if err != nil {
			fmt.Println(err)
		}
	}
	defer result.Close()

	err = result.Err()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(account)

	return account
}

func getAccountID(id int) Account {
	account := Account{
		ID: id,
	}

	db, err := sql.Open("sqlite", ACCOUNTS.FileURL)
	if err != nil {
		fmt.Println(err)
		return account
	}

	query := `
	SELECT * FROM Accounts WHERE accID = ?
	`
	result, err := db.Query(query, account.ID)

	if err != nil {
		fmt.Println(err)
	}
	for result.Next() {
		err = result.Scan(&account.ID, &account.Username, &account.Password, &account.Name)

		if err != nil {
			fmt.Println(err)
		}
	}
	defer result.Close()

	err = result.Err()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(account)

	return account
}

func accountInfo(w http.ResponseWriter, r *http.Request) Account {
	fmt.Println("Retrieving account info")
	loggedIn := checkLoggedIn(w, r)
	account := Account{ID: -1}
	if loggedIn != "LOGGEDIN" {
		return account
	}
	cookie, _ := r.Cookie("session_id")

	account = getAccountID(sessions[cookie.Value])
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(account)
	return account
}

func accountName(w http.ResponseWriter, r *http.Request) {
	account := accountInfo(w, r)

	if account.ID == -1 {
		return
	}
	name := account.Name
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(name)
}

func accountID(w http.ResponseWriter, r *http.Request) int {
	account := accountInfo(w, r)

	if account.ID == -1 {
		return -1
	}
	id := account.ID
	return id
}

func renderAccountID(w http.ResponseWriter, r *http.Request) {
	id := accountID(w, r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

func getAccounts() []Account {
	var accounts []Account

	db, err := sql.Open("sqlite", ACCOUNTS.FileURL)
	if err != nil {
		fmt.Println("In the first one")
		fmt.Println(err)
		return accounts
	}

	query := `
	SELECT * FROM Accounts;
	`

	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return accounts
	}
	defer result.Close()

	for result.Next() {
		var acc Account
		err := result.Scan(&acc.ID, &acc.Username, &acc.Password, &acc.Name)
		if err != nil {
			fmt.Println(err)
			return accounts
		}

		accounts = append(accounts, acc)
	}

	err = result.Err()

	if err != nil {
		fmt.Println(err)
		return accounts
	}
	fmt.Print(accounts)
	return accounts
}

func deleteAccount(acc Account) bool {
	db, err := sql.Open("sqlite", ACCOUNTS.FileURL)
	if err != nil {
		fmt.Println(err)
		return false
	}

	query := `
	DELETE FROM Accounts WHERE accID = ?
	`
	result, err := db.Exec(query, acc.ID)

	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("Deleted row %d", result)
	defer db.Close()
	return true
}

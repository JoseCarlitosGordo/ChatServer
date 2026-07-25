package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func main() {
	database, err := sql.Open("sqlite", "ChatServerDB")

	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	defer database.Close()
	database.Exec("DROP TABLE IF EXISTS Users")
	database.Exec(`CREATE TABLE IF NOT EXISTS Users
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		description TEXT,
		hashedPassword TEXT,
		salt BLOB

	)`)
	//database.Exec("CREATE INDEX idx_users_username ON Users (username)")
	getRows(database)

}
func getRows(db *sql.DB) {
	var rows int
	err := db.QueryRow("Select 1 from Users Limit 1").Scan(&rows)
	if err == sql.ErrNoRows {
		fmt.Print("Is Empty")
	} else {
		fmt.Print("what")
	}

}

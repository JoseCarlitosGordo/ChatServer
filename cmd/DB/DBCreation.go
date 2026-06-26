package db

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
	database.Exec(`CREATE TABLE IF NOT EXISTS Users
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		description TEXT,
		hashedPassword TEXT,
		salt TEXT

	)`)

}

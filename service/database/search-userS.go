package database

import (
	"wasaText/service/structs"

	_ "github.com/mattn/go-sqlite3"
)

var queryGetUsers = `SELECT userID, username FROM User WHERE username LIKE ?`

func (db *appdbimpl) SearchUsers(search string) ([]structs.User, error) {
	var users []structs.User

	// Ensure search includes wildcards for partial matches
	search = "%" + search + "%"

	rows, err := db.c.Query(queryGetUsers, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Properly closing rows

	for rows.Next() {
		var u structs.User
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	// Check for any iteration errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

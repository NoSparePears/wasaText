package database

import (
	"wasaText/service/structs"

	_ "github.com/mattn/go-sqlite3"
)

var queryGetUsers = `SELECT userID, username FROM User WHERE username = ?`

func (db *appdbimpl) SearchUsers(search string) ([]structs.User, error) {
	var users []structs.User

	rows, err := db.c.Query(queryGetUsers, search)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}

		var u structs.User
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

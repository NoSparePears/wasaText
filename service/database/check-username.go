package database

import (
	"database/sql"
	"errors"
)

// Query for search user with a specified username
var queryCheckUsername = `SELECT username FROM user WHERE username = ?`

func (db *appdbimpl) CheckUsername(username string) (bool, error) {
	var existsName string
	err := db.c.QueryRow(queryCheckUsername, username).Scan(&existsName)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return true, err
}

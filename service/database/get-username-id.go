package database

import (
	"database/sql"
	"errors"
	"wasaText/service/structs"
)

// Query for search user with a specified username
var queryFindUsername = `SELECT username FROM user WHERE userID = ?`

func (db *appdbimpl) GetUsernameByID(id int) (structs.User, error) {
	var user structs.User

	err := db.c.QueryRow(queryFindUsername, id).Scan(&user.Username)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return structs.User{}, nil
	}

	return user, err
}

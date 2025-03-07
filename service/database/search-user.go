package database

import (
	"database/sql"
	"errors"
	"wasaText/service/structs"
)

var query_SEARCHUSER = `SELECT userID, username FROM User WHERE username = ?;`

func (db *appdbimpl) SearchUser(username string) (structs.User, error) {
	var user structs.User
	err := db.c.QueryRow(query_SEARCHUSER, username).Scan(&user.ID, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structs.User{}, errors.New("user not found")
		}
		return structs.User{}, errors.New("internal server error")
	}
	return user, nil
}

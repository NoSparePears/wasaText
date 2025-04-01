package database

import (
	"database/sql"
	"errors"
	"wasaText/service/structs"
)

var query_SEARCHUSERID = `SELECT userID, username FROM User WHERE userID = ?;`

func (db *appdbimpl) SearchUserID(id int) (structs.User, error) {
	var user structs.User
	err := db.c.QueryRow(query_SEARCHUSERID, id).Scan(&user.ID, &user.Username)
	// Check for errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structs.User{}, errors.New("user not found")
		}
		return structs.User{}, errors.New("internal server error")
	}
	return user, nil
}

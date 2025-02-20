package database

import "wasaText/service/structs"

var query_SEARCHUSER = `SELECT userID, username FROM User WHERE username = ?;`

func (db *appdbimpl) SearchUser(username string) (structs.User, error) {
	var user structs.User
	user.ID = -1
	err := db.c.QueryRow(query_SEARCHUSER, username).Scan(&user.ID, &user.Username)
	return user, err
}

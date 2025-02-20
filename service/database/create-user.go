package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"wasaText/service/api/utils"
	"wasaText/service/structs"
)

var query_ADDUSER = `INSERT INTO User (userID, username) VALUES (?, ?);`

var query_MAXID = `SELECT MAX(userID) FROM User`

func (db *appdbimpl) CreateUser(username string) (structs.User, error) {
	var user structs.User
	user.Username = username

	//FIND userID
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query_MAXID)
	if err != nil {
		return user, err
	}

	var maxID int
	for row.Next() {
		if row.Err() != nil {
			return user, err
		}

		err = row.Scan(&_maxID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return user, err
		}

		if !_maxID.Valid {
			maxID = 0
		} else {
			maxID = int(_maxID.Int64)
		}
	}

	//SET USERID
	user.ID = maxID + 1

	//CREATE USER FOLDER
	path := "./storage/profiles" + fmt.Sprint(user.ID) + "/media"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return user, err
	}

	//SET DEFAULT PROPIC
	source, err := os.Open("./storage/default_propic.jpg")
	if err != nil {
		return user, err
	}
	defer source.Close()

	destination, err := os.Create(utils.GetProfilePicPath(user.ID))
	if err != nil {
		return user, err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return user, err
	}

	//INSERT USER
	_, err = db.c.Exec(query_ADDUSER, user.ID, user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}

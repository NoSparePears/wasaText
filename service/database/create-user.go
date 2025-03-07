package database

import (
	"fmt"
	"io"
	"os"
	"wasaText/service/api/utils"
	"wasaText/service/structs"
)

func (db *appdbimpl) CreateUser(username string) (structs.User, error) {
	var user structs.User
	user.Username = username

	// Creiamo un nuovo userID
	result, err := db.c.Exec("INSERT INTO User (username) VALUES (?);", username)
	if err != nil {
		return structs.User{}, err
	}
	// Otteniamo l'ID appena creato
	userID64, err := result.LastInsertId()
	if err != nil {
		return structs.User{}, err
	}
	user.ID = int(userID64)

	// CREATE USER FOLDER
	path := "./storage/" + fmt.Sprint(user.ID) + "/media"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return user, err
	}

	// SET DEFAULT PROPIC
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
	user.UserPropic64, err = utils.ImageToBase64(utils.GetProfilePicPath(user.ID))
	if err != nil {
		return user, err
	}
	return user, nil
}

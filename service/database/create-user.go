package database

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"wasaText/service/api/utils"
	"wasaText/service/structs"
)

func (db *appdbimpl) CreateUser(username string) (structs.User, error) {
	var user structs.User
	user.Username = username

	// Insert new user into the database
	result, err := db.c.Exec("INSERT INTO `User` (username) VALUES (?);", username)
	if err != nil {
		return structs.User{}, err
	}

	// Get the newly created user ID
	userID64, err := result.LastInsertId()
	if err != nil {
		return structs.User{}, err
	}
	user.ID = int(userID64)

	// Create user folder
	userStoragePath := fmt.Sprintf("./storage/%d/media", user.ID)
	if err := os.MkdirAll(userStoragePath, os.ModePerm); err != nil {
		return user, err
	}

	// Set default profile picture
	defaultPfp := "./storage/default_propic.jpg"
	profilePicPath := utils.GetProfilePicPath(user.ID)

	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(profilePicPath), os.ModePerm); err != nil {
		return user, err
	}

	// Copy default profile picture
	source, err := os.Open(defaultPfp)
	if err != nil {
		return user, err
	}
	defer source.Close()

	destination, err := os.Create(profilePicPath)
	if err != nil {
		return user, err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, source); err != nil {
		return user, err
	}

	// Store profile picture **path** in the database instead of Base64
	_, err = db.c.Exec("UPDATE `User` SET pfpPath = ? WHERE userID = ?;", profilePicPath, user.ID)
	if err != nil {
		return user, err
	}

	return user, nil
}

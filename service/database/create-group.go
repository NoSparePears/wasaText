package database

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"wasaText/service/api/utils"
	"wasaText/service/structs"
)

func (db *appdbimpl) CreateGroup(groupName string, userID int) (structs.Group, error) {
	result, err := db.c.Exec("INSERT INTO GlobalConversation (isGroup, groupName) VALUES (?, ?);", true, groupName)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to create group in GlobalConversation: %w", err)
	}

	globalConvoID64, err := result.LastInsertId()
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to retrieve last inserted group ID: %w", err)
	}

	group := structs.Group{
		GlobalConvoID: int(globalConvoID64),
		GroupName:     groupName,
	}

	creator, err := db.c.Exec("INSERT INTO GroupMember (groupID, userID) VALUES (?, ?);", group.GlobalConvoID, userID)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to insert user into GroupMember: %w", err)
	}

	rowsAffected, err := creator.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return structs.Group{}, errors.New("failed to insert user into GroupMember: no rows affected")
	}

	convoGroup, err := db.c.Exec("INSERT INTO Conversation (userID, destUserID, globalConvoID, lastMsgId) VALUES (?, ?, ?, ?);", userID, 0, group.GlobalConvoID, 0)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to create conversation: %w", err)
	}

	rowsAffected, err = convoGroup.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return structs.Group{}, errors.New("failed to insert conversation: no rows affected")
	}

	// Create group storage folder
	userStoragePath := fmt.Sprintf("./storage/%d/media", group.GlobalConvoID)
	if err := os.MkdirAll(userStoragePath, os.ModePerm); err != nil {
		return structs.Group{}, fmt.Errorf("failed to create group media directory: %w", err)
	}

	// Set default profile picture
	defaultPfp := "./storage/default_propic.jpg"
	profilePicPath := utils.GetProfilePicPath(group.GlobalConvoID)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(profilePicPath), os.ModePerm); err != nil {
		return structs.Group{}, fmt.Errorf("failed to create profile picture directory: %w", err)
	}

	// Copy the default profile picture
	source, err := os.Open(defaultPfp)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to open default profile picture: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(profilePicPath)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to create new profile picture file: %w", err)
	}
	defer func() {
		if cerr := destination.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close profile picture file: %w", cerr)
		}
	}()

	if _, err = io.Copy(destination, source); err != nil {
		return structs.Group{}, fmt.Errorf("failed to copy profile picture: %w", err)
	}

	// Store profile picture path in DB
	_, err = db.c.Exec("UPDATE GlobalConversation SET photoPath = ? WHERE globalConvoID = ?;", profilePicPath, group.GlobalConvoID)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to update profile picture path in database: %w", err)
	}

	return group, nil
}

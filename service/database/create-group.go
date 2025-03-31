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
		return structs.Group{}, err
	}
	// Otteniamo l'ID appena creato
	globalConvoID64, err := result.LastInsertId()
	if err != nil {
		return structs.Group{}, err
	}
	var group structs.Group
	group.GlobalConvoID = int(globalConvoID64)
	group.GroupName = groupName

	creator, err := db.c.Exec("INSERT INTO GroupMember (groupID, userID) VALUES (?, ?);", group.GlobalConvoID, userID)
	if err != nil {
		return structs.Group{}, err
	}
	// Check if any rows were affected
	rowsAffected, err := creator.RowsAffected()
	if err != nil {
		return structs.Group{}, errors.New("could not verify group membership insertion: " + err.Error())
	}
	if rowsAffected == 0 {
		return structs.Group{}, errors.New("no rows were inserted, possible issue with groupID or userID")
	}

	convoGroup, err := db.c.Exec("INSERT INTO Conversation (userID, destUserID, globalConvoID, lastMsgId) VALUES (?, ?, ?, ?);", userID, 0, group.GlobalConvoID, 0)
	if err != nil {
		return structs.Group{}, err
	}
	// Check if any rows were affected
	rowsAffected, err = convoGroup.RowsAffected()
	if err != nil {
		return structs.Group{}, errors.New("could not verify group membership insertion: " + err.Error())
	}
	if rowsAffected == 0 {
		return structs.Group{}, errors.New("no rows were inserted, possible issue with groupID or userID")
	}

	// Create group folder
	userStoragePath := fmt.Sprintf("./storage/%d/media", group.GlobalConvoID)
	if err := os.MkdirAll(userStoragePath, os.ModePerm); err != nil {
		return group, err
	}

	// Set default profile picture
	defaultPfp := "./storage/default_propic.jpg"
	profilePicPath := utils.GetProfilePicPath(group.GlobalConvoID)

	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(profilePicPath), os.ModePerm); err != nil {
		return group, err
	}

	// Copy default profile picture
	source, err := os.Open(defaultPfp)
	if err != nil {
		return group, err
	}
	defer source.Close()

	destination, err := os.Create(profilePicPath)
	if err != nil {
		return group, err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, source); err != nil {
		return group, err
	}

	// Store profile picture **path** in the database instead of Base64
	_, err = db.c.Exec("UPDATE GlobalConversation SET photoPath = ? WHERE globalConvoID = ?;", profilePicPath, group.GlobalConvoID)
	if err != nil {
		return group, err
	}

	return group, nil

}

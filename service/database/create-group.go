package database

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"wasaText/service/api/utils"
	"wasaText/service/structs"
)

func (db *appdbimpl) CreateGroup(groupName string, userID int) (structs.Group, error) {
	// Insert new group into GlobalConversation
	result, err := db.c.Exec("INSERT INTO GlobalConversation (isGroup, groupName) VALUES (?, ?);", true, groupName)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to create group in GlobalConversation: %w", err)
	}

	// Retrieve the last inserted group ID
	globalConvoID64, err := result.LastInsertId()
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to retrieve last inserted group ID: %w", err)
	}
	groupID := int(globalConvoID64)

	group := structs.Group{
		GlobalConvoID: groupID,
		GroupName:     groupName,
	}

	// Insert creator as a member of the group
	_, err = db.c.Exec("INSERT INTO GroupMember (groupID, userID) VALUES (?, ?);", groupID, userID)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to insert user into GroupMember: %w", err)
	}

	// Insert into Conversation table
	_, err = db.c.Exec("INSERT INTO Conversation (userID, destUserID, globalConvoID, lastMsgId) VALUES (?, ?, ?, ?);", userID, 0, groupID, 0)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Default profile picture
	defaultPfp := "./storage/default_propic.jpg"
	groupPicPath := utils.GetGroupPhotoPath(groupID)

	// Ensure profile picture directory exists
	if err := os.MkdirAll(filepath.Dir(groupPicPath), os.ModePerm); err != nil {
		return structs.Group{}, fmt.Errorf("failed to create profile picture directory: %w", err)
	}

	// Copy default profile picture
	if err := copyFile(defaultPfp, groupPicPath); err != nil {
		return structs.Group{}, fmt.Errorf("failed to set default group profile picture: %w", err)
	}

	// Store profile picture path in DB
	_, err = db.c.Exec("UPDATE GlobalConversation SET photoPath = ? WHERE globalConvoID = ?;", groupPicPath, groupID)
	if err != nil {
		return structs.Group{}, fmt.Errorf("failed to update profile picture path in database: %w", err)
	}

	return group, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if cerr := destination.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close destination file: %w", cerr)
		}
	}()

	if _, err = io.Copy(destination, source); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

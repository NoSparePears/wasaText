package database

import (
	"errors"
	"fmt"
)

// CheckCommentEntry checks if there is an entry in the Comment table for the given msgID and userID.
// Returns nil if no entry exists, otherwise returns an error.
func (db *appdbimpl) CheckComment(msgID int, userID int) error {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM Comment WHERE msgID = ? AND senderID = ?)`
	err := db.c.QueryRow(query, msgID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("database query error: %w", err)
	}

	if exists {
		return errors.New("entry already exists in Comment table")
	}

	return nil
}

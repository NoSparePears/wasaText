package database

import (
	"strings"
)

func (db *appdbimpl) MarkGroupMessagesAsRead(msgIDs []int, groupID int, userID int) error {
	if len(msgIDs) == 0 {
		return nil
	}

	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update Checkmarks to mark as read for the current user
	query := `
		UPDATE Checkmarks
		SET read = TRUE
		WHERE msgID IN (` + placeholders(len(msgIDs)) + `) AND viewerID = ?`

	args := make([]interface{}, len(msgIDs)+1)
	for i, id := range msgIDs {
		args[i] = id
	}
	args[len(msgIDs)] = userID // Ensure only this user's checkmark is updated

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	// Check if each message has been read by all group members
	for _, msgID := range msgIDs {
		var totalUsers, readCount int

		// Get total members in the group, excluding the sender
		err := tx.QueryRow(`
			SELECT COUNT(*) 
			FROM GroupMember 
			WHERE groupID = ? 
			`, groupID).Scan(&totalUsers)
		if err != nil {
			return err
		}
		totalUsers -= 1

		// Get count of users who have read the message
		err = tx.QueryRow(`
			SELECT COUNT(*) 
			FROM Checkmarks 
			WHERE msgID = ? AND read = TRUE`, msgID).Scan(&readCount)
		if err != nil {
			return err
		}
		// If all users have read the message, delete Checkmarks and update Message table
		if readCount == totalUsers {

			/*
				_, err = tx.Exec(`DELETE FROM Checkmarks WHERE msgID = ?`, msgID)
				if err != nil {
					log.Printf("Error deleting Checkmarks for msgID %d: %v\n", msgID, err)
					return err
				}
			*/
			_, err = tx.Exec(`UPDATE Message SET groupRead = TRUE WHERE msgID = ?`, msgID)
			if err != nil {
				return err
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Helper function to generate SQL placeholders for IN (?) queries
func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}

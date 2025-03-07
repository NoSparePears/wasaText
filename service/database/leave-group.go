package database

import (
	"errors"
)

func (db *appdbimpl) LeaveGroup(userID int, groupID int) error {

	result, err := db.c.Exec(`UPDATE Conversation SET delByUser = ? WHERE userID = ? AND globalConvoID = ?;`, true, userID, groupID)
	if err != nil {
		return errors.New("internal server error")
	}

	// check delete update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve update status")
	}
	if rowsAffected == 0 {
		return errors.New("group not found or unauthorized action")
	}

	return nil
}

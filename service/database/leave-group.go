package database

import (
	"errors"
)

func (db *appdbimpl) LeaveGroup(userID int, groupID int) error {

	result, err := db.c.Exec(`DELETE FROM Conversation WHERE userID = ? AND globalConvoID = ?;`, userID, groupID)
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

	result2, err := db.c.Exec(`DELETE FROM GroupMember WHERE userID = ? AND groupID = ?;`, userID, groupID)
	if err != nil {
		return errors.New("internal server error")
	}

	// check delete update
	rowsAffected2, err := result2.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve update status")
	}
	if rowsAffected2 == 0 {
		return errors.New("group not found or unauthorized action")
	}
	return nil
}

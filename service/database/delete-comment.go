package database

import (
	"errors"
)

var query_DELETE_COMMENT = `DELETE FROM Comment WHERE msgID = ? AND senderID = ?;`

func (db *appdbimpl) DeleteComment(userID int, msgID int) error {
	// execute deletion inside db
	result, err := db.c.Exec(query_DELETE_COMMENT, msgID, userID)
	if err != nil {
		return errors.New("internal server error")
	}
	// check comment deletion
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve deletion status")
	}
	if rowsAffected == 0 {
		return errors.New("message not found or unauthorized action")
	}

	return nil
}

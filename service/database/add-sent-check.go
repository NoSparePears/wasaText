package database

import (
	"errors"
)

func (db *appdbimpl) AddSentCheck(msgID int) error {
	// Insert the checkmark into the database
	result, err := db.c.Exec(`INSERT INTO Checkmarks (msgID, sent) VALUES (?, ?);`, msgID, true)
	if err != nil {
		return errors.New("failed to insert sent check: " + err.Error())
	}

	// Ensure result is valid before calling RowsAffected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve affected rows: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("message not found or unauthorized action")
	}

	return nil
}

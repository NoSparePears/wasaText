package database

import (
	"errors"
)

func (db *appdbimpl) AddSentCheck(msgID int) error {

	// Inserisce il checkmark nel database
	result, err := db.c.Exec(`INSERT INTO Checkmarks (msgID, sent) VALUES (?, ?);`, msgID, true)
	if err != nil {
		return errors.New("failed to insert sent check")
	}
	// Controlla se il checkmark è stato aggiunto
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve deletion status")
	}
	if rowsAffected == 0 {
		return errors.New("message not found or unauthorized action")
	}
	// Recupera il timestamp dal database
	var timestamp string
	err = db.c.QueryRow(query_GET_TIMESTAMP, msgID).Scan(&timestamp)
	if err != nil {
		return errors.New("failed to retrieve checkmark timestamp")
	}

	update, err := db.c.Exec("UPDATE Checkmarks SET sentTime = ? WHERE msgID = ? AND sent = ?", timestamp, msgID, true)
	if err != nil {
		return errors.New("internal server error")
	}

	// check name update
	rowsUpdated, err := update.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve update status")
	}
	if rowsUpdated == 0 {
		return errors.New("group not found or unauthorized action")
	}

	return nil
}

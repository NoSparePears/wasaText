package database

import (
	"errors"
)

var query_GET_MESSAGES = `SELECT m.*
FROM Message m
JOIN Checkmarks c ON m.msgID = c.msgID
WHERE m.convoID = ? 
  AND c.sent = TRUE 
  AND c.read = FALSE;`

func (db *appdbimpl) AddReadCheck(convoID int) error {
	// controllo se ci sono messaggi nella data conversazione
	// controllo quali di quei messsaggi non hanno il campo read = true
	// modifico quel campo a false
	// aggiungo timestamp
	rows, err := db.c.Query(query_GET_MESSAGES, convoID)
	if err != nil {
		return errors.New("internal server error")
	}
	defer rows.Close()

	// Store message IDs that need to be updated
	var msgIDs []int

	for rows.Next() {
		var msgID int
		if err := rows.Scan(&msgID); err != nil {
			return errors.New("internal server error")
		}
		msgIDs = append(msgIDs, msgID)
	}

	if err := rows.Err(); err != nil {
		return errors.New("internal server error")
	}

	// If no messages need updating, return early
	if len(msgIDs) == 0 {
		return nil
	}

	// Update all selected messages
	tx, err := db.c.Begin() // Start a transaction, so that the update is atomic
	if err != nil {
		return errors.New("failed to start transaction")
	}

	stmt, err := tx.Prepare("UPDATE Checkmarks SET read = ?, readTime = CURRENT_TIMESTAMP WHERE msgID = ?")
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.New("failed to prepare update statement; rollback also failed")
		}
		return errors.New("failed to prepare update statement")
	}
	defer stmt.Close()

	for _, msgID := range msgIDs {
		_, err := stmt.Exec(true, msgID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return errors.New("failed to update message read status; rollback also failed")
			}
			return errors.New("failed to update message read status")
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.New("failed to commit transaction; rollback also failed")
		}
		return errors.New("failed to commit transaction")
	}

	return nil

}

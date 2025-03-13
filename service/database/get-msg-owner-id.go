package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) GetMsgOwnerID(msgID int) (int, error) {
	var ownerID int
	err := db.c.QueryRow("SELECT senderID FROM Message WHERE msgID = ?", msgID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("message not found: %w", err)

		}
		return 0, err
	}
	return ownerID, err
}

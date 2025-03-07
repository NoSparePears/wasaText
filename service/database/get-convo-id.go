package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetConvoID(userID int, recID int) (int, error) {
	var convoID int
	// execute query and find convo
	row := db.c.QueryRow("SELECT globalConvoID FROM Conversation WHERE userID = ? AND destUserID = ?;", userID, recID)
	err := row.Scan(&convoID)

	if err != nil {
		// Check if no rows were found
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("conversation not found")
		}
		return 0, errors.New("failed to retrieve conversation ID")
	}

	return convoID, nil

}

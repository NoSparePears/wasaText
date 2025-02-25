package database

import (
	"database/sql"
	"errors"
)

var query_INSERT_MEMBER = `INSERT INTO GroupMember (groupID, userID) VALUES (?, ?);`

func (db *appdbimpl) AddToGroup(userID int, groupID int) error {
	// Execute the insert query
	result, err := db.c.Exec(query_INSERT_MEMBER, groupID, userID)
	if err != nil {
		// Handle duplicate entry error (specific to MySQL)
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user is already a member of this group")
		}
		return errors.New("failed to insert member into group: " + err.Error())
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not verify group membership insertion: " + err.Error())
	}
	if rowsAffected == 0 {
		return errors.New("no rows were inserted, possible issue with groupID or userID")
	}

	return nil

}

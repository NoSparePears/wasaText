package database

import (
	"errors"
)

var query_CHANGEUSERNAME = `UPDATE User SET username = ? WHERE userID = ?;`

func (db *appdbimpl) ChangeUsername(userID int, newUsername string) error {
	// Check if the new username is already taken before attempting to update
	taken, err := db.CheckUsername(newUsername)
	if err != nil {
		return errors.New("internal server error") // Error while checking the username
	}
	if taken {
		return errors.New("Username already taken. Please choose another one.") // Username already taken
	}

	// Update the username in the database
	result, err := db.c.Exec(query_CHANGEUSERNAME, newUsername, userID)
	if err != nil {
		return errors.New("internal server error") // Internal server error
	}

	// Check if the update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve update status") // Error retrieving update status
	}
	if rowsAffected == 0 {
		return errors.New("user not found or unauthorized action") // No rows affected
	}

	return nil
}

package database

import "errors"

var query_CHANGEUSERNAME = `UPDATE User SET username = ? WHERE userID = ?;`

func (db *appdbimpl) ChangeUsername(userID int, newUsername string) error {
	result, err := db.c.Exec(query_CHANGEUSERNAME, newUsername, userID)
	if err != nil {
		return errors.New("internal server error")
	}

	// check username update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve update status")
	}
	if rowsAffected == 0 {
		return errors.New("user not found or unauthorized action")
	}

	return nil
}

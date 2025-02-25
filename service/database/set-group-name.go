package database

import "errors"

func (db *appdbimpl) SetGroupName(groupID int, name string) error {
	result, err := db.c.Exec("UPDATE GlobalConversation SET groupName = ? WHERE globalConvoID = ? AND isGroup = ?", name, groupID, true)
	if err != nil {
		return errors.New("internal server error")
	}

	//check name update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to retrieve update status")
	}
	if rowsAffected == 0 {
		return errors.New("group not found or unauthorized action")
	}

	return nil

}

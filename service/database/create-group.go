package database

import (
	"errors"
	"wasaText/service/structs"
)

func (db *appdbimpl) CreateGroup(groupName string, userID int) (structs.Group, error) {

	result, err := db.c.Exec("INSERT INTO GlobalConversation (isGroup, groupName) VALUES (?, ?);", true, groupName)
	if err != nil {
		return structs.Group{}, err
	}
	// Otteniamo l'ID appena creato
	globalConvoID64, err := result.LastInsertId()
	if err != nil {
		return structs.Group{}, err
	}
	var group structs.Group
	group.GlobalConvoID = int(globalConvoID64)
	group.GroupName = groupName

	admin, err := db.c.Exec("INSERT INTO GroupMember (groupID, userID, role) VALUES (?, ?. ?);", group.GlobalConvoID, userID, "admin")
	if err != nil {
		return structs.Group{}, err
	}
	// Check if any rows were affected
	rowsAffected, err := admin.RowsAffected()
	if err != nil {
		return structs.Group{}, errors.New("could not verify group membership insertion: " + err.Error())
	}
	if rowsAffected == 0 {
		return structs.Group{}, errors.New("no rows were inserted, possible issue with groupID or userID")
	}

	return group, nil

}

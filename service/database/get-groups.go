package database

import (
	"errors"
	"wasaText/service/structs"
)

var query_GETGROUPS = ` SELECT 
					    c.globalConvoID, 
					    c.lastMsgID, 
					    g.groupName
					FROM Conversation c
					JOIN GlobalConversation g ON c.globalConvoID = g.globalConvoID
					WHERE c.userID = ? 
					AND g.isGroup = TRUE;`

func (db *appdbimpl) GetGroups(userID int) ([]structs.Group, error) {
	rows, err := db.c.Query(query_GETGROUPS, userID)
	if err != nil {
		return []structs.Group{}, errors.New("internal server error")
	}
	defer rows.Close()

	var groups []structs.Group

	for rows.Next() {
		var group structs.Group
		err = rows.Scan(&group.GlobalConvoID, &group.LastMsgID, &group.GroupName)
		if err != nil {
			return []structs.Group{}, errors.New("internal server error")
		}
		groups = append(groups, group)
	}
	// check for errors encountered during iteration
	if rows.Err() != nil {
		return nil, errors.New("internal server error")
	}

	return groups, nil
}

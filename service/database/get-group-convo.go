package database

import "wasaText/service/structs"

var query_GETGROUPCONVO = `SELECT gc.groupName, c.lastMsgID
	FROM GlobalConversation gc
	JOIN Conversation c ON gc.globalConvoID = c.globalConvoID
	WHERE gc.globalConvoID = ?;`

func (db *appdbimpl) GetGroupConvo(groupID int) (structs.Group, error) {
	// execute  query and find conversation
	rows, err := db.c.Query(query_GETGROUPCONVO, groupID)
	if err != nil {
		return structs.Group{}, err
	}
	defer rows.Close()

	var group structs.Group
	// check if query returned a result
	if rows.Next() {
		// Scan values into struct fields
		err := rows.Scan(&group.GroupName, &group.LastMsgID)
		if err != nil {
			return structs.Group{}, err
		}

		// Check for additional errors in row iteration
		if rows.Err() != nil {
			return structs.Group{}, err
		}
		return group, nil //success
	}

	return structs.Group{}, err
}

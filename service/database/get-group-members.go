package database

import (
	"errors"
	"wasaText/service/structs"
)

var query_GETMEMBERS = `SELECT userID FROM GroupMember WHERE groupID = ?`

func (db *appdbimpl) GetGroupMembers(group_id int) ([]structs.User, error) {
	// get members from database
	rows, err := db.c.Query(query_GETMEMBERS, group_id)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	defer rows.Close()

	var members []structs.User
	for rows.Next() {
		var member structs.User
		err = rows.Scan(&member.ID)
		if err != nil {
			return nil, errors.New("failed to retrieve member")
		}
		members = append(members, member)
	}
	// check for errors encountered during iteration
	if rows.Err() != nil {
		return nil, errors.New("internal server error")
	}

	return members, nil
}

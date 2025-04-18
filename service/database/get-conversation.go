package database

import (
	"errors"
	"wasaText/service/structs"
)

var query_GETCONVERSATION = `SELECT * from Conversation WHERE userID = ? AND destUserID = ?;`

func (db *appdbimpl) GetConversation(userID int, recID int) (structs.Conversation, error) {
	// execute  query and find conversation
	rows, err := db.c.Query(query_GETCONVERSATION, userID, recID)
	if err != nil {
		return structs.Conversation{}, errors.New("internal server error")
	}
	defer rows.Close()

	var convo structs.Conversation
	// check if query returned a result
	if rows.Next() {
		// Scan values into struct fields
		err := rows.Scan(&convo.UserID, &convo.DestUserID, &convo.GlobalConvoID, &convo.LastMsgID, &convo.DelByUser, &convo.Visible)
		if err != nil {
			return structs.Conversation{}, errors.New("failed to retrieve conversation")
		}

		// Ensure that the retrieved `destUserID` matches `recID`
		if convo.DestUserID != recID {
			return structs.Conversation{}, errors.New("conversation mismatch error")
		}

		// Check for additional errors in row iteration
		if rows.Err() != nil {
			return structs.Conversation{}, errors.New("internal server error")
		}

		return convo, nil // success
	}

	// no conversation found, so create one
	return db.CreateConversation(userID, recID)

}

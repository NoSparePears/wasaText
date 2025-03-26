package database

import (
	"errors"
	"wasaText/service/structs"
)

var query_GETCONVERSATIONS = `
	SELECT c.destUserID, c.globalConvoID, c.lastMsgID 
  	FROM Conversation c
  	JOIN GlobalConversation g ON c.globalConvoID = g.globalConvoID
  	WHERE c.userID = ? 
  	AND g.isGroup = FALSE;`

func (db *appdbimpl) GetConversations(user_id int) ([]structs.Conversation, error) {
	// get conversations from database
	rows, err := db.c.Query(query_GETCONVERSATIONS, user_id)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	defer rows.Close()

	var convos []structs.Conversation

	for rows.Next() {
		var convo structs.Conversation
		err = rows.Scan(&convo.DestUserID, &convo.GlobalConvoID, &convo.LastMsgID)
		if err != nil {
			return nil, err
		}
		convos = append(convos, convo)
	}
	// check for errors encountered during iteration
	if rows.Err() != nil {
		return nil, errors.New("internal server error")
	}

	return convos, nil
}

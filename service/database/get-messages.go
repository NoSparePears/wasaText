package database

import (
	"database/sql"
	"errors"
	"wasaText/service/structs"
)

var query_GETCONVERSATIONID = `SELECT globalConvoID from Conversation WHERE userID = ? AND destUserID = ?;`

func (db *appdbimpl) GetMessages(userID int, recID int) ([]structs.Message, error) {
	var globalConvoID int

	// Execute query to find the conversation ID
	err := db.c.QueryRow(query_GETCONVERSATIONID, userID, recID).Scan(&globalConvoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("conversation not found")
		}
		return nil, errors.New("failed to retrieve conversationID")
	}

	//get messages from database
	rows, err := db.c.Query("SELECT msgID, convoID, senderID, content, timestamp FROM Message WHERE convoID = ? ORDER BY timestamp ASC;", globalConvoID)
	if err != nil {
		return nil, errors.New("couldnt get messages from selected convo")
	}
	if rows == nil {
		return nil, nil
	}
	defer rows.Close()

	var messages []structs.Message

	for rows.Next() {
		var message structs.Message
		err = rows.Scan(&message.MsgID, &message.ConvoID, &message.SenderID, &message.Content, &message.Timestamp)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}
	//check for errors encountered during iteration
	if rows.Err() != nil {
		return nil, errors.New("errore durante scan rows")
	}

	return messages, nil
}

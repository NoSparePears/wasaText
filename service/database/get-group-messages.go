package database

import (
	"errors"
	"wasaText/service/structs"
)

func (db *appdbimpl) GetGroupMessages(groupID int) ([]structs.Message, error) {

	// get messages from database
	rows, err := db.c.Query("SELECT msgID, convoID, senderID, content, timestamp FROM Message WHERE convoID = ? ORDER BY timestamp ASC;", groupID)
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
	// check for errors encountered during iteration
	if rows.Err() != nil {
		return nil, errors.New("errore durante scan rows")
	}

	return messages, nil
}

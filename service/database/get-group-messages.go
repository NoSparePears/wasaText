package database

import (
	"errors"
	"wasaText/service/structs"
)

// GetGroupMessages retrieves messages from a group and ensures CheckReceived is correctly set
func (db *appdbimpl) GetGroupMessages(userID int, groupID int) ([]structs.Message, error) {

	// Fetch messages
	rows, err := db.c.Query(query_MSGINFO, groupID)
	if err != nil {
		return nil, errors.New("couldn't get messages from selected conversation")
	}
	defer rows.Close()

	var messages []structs.Message

	for rows.Next() {
		var message structs.Message

		// Scan the message fields
		err = rows.Scan(
			&message.MsgID, &message.ConvoID, &message.SenderID, &message.Content,
			&message.Timestamp, &message.IsPhoto, &message.IsForwarded, &message.CheckSent, &message.CheckReceived,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	// Handle row scanning error
	if err = rows.Err(); err != nil {
		return nil, errors.New("error during row scanning")
	}

	return messages, nil
}

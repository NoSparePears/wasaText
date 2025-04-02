package database

import (
	"errors"
	"wasaText/service/structs"
)

// GetGroupMessages retrieves messages from a group and ensures CheckReceived is correctly set
func (db *appdbimpl) GetGroupMessages(userID int, groupID int) ([]structs.Message, error) {

	query := `
		SELECT 
		    msgID, 
		    convoID, 
		    senderID, 
		    content, 
		    timestamp, 
		    isPhoto, 
		    isForwarded, 
		    groupRead 
		FROM Message 
		WHERE convoID = ?
		ORDER BY timestamp ASC;
	`

	// Fetch messages
	rows, err := db.c.Query(query, groupID)
	if err != nil {
		return nil, errors.New("could not retrieve messages from selected conversation")
	}
	defer rows.Close()

	var messages []structs.Message
	var unreadMsgIDs []int

	for rows.Next() {
		var message structs.Message

		// Scan the message fields
		err = rows.Scan(
			&message.MsgID, &message.ConvoID, &message.SenderID, &message.Content,
			&message.Timestamp, &message.IsPhoto, &message.IsForwarded,
			&message.CheckReceived,
		)
		if err != nil {
			return nil, err
		}

		// Set CheckSent to true for all messages
		message.CheckSent = true

		// If message is unread by the current user and was not sent by them, mark for update
		if !message.CheckReceived && message.SenderID != userID {
			unreadMsgIDs = append(unreadMsgIDs, message.MsgID)
		}

		messages = append(messages, message)
	}

	// Check if there are any unread messages to update
	if len(unreadMsgIDs) > 0 {
		err = db.MarkGroupMessagesAsRead(unreadMsgIDs, groupID, userID)
		if err != nil {
			return nil, err
		}
	}
	// Check if there were any errors during row scanning
	if rows.Err() != nil {
		return nil, errors.New("error during row scan")
	}

	return messages, nil
}

package database

import (
	"database/sql"
	"errors"
	"wasaText/service/structs"
)

var query_GETCONVERSATIONID = `SELECT globalConvoID FROM Conversation WHERE userID = ? AND destUserID = ?;`

var query_MSGINFO = `
	SELECT m.msgID, m.convoID, m.senderID, m.content, m.timestamp, 
	       COALESCE(c.sent, FALSE), COALESCE(c.read, FALSE)
	FROM Message m
	LEFT JOIN Checkmarks c ON m.msgID = c.msgID
	WHERE m.convoID = ?
	ORDER BY m.timestamp ASC;
`

func (db *appdbimpl) GetMessages(userID int, recID int) ([]structs.Message, error) {
	var globalConvoID int

	// Execute query to find the conversation ID
	err := db.c.QueryRow(query_GETCONVERSATIONID, userID, recID).Scan(&globalConvoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("conversation not found")
		}
		return nil, errors.New("failed to retrieve conversation ID")
	}

	// Get messages from database
	rows, err := db.c.Query(query_MSGINFO, globalConvoID)
	if err != nil {
		return nil, errors.New("couldn't get messages from selected conversation")
	}
	defer rows.Close()

	var messages []structs.Message
	var unreadMsgIDs []int // Collect unread message IDs

	for rows.Next() {
		var message structs.Message

		// Scan the message fields, ensuring boolean values don't become NULL
		err = rows.Scan(
			&message.MsgID, &message.ConvoID, &message.SenderID, &message.Content,
			&message.Timestamp, &message.CheckSent, &message.CheckReceived,
		)
		if err != nil {
			return nil, err
		}

		// If the message is unread and sent by the recipient, mark it as read
		if !message.CheckReceived && message.SenderID != userID {
			unreadMsgIDs = append(unreadMsgIDs, message.MsgID)
			message.CheckReceived = true
		}

		messages = append(messages, message)
	}

	// Batch update all unread messages
	if len(unreadMsgIDs) > 0 {
		err = db.MarkMessagesAsRead(unreadMsgIDs)
		if err != nil {
			return nil, err
		}
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, errors.New("error during row scanning")
	}

	return messages, nil
}

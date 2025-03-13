package database

import (
	"database/sql"
	"fmt"
	"wasaText/service/structs"
)

func (db *appdbimpl) GetLastMessage(convoID int, msgID int) (structs.Message, error) {
	var msg structs.Message

	query := `SELECT 
		      	m.senderID, 
		      	m.content, 
		      	m.timestamp, 
		      	m.isPhoto, 
		      	c.sent, 
		      	c.read
		      FROM Message m
		      LEFT JOIN Checkmarks c ON m.msgID = c.msgID
		      WHERE m.convoID = ? AND m.msgID = ?`

	err := db.c.QueryRow(query, convoID, msgID).Scan(
		&msg.SenderID,
		&msg.Content,
		&msg.Timestamp,
		&msg.IsPhoto,
		&msg.CheckSent,
		&msg.CheckReceived,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return structs.Message{}, fmt.Errorf("last message not found: %w", err)
		}
		return structs.Message{}, fmt.Errorf("database error: %w", err)
	}

	msg.MsgID = msgID
	msg.ConvoID = convoID

	return msg, nil
}

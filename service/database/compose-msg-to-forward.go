package database

import (
	"database/sql"
	"errors"
	"wasaText/service/structs"
)

var query_GETMSG = `SELECT content FROM Message WHERE msgID = ? AND convoID = ?;`

func (db *appdbimpl) ComposeMsgToForward(msgID int, convoID int, customContent string) (structs.Message, error) {
	var msg structs.Message
	err := db.c.QueryRow(query_GETMSG, msgID, convoID).Scan(&msg.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structs.Message{}, errors.New("message not found")
		}
		return structs.Message{}, errors.New("failed to retrieve original message")
	}
	// Costruzione del contenuto finale del messaggio inoltrato
	finalContent := msg.Content
	if customContent != "" {
		finalContent = customContent + " Inoltrato:  + originalMsg.Content + "
	}

	msg.Content = finalContent

	return msg, nil
}

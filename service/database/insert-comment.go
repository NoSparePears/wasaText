package database

import (
	"errors"
	"wasaText/service/structs"
)

var query_INSERT_COMMENT = `INSERT INTO Comment (msgID, senderID, emoji) VALUES (?, ?, ?);`

func (db *appdbimpl) InsertComment(c structs.Comment) (structs.Comment, error) {
	// insert comment into table
	result, err := db.c.Exec(query_INSERT_COMMENT, c.MsgID, c.SenderID, c.Emoji)
	if err != nil {
		return structs.Comment{}, errors.New("failed to insert comment")
	}
	// Ottiene l'ID del commento appena creato
	cID, err := result.LastInsertId()
	if err != nil {
		return structs.Comment{}, errors.New("failed to retrieve comment ID")
	}
	// Set the CommID field of the comment struct
	c.CommID = int(cID)

	return c, nil
}

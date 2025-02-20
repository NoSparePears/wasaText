package database

import (
	"errors"
	"wasaText/service/structs"
)

var (
	query_INSERT_MESSAGE  = `INSERT INTO Message (globalConvoID, senderID, content) VALUES (?, ?, ?);`
	query_UPDATE_LAST_MSG = `UPDATE Conversation SET lastMsgId = ? WHERE userID = ? AND globalConvoID = ?;`
)

// SendMessage inserisce un messaggio in una conversazione e aggiorna lo stato della chat.
func (db *appdbimpl) InsertMessage(msg structs.Message, recID int) (structs.Message, error) {
	// Inserisce il messaggio nel database
	result, err := db.c.Exec(query_INSERT_MESSAGE, msg.ConvoID, msg.SenderID, msg.Content)
	if err != nil {
		return structs.Message{}, errors.New("failed to insert message")
	}

	// Ottiene l'ID del messaggio appena creato
	msgID, err := result.LastInsertId()
	if err != nil {
		return structs.Message{}, errors.New("failed to retrieve message ID")
	}

	// Aggiorna il lastMsgId per il mittente
	_, err = db.c.Exec(query_UPDATE_LAST_MSG, msgID, msg.SenderID, msg.ConvoID)
	if err != nil {
		return structs.Message{}, errors.New("failed to update last message ID")
	}

	//crea conversazione per destinatario se non esiste
	_, err = db.GetConversation(recID, msg.SenderID)
	if err != nil {
		return structs.Message{}, errors.New("failed to create conversation for receiver")
	}
	msg.ConvoID = int(msgID)
	return msg, nil
}

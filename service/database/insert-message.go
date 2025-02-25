package database

import (
	"errors"
	"wasaText/service/structs"
)

var (
	query_INSERT_MESSAGE  = `INSERT INTO Message (globalConvoID, senderID, content) VALUES (?, ?, ?);`
	query_UPDATE_LAST_MSG = `UPDATE Conversation SET lastMsgId = ? WHERE userID = ? AND globalConvoID = ?;`
	query_GET_TIMESTAMP   = `SELECT timestamp FROM Message WHERE msgID = ?;`
)

// SendMessage inserisce un messaggio in una conversazione e aggiorna lo stato della chat.
func (db *appdbimpl) InsertMessage(msg structs.Message, recID int) (structs.Message, error) {
	// Inserisce il messaggio nel database
	result, err := db.c.Exec(query_INSERT_MESSAGE, msg.ConvoID, msg.SenderID, msg.Content)
	if err != nil {
		return structs.Message{}, errors.New("failed to insert message")
	}
	// Controlla se il messaggio è stato aggiunto
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return structs.Message{}, errors.New("failed to retrieve deletion status")
	}
	if rowsAffected == 0 {
		return structs.Message{}, errors.New("message not found or unauthorized action")
	}
	// Ottiene l'ID del messaggio appena creato
	msgID, err := result.LastInsertId()
	if err != nil {
		return structs.Message{}, errors.New("failed to retrieve message ID")
	}
	// Recupera il timestamp dal database
	var timestamp string
	err = db.c.QueryRow(query_GET_TIMESTAMP, msgID).Scan(&timestamp)
	if err != nil {
		return structs.Message{}, errors.New("failed to retrieve message timestamp")
	}
	msg.Timestamp = timestamp // Assegna il timestamp al messaggio
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

	err = db.AddSentCheck(msg.MsgID)
	if err != nil {
		return structs.Message{}, errors.New("failed to create sent checkmark")
	}
	return msg, nil
}

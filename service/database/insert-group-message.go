package database

import (
	"errors"
	"fmt"
	"wasaText/service/structs"
)

// SendMessage inserisce un messaggio in una conversazione e aggiorna lo stato della chat.
func (db *appdbimpl) InsertGroupMessage(msg structs.Message) (int, string, error) {
	// Inserisce il messaggio nel database
	result, err := db.c.Exec(query_INSERT_MESSAGE, msg.ConvoID, msg.SenderID, msg.Content, msg.IsPhoto, msg.IsForwarded)
	if err != nil {
		return 0, "", errors.New("failed to insert message")
	}

	// Ottiene l'ID del messaggio appena creato
	msgID64, err := result.LastInsertId()
	if err != nil {
		return 0, "", errors.New("failed to retrieve message ID")
	}
	msgID := int(msgID64)

	// Recupera il timestamp dal database
	var timestamp string
	err = db.c.QueryRow(query_GET_TIMESTAMP, msgID).Scan(&timestamp)
	if err != nil {
		return 0, "", errors.New("failed to retrieve message timestamp")
	}

	// Aggiorna il lastMsgId per il mittente
	result2, err := db.c.Exec(query_UPDATE_LAST_MSG, msgID, msg.ConvoID)
	if err != nil {
		return 0, "", errors.New("failed to update last message")

	}
	// Controlla se il last msg id è stato aggiornato
	rowsAffected, err := result2.RowsAffected()
	if err != nil {
		return 0, "", errors.New("failed to retrieve update status")
	}
	if rowsAffected == 0 {
		return 0, "", fmt.Errorf("no conversation updated: conversation with ID %d not found or lastMsgId was already set to %d", msg.ConvoID, msgID)
	}

	err = db.AddSentCheck(msgID)
	if err != nil {
		return 0, "", errors.New("failed to create sent checkmark")
	}

	return msgID, timestamp, nil
}

package database

import (
	"database/sql"
	"errors"
	"wasaText/service/structs"
)

var (
	query_DELETEMSG       = `DELETE FROM Message WHERE msgID = ? AND globalConvoID = ?;`
	query_GETUSERCONVO    = `SELECT globalConvoID FROM Conversation WHERE userID = ? AND globalConvoID = ?;`
	query_GET_LAST_MSG_ID = `SELECT MAX(msgID) FROM Message WHERE globalConvoID = ?;`                       // ✅ Trova l'ultimo messaggio rimanente
	query_UPDATE_LASTMSG  = `UPDATE Conversation SET lastMsgId = ? WHERE userID = ? AND globalConvoID = ?;` // ✅ Aggiorna il lastMsgId
)

func (db *appdbimpl) DeleteMessage(msgID int, convoID int, userID int) (structs.Conversation, error) {
	// Esegui la query di eliminazione
	result, err := db.c.Exec(query_DELETEMSG, msgID, convoID)
	if err != nil {
		return structs.Conversation{}, errors.New("internal server error")
	}

	// Controlla se il messaggio è stato eliminato
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return structs.Conversation{}, errors.New("failed to retrieve deletion status")
	}
	if rowsAffected == 0 {
		return structs.Conversation{}, errors.New("message not found or unauthorized action")
	}

	// Trova il nuovo lastMsgId per la conversazione (l'ultimo messaggio rimasto)
	var lastMsgID sql.NullInt64
	err = db.c.QueryRow(query_GET_LAST_MSG_ID, convoID).Scan(&lastMsgID)
	if err != nil {
		return structs.Conversation{}, errors.New("failed to retrieve last message ID")
	}

	// Se non ci sono più messaggi nella conversazione, lastMsgID deve essere NULL
	newLastMsgID := sql.NullInt64{Int64: 0, Valid: false}
	if lastMsgID.Valid {
		newLastMsgID = lastMsgID
	}

	// Aggiorna lastMsgId nella tabella Conversation
	_, err = db.c.Exec(query_UPDATE_LAST_MSG, newLastMsgID.Int64, userID, convoID)
	if err != nil {
		return structs.Conversation{}, errors.New("failed to update last message ID in conversation")
	}

	// Recupera la conversazione aggiornata
	row := db.c.QueryRow(query_GETUSERCONVO, userID, convoID)
	var convo structs.Conversation
	err = row.Scan(&convo.GlobalConvoID)
	if err != nil {
		return structs.Conversation{}, errors.New("failed to retrieve updated conversation")
	}

	return convo, nil

}

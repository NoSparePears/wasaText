package database

import (
	"database/sql"
	"errors"
	"wasaText/service/structs"
)

var (
	query_DELETEMSG           = `DELETE FROM Message WHERE msgID = ? AND globalConvoID = ?;`
	query_GETUSERCONVO        = `SELECT globalConvoID FROM Conversation WHERE userID = ? AND globalConvoID = ?;`
	query_GET_LAST_MSG_ID     = `SELECT MAX(msgID) FROM Message WHERE globalConvoID = ?;`
	query_UPDATE_LASTMSG      = `UPDATE Conversation SET lastMsgId = ? WHERE userID = ? AND globalConvoID = ?;`
	query_UPDATE_LASTMSG_NULL = `UPDATE Conversation SET lastMsgId = NULL WHERE userID = ? AND globalConvoID = ?;`
)

// DeleteMessage elimina un messaggio e aggiorna lastMsgId nella tabella Conversation
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

	// Trova il nuovo lastMsgId per la conversazione
	var lastMsgID sql.NullInt64
	err = db.c.QueryRow(query_GET_LAST_MSG_ID, convoID).Scan(&lastMsgID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return structs.Conversation{}, errors.New("failed to retrieve last message ID")
	}

	// Se non ci sono più messaggi, aggiorniamo lastMsgId a NULL
	if !lastMsgID.Valid {
		_, err = db.c.Exec(query_UPDATE_LASTMSG_NULL, userID, convoID)
	} else {
		_, err = db.c.Exec(query_UPDATE_LASTMSG, lastMsgID.Int64, userID, convoID)
	}
	if err != nil {
		return structs.Conversation{}, errors.New("failed to update last message ID in conversation")
	}

	// Recupera la conversazione aggiornata
	row := db.c.QueryRow(query_GETUSERCONVO, userID, convoID)
	var convo structs.Conversation
	err = row.Scan(&convo.GlobalConvoID)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.Conversation{}, errors.New("conversation not found")
	}
	if err != nil {
		return structs.Conversation{}, errors.New("failed to retrieve updated conversation")
	}

	return convo, nil
}

package database

import (
	"errors"
	"fmt"
	"strings"
	"wasaText/service/structs"
)

// SendMessage inserisce un messaggio in una conversazione e aggiorna lo stato della chat.
func (db *appdbimpl) InsertGroupMessage(msg structs.Message) (int, string, error) {

	// Start transaction
	tx, err := db.c.Begin()
	if err != nil {
		return 0, "", errors.New("failed to start transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Inserisce il messaggio nel database
	result, err := tx.Exec(query_INSERT_MESSAGE, msg.ConvoID, msg.SenderID, msg.Content, msg.IsPhoto, msg.IsForwarded)
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
	err = tx.QueryRow(query_GET_TIMESTAMP, msgID).Scan(&timestamp)
	if err != nil {
		return 0, "", errors.New("failed to retrieve message timestamp")
	}

	// Aggiorna il lastMsgId per il mittente
	result2, err := tx.Exec(query_UPDATE_LAST_MSG, msgID, msg.ConvoID)
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

	// Retrieve all group members except the sender
	rows, err := tx.Query(`
		SELECT userID FROM GroupMember WHERE groupID = ? AND userID != ?`,
		msg.ConvoID, msg.SenderID)
	if err != nil {
		return 0, "", errors.New("failed to retrieve group members")
	}
	defer rows.Close()

	// Prepare bulk insert for Checkmarks
	var checkmarkInserts []string
	var args []interface{}

	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return 0, "", errors.New("failed to retrieve group member")
		}
		checkmarkInserts = append(checkmarkInserts, "(?, ?)")
		args = append(args, msgID, userID)
	}

	// Insert into Checkmarks
	if len(checkmarkInserts) > 0 {
		query := `INSERT INTO Checkmarks (msgID, viewerID) VALUES ` +
			strings.Join(checkmarkInserts, ",")
		_, err = tx.Exec(query, args...)
		if err != nil {
			return 0, "", errors.New("failed to insert checkmarks")
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return 0, "", errors.New("failed to commit transaction")
	}

	return msgID, timestamp, nil
}

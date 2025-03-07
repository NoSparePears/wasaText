package database

import (
	"errors"
	"wasaText/service/structs"
)

var query_DELETECONVO = `DELETE FROM Conversation WHERE sendID = ? AND recID = ?;`
var query_GETUSERCONVOS = `SELECT convoID FROM Conversation WHERE sendID = ?;`

func (db *appdbimpl) DeleteConversation(user_id int, rec_id int) ([]structs.Conversation, error) {
	// execute delete query
	result, err := db.c.Exec(query_DELETECONVO, user_id, rec_id)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	// check if any rows was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.New("failed to retrieve deletion status")
	}
	if rowsAffected == 0 {
		return nil, errors.New("conversation not found or unauthorized action")
	}
	// retrieve updated list of user's convos
	rows, err := db.c.Query(query_GETUSERCONVOS, user_id)
	if err != nil {
		return nil, errors.New("failed to retrieve updated conversations")
	}
	defer rows.Close()

	var convos []structs.Conversation
	for rows.Next() {
		var convo structs.Conversation
		if err := rows.Scan(&convo.GlobalConvoID); err != nil {
			return nil, errors.New("failed to process conversations list")
		}
		convos = append(convos, convo)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("failed to retrieve complete conversations list")
	}

	return convos, nil

}

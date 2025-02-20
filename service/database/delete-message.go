package database

import (
	"errors"
	"wasaText/service/structs"
)

var query_DELETEMSG = `DELETE FROM Message WHERE msgID = ? AND convoID = ?;`
var query_GETUSERCONVO = `SELECT FROM Conversation WHERE userID = ? AND globalConvoID = ?;`

func (db *appdbimpl) DeleteMessage(msgID int, convoID int, userID int) (structs.Conversation, error) {
	//execute delete query
	result, err := db.c.Exec(query_DELETEMSG, msgID, convoID)
	if err != nil {
		return structs.Conversation{}, errors.New("internal server error")
	}

	//check if any rows was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return structs.Conversation{}, errors.New("failed to retrieve deletion status")
	}
	if rowsAffected == 0 {
		return structs.Conversation{}, errors.New("message not found or unauthorized action")
	}

	//retrieve updated convo
	row, err := db.c.Query(query_GETUSERCONVO, userID, convoID)
	if err != nil {
		return structs.Conversation{}, errors.New("failed to retrieve updated conversation")
	}
	defer row.Close()

	var convo structs.Conversation
	//check if query returned a result
	if row.Next() {
		if err := row.Scan(&convo.GlobalConvoID); err != nil {
			return structs.Conversation{}, errors.New(("failed to retrieve conversation"))
		}
		return convo, nil //success
	}

	return structs.Conversation{}, errors.New(("failed to retrieve conversation"))

}

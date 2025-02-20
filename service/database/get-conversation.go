package database

import (
	"errors"
	"wasaText/service/structs"
)

var query_GETCONVERSATIONID = `SELECT c1.globalConvoID
FROM Conversation c1
JOIN Conversation c2 ON c1.globalConvoID = c2.globalConvoID
WHERE c1.userID = ? AND c2.userID = ?;`

func (db *appdbimpl) GetConversation(userID int, recID int) (structs.Conversation, error) {
	//execute  query and find conversation
	rows, err := db.c.Query(query_GETCONVERSATIONID, userID, recID)
	if err != nil {
		return structs.Conversation{}, errors.New("internal server error")
	}
	defer rows.Close()

	var convo structs.Conversation
	//check if query returned a result
	if rows.Next() {
		if err := rows.Scan(&convo.GlobalConvoID); err != nil {
			return structs.Conversation{}, errors.New(("failed to retrieve conversation"))
		}
		return convo, nil //success
	}

	//no conversation found, so create one
	return db.CreateConversation(userID, recID)

}

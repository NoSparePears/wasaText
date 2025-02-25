package database

import (
	"wasaText/service/structs"
)

func (db *appdbimpl) CreateConversation(userID int, recID int) (structs.Conversation, error) {

	// Creiamo un nuovo globalConvoID
	result, err := db.c.Exec("INSERT INTO GlobalConversation (isGroup, groupName) VALUES (?, ?);", false, "")
	if err != nil {
		return structs.Conversation{}, err
	}

	// Otteniamo l'ID appena creato
	globalConvoID64, err := result.LastInsertId()
	if err != nil {
		return structs.Conversation{}, err
	}

	// Creiamo la conversazione per il primo utente
	_, err = db.c.Exec("INSERT INTO Conversation (userID, globalConvoID, visible) VALUES (?, ?, ?);", userID, int(globalConvoID64), true)
	if err != nil {
		return structs.Conversation{}, err
	}

	var convo structs.Conversation
	convo.UserID = userID
	convo.GlobalConvoID = int(globalConvoID64)

	return convo, nil

}

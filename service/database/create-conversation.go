package database

import (
	"wasaText/service/structs"
)

func (db *appdbimpl) CreateConversation(userID int, recID int) (structs.Conversation, error) {
	var globalConvoID int

	// Creiamo un nuovo globalConvoID
	result, err := db.c.Exec("INSERT INTO GlobalConversation (isGroup) VALUES (?);", false)
	if err != nil {
		return structs.Conversation{}, err
	}

	// Otteniamo l'ID appena creato
	globalConvoID64, err := result.LastInsertId()
	if err != nil {
		return structs.Conversation{}, err
	}
	globalConvoID = int(globalConvoID64)

	// Creiamo la conversazione per il primo utente
	_, err = db.c.Exec("INSERT INTO Conversation (userID, globalConvoID, visible) VALUES (?, ?, ?);", userID, globalConvoID, true)
	if err != nil {
		return structs.Conversation{}, err
	}

	var convo structs.Conversation
	convo.UserID = userID
	convo.GlobalConvoID = globalConvoID

	return convo, nil

}

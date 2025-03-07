package database

import (
	"errors"
	"wasaText/service/structs"
)

func (db *appdbimpl) CreateConversation(userID int, recID int) (structs.Conversation, error) {

	var existingCount int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Conversation WHERE userID = ? AND destUserID = ?;", userID, recID).Scan(&existingCount)
	if err != nil {
		return structs.Conversation{}, err
	}
	// If a conversation already exists, return an error (or handle it as needed)
	if existingCount > 0 {
		return structs.Conversation{}, errors.New("conversation already exists")
	}

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
	globalConvoID := int(globalConvoID64)

	// Creiamo la conversazione per il primo utente
	_, err = db.c.Exec("INSERT INTO Conversation (userID, globalConvoID, visible, destUserID) VALUES (?, ?, ?, ?);", userID, globalConvoID, true, recID)
	if err != nil {
		return structs.Conversation{}, err
	}
	// Creiamo la conversazione per il secondo utente
	_, err = db.c.Exec("INSERT INTO Conversation (userID, globalConvoID, visible, destUserID) VALUES (?, ?, ?, ?);", recID, globalConvoID, false, userID)
	if err != nil {
		return structs.Conversation{}, err
	}

	var convo structs.Conversation
	convo.UserID = userID
	convo.GlobalConvoID = globalConvoID
	convo.DestUserID = recID
	convo.Visible = true

	return convo, nil

}

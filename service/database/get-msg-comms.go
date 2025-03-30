package database

import (
	"errors"
	"wasaText/service/structs"
)

// GetMessageComments retrieves comments for a specific message and conversation
func (db *appdbimpl) GetMessageComments(msgID int, convoID int) ([]structs.Comment, error) {
	// Fetch comments for the given message and conversation
	dbComms, err := db.c.Query(`
		SELECT c.commID, c.msgID, c.senderID, c.emoji 
		FROM Comment c
		JOIN Message m ON c.msgID = m.msgID
		WHERE c.msgID = ? AND m.convoID = ?`, msgID, convoID)
	if err != nil {
		return nil, errors.New("couldn't get comments from selected message: " + err.Error())
	}
	defer dbComms.Close()

	// Iterate through the query results
	var comments []structs.Comment
	for dbComms.Next() {
		var comment structs.Comment
		err = dbComms.Scan(&comment.CommID, &comment.MsgID, &comment.SenderID, &comment.Emoji)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Check if there was an error during iteration
	if err = dbComms.Err(); err != nil {
		return nil, err
	}

	// Get usernames for each comment (Optimized version can be done)
	for i, comment := range comments {
		user, err := db.GetUsernameByID(comment.SenderID)
		if err != nil {
			return nil, err
		}
		comments[i].SendUsername = user.Username
	}

	// Return the list of comments
	return comments, nil
}

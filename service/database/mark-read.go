package database

import "errors"

func (db *appdbimpl) MarkMessagesAsRead(msgIDs []int) error {
	if len(msgIDs) == 0 {
		return nil
	}

	// Construct the SQL query dynamically
	query := `UPDATE Checkmarks SET read = TRUE WHERE msgID IN (`
	params := make([]interface{}, len(msgIDs))
	for i, id := range msgIDs {
		params[i] = id
		query += "?,"
	}
	query = query[:len(query)-1] + `);` // Remove last comma and close parenthesis

	// Execute the query
	_, err := db.c.Exec(query, params...)
	if err != nil {
		return errors.New("failed to mark messages as read")
	}

	return nil
}

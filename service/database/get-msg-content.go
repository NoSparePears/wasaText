package database

func (db *appdbimpl) GetMessageContent(msgID int) (string, int, error) {
	var content string
	var isPhoto int

	query := "SELECT content, isPhoto FROM Message WHERE msgID = ?"
	err := db.c.QueryRow(query, msgID).Scan(&content, &isPhoto)
	if err != nil {
		return "", 0, err
	}
	return content, isPhoto, nil
}
